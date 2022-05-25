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

// Encode encode proto message to url path.
func Encode(pathTemplate string, msg interface{}, needQuery bool) string {
	if msg == nil || (reflect.ValueOf(msg).Kind() == reflect.Ptr && reflect.ValueOf(msg).IsNil()) {
		return pathTemplate
	}
	reg := regexp.MustCompile(`/{[.\w]+}`)
	if reg == nil {
		return pathTemplate
	}

	path := pathTemplate
	pathParams := make(map[string]struct{})
	if mg, ok := msg.(proto.Message); ok {
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
		u, err := form.Encode(msg)
		if err == nil && len(u) > 0 {
			for key := range pathParams {
				delete(u, key)
			}
			query := u.Encode()
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
