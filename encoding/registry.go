package encoding

import (
	"errors"
	"fmt"
	"mime"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/things-go/dyn/encoding/codec"
	"github.com/things-go/dyn/encoding/form"
	"github.com/things-go/dyn/encoding/json"
	"github.com/things-go/dyn/encoding/jsonpb"
	"github.com/things-go/dyn/encoding/msgpack"
	"github.com/things-go/dyn/encoding/proto"
	"github.com/things-go/dyn/encoding/toml"
	"github.com/things-go/dyn/encoding/xml"
	"github.com/things-go/dyn/encoding/yaml"
)

var (
	acceptHeader      = http.CanonicalHeaderKey("Accept")
	contentTypeHeader = http.CanonicalHeaderKey("Content-Type")

	globalRegistry = registry{
		mimeMap: map[string]codec.Marshaler{
			MIMEWildcard: &HTTPBodyCodec{
				Marshaler: &jsonpb.Codec{
					MarshalOptions: protojson.MarshalOptions{
						UseEnumNumbers: true,
					},
					UnmarshalOptions: protojson.UnmarshalOptions{
						DiscardUnknown: true,
					},
				},
			},
			MIMEQuery: &form.QueryCodec{Codec: form.New("json")},
			MIMEURI:   &form.UriCodec{Codec: form.New("json")},

			MIMEPOSTForm:          form.New("json"),
			MIMEMultipartPOSTForm: &form.MultipartCodec{Codec: form.New("json")},
			MIMEJSON:              &json.Codec{},
			MIMEXML:               &xml.Codec{},
			MIMEXML2:              &xml.Codec{},
			MIMEPROTOBUF:          &proto.Codec{},
			MIMEMSGPACK:           &msgpack.Codec{},
			MIMEMSGPACK2:          &msgpack.Codec{},
			MIMEYAML:              &yaml.Codec{},
			MIMETOML:              &toml.Codec{},
		},
	}
)

// registry is a mapping from MIME types to Marshalers.
type registry struct {
	mimeMap map[string]codec.Marshaler
}

// RegisterMarshaler register a marshaler for a case-sensitive MIME type string
// ("*" to match any MIME type).
// you can override default marshaler with same MIME type
// if marshaler is nil, it will remove the MIME type marshaler.
func RegisterMarshaler(mime string, marshaler codec.Marshaler) error {
	if len(mime) == 0 {
		return errors.New("empty MIME type")
	}
	if marshaler == nil {
		if mime == MIMEWildcard ||
			mime == MIMEQuery ||
			mime == MIMEURI {
			return fmt.Errorf("MIME(%s) can't delete, you can override", mime)
		}
		delete(globalRegistry.mimeMap, mime)
	} else {
		globalRegistry.mimeMap[mime] = marshaler
	}
	return nil
}

func GetMarshaler(mime string) codec.Marshaler {
	m := globalRegistry.mimeMap[mime]
	if m == nil {
		m = globalRegistry.mimeMap[MIMEWildcard]
	}
	return m
}

// MarshalerFromRequest returns the inbound or outbound marshalers for this request.
// It checks the registry on the globalRegistry for the MIME type set by the Content-Type header.
// If it isn't set (or the request Content-Type is empty), checks for "*".
// If there are multiple Content-Type headers set, choose the first one that it can
// exactly match in the registry.
// Otherwise, it follows the above logic for "*"/Inbound/Outbound Marshaler.
func MarshalerFromRequest(r *http.Request, inbound bool) codec.Marshaler {
	var marshaler codec.Marshaler

	if inbound {
		for _, contentTypeVal := range r.Header[contentTypeHeader] {
			contentType, _, err := mime.ParseMediaType(contentTypeVal)
			if err != nil {
				continue
			}
			if m, ok := globalRegistry.mimeMap[contentType]; ok {
				marshaler = m
				break
			}
		}
	} else {
		for _, acceptVal := range r.Header[acceptHeader] {
			if m, ok := globalRegistry.mimeMap[acceptVal]; ok {
				marshaler = m
				break
			}
		}
	}
	if marshaler == nil {
		marshaler = globalRegistry.mimeMap[MIMEWildcard]
	}
	return marshaler
}
