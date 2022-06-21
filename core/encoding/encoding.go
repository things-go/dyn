package encoding

import (
	"sort"
	"sync"
)

// Codec defines the interface Transport uses to encode and decode messages.  Note
// that implementations of this interface must be thread safe; a Codec's
// methods can be called from concurrent goroutines.
type Codec interface {
	// Name returns the name of the Codec implementation. The returned string
	// will be used as part of content type in transmission.  The result must be
	// static; the result cannot change between calls.
	Name() string
	// Marshal returns the wire format of v.
	Marshal(v any) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v any) error
}

var (
	codecsMu sync.RWMutex
	codecs   = make(map[string]Codec)
)

// Register registers the provided Codec.
// if register called twice for Codec name, it will overwrite previous Codec.
func Register(codec Codec) {
	codecsMu.Lock()
	defer codecsMu.Unlock()
	if codec == nil {
		panic("encoding: register a nil Codec")
	}
	if codec.Name() == "" {
		panic("encoding: register Codec with empty string for Name()")
	}
	codecs[codec.Name()] = codec
}

// Codecs returns a sorted list of the names of the registered Codec.
func Codecs() []string {
	codecsMu.RLock()
	defer codecsMu.RUnlock()
	names := make([]string, 0, len(codecs))
	for name := range codecs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetCodec gets a registered Codec by name, or nil if no Codec is
// registered for the name.
func GetCodec(name string) Codec {
	codecsMu.RLock()
	defer codecsMu.RUnlock()
	return codecs[name]
}
