package signature

import (
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// ConcatMapWithSort 拼接对象, 按key排序
// 格式: k1=v1&k2=v2 (其中sep1='=', sep2='&')
func ConcatMapWithSort(mp map[string]string, sep1, sep2 string) string {
	if len(mp) == 0 {
		return ""
	}
	ks := make([]string, 0, len(mp))
	n := 0
	for k, v := range mp {
		if v != "" {
			ks = append(ks, k)
			n += len(k) + len(v)
		}
	}
	n += len(ks)*2 - 1 // 加上sep1和sep2的总个数

	sort.Strings(ks)

	bs := strings.Builder{}
	bs.Grow(n)
	for _, k := range ks {
		if bs.Len() > 0 {
			bs.WriteString(sep2)
		}
		bs.WriteString(k)
		bs.WriteString(sep1)
		bs.WriteString(mp[k])
	}
	return bs.String()
}

// ConcatMap 拼接对象,按key排序
// hasBrace 前后是否带有大括号
// 格式: hasBrace=false, k1=v1&k2=v2
// 格式: hasBrace=true, {k1=v1&k2=v2}
func ConcatMap(mp map[string]any, hasBrace bool) string {
	if len(mp) == 0 {
		return ""
	}

	keys := make([]string, 0, len(mp))
	for k, v := range mp {
		if !(v == nil || v == "") { // ignore nil and empty string
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	first := true
	buff := &strings.Builder{}
	if hasBrace {
		buff.WriteString("{")
	}
	for _, key := range keys {
		if v := toString(mp[key]); v != "" {
			if !first {
				buff.WriteString("&")
			}
			first = false
			buff.WriteString(key)
			buff.WriteString("=")
			buff.WriteString(v)
		}
	}
	if hasBrace {
		buff.WriteString("}")
	}
	return buff.String()
}

// ConcatArray 拼接数组, 不排序, 按","分隔, 忽略空值, 前后带有方括号
func ConcatArray(v any) string {
	value := reflect.ValueOf(v)
	if !(value.Kind() == reflect.Array || value.Kind() == reflect.Slice) {
		return toString(v)
	}
	if value.Len() == 0 {
		return ""
	}

	first := true
	buff := &strings.Builder{}
	buff.WriteString("[")
	for i := 0; i < value.Len(); i++ {
		o := value.Index(i).Interface()
		if vv := toString(o); vv != "" {
			if !first {
				buff.WriteString(",")
			}
			first = false
			buff.WriteString(vv)
		}
	}
	buff.WriteString("]")
	return buff.String()
}

func toString(vv any) string {
	vv = indirectToStringer(vv)

	switch s := vv.(type) {
	case nil:
		return ""
	case string:
		return s
	case bool:
		return strconv.FormatBool(s)
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32)
	case int:
		return strconv.Itoa(s)
	case int64:
		return strconv.FormatInt(s, 10)
	case int32:
		return strconv.Itoa(int(s))
	case int16:
		return strconv.FormatInt(int64(s), 10)
	case int8:
		return strconv.FormatInt(int64(s), 10)
	case uint:
		return strconv.FormatUint(uint64(s), 10)
	case uint64:
		return strconv.FormatUint(s, 10)
	case uint32:
		return strconv.FormatUint(uint64(s), 10)
	case uint16:
		return strconv.FormatUint(uint64(s), 10)
	case uint8:
		return strconv.FormatUint(uint64(s), 10)
	case []byte:
		return string(s)
	case template.HTML:
		return string(s)
	case template.URL:
		return string(s)
	case template.JS:
		return string(s)
	case template.CSS:
		return string(s)
	case template.HTMLAttr:
		return string(s)
	case json.Number:
		return string(s)
	case fmt.Stringer:
		return s.String()
	case error:
		return s.Error()
	case map[string]any:
		return ConcatMap(s, true)
	case []any:
		if len(s) == 0 {
			return ""
		}
		return ConcatArray(s)
	default:
		value := reflect.ValueOf(s)
		if kind := value.Kind(); kind == reflect.Array || kind == reflect.Slice {
			if value.Len() == 0 {
				return ""
			}
			return ConcatArray(s)
		}
		return ""
	}
}

func indirectToStringer(vv any) any {
	if vv == nil {
		return nil
	}

	var errorType = reflect.TypeOf((*error)(nil)).Elem()
	var fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	v := reflect.ValueOf(vv)
	if !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}
