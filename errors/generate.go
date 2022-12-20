//go:generate protoc -I . -I ../example/third_party --go_out=paths=source_relative:. errors.proto
package errors
