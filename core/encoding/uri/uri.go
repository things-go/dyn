package uri

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/things-go/dyn/core/encoding/form"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var defaultCodec = New()

func ReplaceDefaultCodec(codec Codec) { defaultCodec = codec }
func Encode(pathTemplate string, msg interface{}, needQuery bool) string {
	return defaultCodec.Encode(pathTemplate, msg, needQuery)
}

type Option func(*Codec)

func WithDisableProto() Option {
	return func(c *Codec) {
		c.disableProto = true
	}
}

type Codec struct {
	f            form.Codec
	disableProto bool
}

// New returns a new Codec,
// default tag name is "form",
// proto use protoJSON tag
func New(opts ...Option) Codec {
	codec := Codec{form.New(), false}
	for _, opt := range opts {
		opt(&codec)
	}
	return codec
}

// Encode encode proto message to url path.
func (c Codec) Encode(pathTemplate string, msg interface{}, needQuery bool) string {
	if msg == nil || (reflect.ValueOf(msg).Kind() == reflect.Ptr && reflect.ValueOf(msg).IsNil()) {
		return pathTemplate
	}
	reg := regexp.MustCompile(`/{[.\w]+}`)
	if reg == nil {
		return pathTemplate
	}

	path := pathTemplate
	pathParams := make(map[string]struct{})
	if c.disableProto {
		// TODO:
	} else if mg, ok := msg.(proto.Message); ok {
		path = reg.ReplaceAllStringFunc(pathTemplate, func(in string) string {
			if len(in) < 4 { //nolint:gomnd // **  explain the 4 number here :-) **
				return in
			}
			key := in[2 : len(in)-1]
			vars := strings.Split(key, ".")
			if value, err := getValueByField(mg.ProtoReflect(), vars); err == nil {
				pathParams[key] = struct{}{}
				return "/" + value
			}
			return in
		})
	} else {
		// TODO:
	}

	if needQuery {
		values, err := c.f.Encode(msg)
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

func getValueByField(v protoreflect.Message, fieldPath []string) (string, error) {
	var fd protoreflect.FieldDescriptor
	for i, fieldName := range fieldPath {
		fields := v.Descriptor().Fields()
		if fd = fields.ByJSONName(fieldName); fd == nil {
			fd = fields.ByName(protoreflect.Name(fieldName))
			if fd == nil {
				return "", fmt.Errorf("field path not found: %q", fieldName)
			}
		}
		if i == len(fieldPath)-1 {
			break
		}
		if fd.Message() == nil || fd.Cardinality() == protoreflect.Repeated {
			return "", fmt.Errorf("invalid path: %q is not a message", fieldName)
		}
		v = v.Get(fd).Message()
	}
	return form.EncodeField(fd, v.Get(fd))
}
