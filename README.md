# dyn

dyn project toolkit for gin

[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/things-go/dyn?tab=doc)
[![codecov](https://codecov.io/gh/things-go/dyn/branch/main/graph/badge.svg)](https://codecov.io/gh/things-go/dyn)
[![Tests](https://github.com/things-go/dyn/actions/workflows/ci.yml/badge.svg)](https://github.com/things-go/dyn/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/things-go/dyn)](https://goreportcard.com/report/github.com/things-go/dyn)
[![License](https://img.shields.io/github/license/things-go/dyn)](https://raw.githubusercontent.com/things-go/dyn/main/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/things-go/dyn)](https://github.com/things-go/dyn/tags)

## Usage

`dyn`是一个`gin`,`protobuf`工程以及代码工具生成器.

- `proto-gen-dyn-gin` 从 `proto` 的生成`gin`的代码.
  ***注意***: 当使用`proto-gen-go-gin`要禁用`gin`自带的`binding`,使用`gin.DisableBindValidation()` 接口
- `proto-gen-dyn-resty` 从 `proto` 的生成`resty`的代码.
- `proto-gen-dyn-enum` 从 `proto` 的生成`enum`的代码.
- `errno-gen` 从枚举生成统一错误
- `dyngen` 简化工程模板生成
  
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
