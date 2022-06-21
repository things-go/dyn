package form

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/things-go/dyn/core/encoding"
)

var defaultCodec = New()

func init() {
	encoding.Register(defaultCodec)
}
func ReplaceDefaultCodec(codec Codec)    { defaultCodec = codec }
func Name() string                       { return defaultCodec.Name() }
func Marshal(v any) ([]byte, error)      { return defaultCodec.Marshal(v) }
func Unmarshal(data []byte, v any) error { return defaultCodec.Unmarshal(data, v) }
func Encode(v any) (url.Values, error)   { return defaultCodec.Encode(v) }
func Decode(vs url.Values, v any) error  { return defaultCodec.Decode(vs, v) }
func EncodeURL(pathTemplate string, msg any, needQuery bool) string {
	return defaultCodec.EncodeURL(pathTemplate, msg, needQuery)
}

type Codec struct {
	encoder      *form.Encoder
	decoder      *form.Decoder
	tagName      string
	disableProto bool
}

type Option func(*Codec)

func WithEncoder(encoder *form.Encoder) Option {
	return func(c *Codec) {
		if encoder != nil {
			c.encoder = encoder
		}
	}
}

func WithDecoder(decoder *form.Decoder) Option {
	return func(c *Codec) {
		if decoder != nil {
			c.decoder = decoder
		}
	}
}

func WithDisableProto() Option {
	return func(c *Codec) {
		c.disableProto = true
	}
}

func WithTagName(tagName string) Option {
	return func(c *Codec) {
		c.tagName = tagName
	}
}

// New returns a new Codec,
// default tag name is "json",
// proto use protoJSON tag
func New(opts ...Option) Codec {
	encoder := form.NewEncoder()
	decoder := form.NewDecoder()
	codec := Codec{
		encoder,
		decoder,
		"json",
		false,
	}
	for _, opt := range opts {
		opt(&codec)
	}
	codec.encoder.SetTagName(codec.tagName)
	codec.decoder.SetTagName(codec.tagName)
	return codec
}

func (Codec) Name() string { return "x-www-form-urlencoded" }

func (c Codec) Marshal(v any) ([]byte, error) {
	vs, err := c.Encode(v)
	if err != nil {
		return nil, err
	}
	return []byte(vs.Encode()), nil
}

func (c Codec) Unmarshal(data []byte, v any) error {
	vs, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}
	return c.Decode(vs, v)
}

func (c Codec) Encode(v any) (url.Values, error) {
	var vs url.Values
	var err error

	if c.disableProto {
		vs, err = c.encoder.Encode(v)
	} else if m, ok := v.(proto.Message); ok {
		vs, err = EncodeValues(m)
	} else {
		vs, err = c.encoder.Encode(v)
	}
	if err != nil {
		return nil, err
	}
	for k, vv := range vs {
		if len(vv) == 0 {
			delete(vs, k)
		}
	}
	return vs, nil
}

func (c Codec) Decode(vs url.Values, v any) error {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	if !c.disableProto {
		if m, ok := v.(proto.Message); ok {
			return DecodeValues(m, vs)
		}
		if m, ok := reflect.Indirect(reflect.ValueOf(v)).Interface().(proto.Message); ok {
			return DecodeValues(m, vs)
		}
	}
	return c.decoder.Decode(v, vs)
}

// EncodeURL encode msg to url path.
// pathTemplate is a template of url path like http://helloworld.dev/{name}/sub/{sub.name},
func (c Codec) EncodeURL(pathTemplate string, msg any, needQuery bool) string {
	if msg == nil || (reflect.ValueOf(msg).Kind() == reflect.Ptr && reflect.ValueOf(msg).IsNil()) {
		return pathTemplate
	}
	reg := regexp.MustCompile(`/{[.\w]+}`)
	if reg == nil {
		return pathTemplate
	}

	pathParams := make(map[string]struct{})
	repl := func(in string) string {
		if len(in) < 4 { //nolint:gomnd
			return in
		}
		key := in[2 : len(in)-1]
		vars := strings.Split(key, ".")
		if value, err := getValueWithField(msg, vars, c.tagName); err == nil {
			pathParams[key] = struct{}{}
			return "/" + value
		}
		return in
	}
	if !c.disableProto {
		if mg, ok := msg.(proto.Message); ok {
			repl = func(in string) string {
				if len(in) < 4 { //nolint:gomnd
					return in
				}
				key := in[2 : len(in)-1]
				vars := strings.Split(key, ".")
				if value, err := getValueFromProtoWithField(mg.ProtoReflect(), vars); err == nil {
					pathParams[key] = struct{}{}
					return "/" + value
				}
				return in
			}
		}
	}
	path := reg.ReplaceAllStringFunc(pathTemplate, repl)

	if needQuery {
		values, err := c.Encode(msg)
		if err == nil && len(values) > 0 {
			for key := range pathParams {
				delete(values, key)
			}
			query := values.Encode()
			if query != "" {
				path += "?" + query
			}
		}
	}
	return path
}

func getValueFromProtoWithField(v protoreflect.Message, fieldPath []string) (string, error) {
	var fd protoreflect.FieldDescriptor

	for i, fieldName := range fieldPath {
		fields := v.Descriptor().Fields()
		if fd = fields.ByJSONName(fieldName); fd == nil {
			fd = fields.ByName(protoreflect.Name(fieldName))
			if fd == nil {
				return "", fmt.Errorf("form: field path not found: %q", fieldName)
			}
		}
		if i == len(fieldPath)-1 {
			break
		}
		if fd.Message() == nil || fd.Cardinality() == protoreflect.Repeated {
			return "", fmt.Errorf("form: invalid path, %q is not a message", fieldName)
		}
		v = v.Get(fd).Message()
	}
	return EncodeField(fd, v.Get(fd))
}

func getValueWithField(s any, fieldPath []string, tagName string) (string, error) {
	v := reflect.ValueOf(s)
	// if pointer get the underlying element
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return "", errors.New("form: not struct")
	}
	for i, fieldName := range fieldPath {
		fields := findField(v, fieldName, tagName)
		if !fields.IsValid() {
			return "", fmt.Errorf("form: field path not found: %q", fieldName)
		}
		v = fields
		if i == len(fieldPath)-1 {
			break
		}
	}
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return cast.ToString(v), nil
}

func findField(v reflect.Value, searchName, tagName string) reflect.Value {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		v = reflect.New(v.Type().Elem())
	}
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return reflect.Value{}
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fv := v.Field(i)
		// we can't access the value of unexported fields
		if !fv.CanInterface() || field.PkgPath != "" {
			continue
		}
		// don't check if it's omitted
		tag := field.Tag.Get(tagName)
		if tag == "-" {
			continue
		}
		name := field.Name
		tagNamed, _ := parseTag(tag)
		if tagNamed != "" {
			name = tagNamed
		}
		if name == searchName {
			return v.FieldByName(field.Name)
		}
	}
	return reflect.Value{}
}
