# dyn
dyn project toolkit for gin 

[![GoDoc](https://godoc.org/github.com/things-go/dyn?status.svg)](https://godoc.org/github.com/things-go/dyn)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/things-go/dyn?tab=doc)
[![codecov](https://codecov.io/gh/things-go/dyn/branch/main/graph/badge.svg)](https://codecov.io/gh/things-go/dyn)
[![Tests](https://github.com/things-go/dyn/actions/workflows/ci.yml/badge.svg)](https://github.com/things-go/dyn/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/things-go/dyn)](https://goreportcard.com/report/github.com/things-go/dyn)
[![Licence](https://img.shields.io/github/license/things-go/dyn)](https://raw.githubusercontent.com/things-go/dyn/main/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/things-go/dyn)](https://github.com/things-go/dyn/tags)

## Usage

`dyn`是一个`gin`工程工具生成器, 集成了 [proto-gen-go-errno](cmd/proto-gen-go-errno) 和 [proto-gen-go-gin](cmd/proto-gen-go-gin)

- `proto-gen-go-errno` 从`proto` 枚举统一生成错误
- `proto-gen-go-gin` 从 `proto` 的生成`gin`的代码. 

***注意***: 当使用`proto-gen-go-gin`要禁用`gin`自带的`binding`,使用`gin.DisableBindValidation()` 接口

### Installation

Use go get.
```bash
    go get github.com/things-go/dyn
```

Then import the package into your own code.
```bash
    import "github.com/things-go/dyn"
```

### Example

[embedmd]:# (_examples/main.go go)
```go

```

## References

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.