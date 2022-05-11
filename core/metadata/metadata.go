package metadata

// Metadata is a map[string]string of key/value pairs.
type Metadata map[string]string

// New creates an MD from a given key-values map.
func New(mds ...map[string]string) Metadata {
	md := Metadata{}
	for _, m := range mds {
		for k, v := range m {
			md.Set(k, v)
		}
	}
	return md
}

// Get returns the value associated with the passed key.
func (m Metadata) Get(key string) string { return m[key] }

// Set stores the key-value pair.
func (m Metadata) Set(key, value string) {
	if key == "" || value == "" {
		return
	}
	m[key] = value
}

// Set stores the key-value pair.
func (m Metadata) Delete(key string) { delete(m, key) }

func (m Metadata) Exist(key string) bool {
	_, ok := m[key]
	return ok
}

// Range iterate over element in metadata.
func (m Metadata) Range(f func(k, v string) bool) {
	for k, v := range m {
		b := f(k, v)
		if !b {
			return
		}
	}
}

// Clone returns a deep copy of Metadata
func (m Metadata) Clone() Metadata {
	md := Metadata{}
	for k, v := range m {
		md[k] = v
	}
	return md
}
