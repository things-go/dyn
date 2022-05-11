package builder

import (
	"io"
	"os"
	"runtime"
	"text/template"
)

// populated by ldflags
var (
	// Version 版本 由外部ldflags指定
	Version string
	// GitCommit git提交版本(短) 由外部ldflags指定
	GitCommit string
	// GitFullCommit git提交版本(完整) 由外部ldflags指定
	GitFullCommit string
	// GitTag git的tag 由外部ldflags指定
	GitTag string
	// BuildDate 编译日期 由外部ldflags指定
	BuildDate string
)

const versionTpl = `##################################################################
#    Version:          {{.Version}}
#    Git Commit:       {{.GitCommit}}
#    Git Full Commit:  {{.GitFullCommit}}
#    Git Tag:          {{.GitTag}}
#    Build Date:       {{.BuildDate}}
#    Go Version:       {{.GoVersion}}
#    OS/Arch:          {{.GOOS}}/{{.GOARCH}}
#    Go Max Procs:     {{.GoMaxProcs}}
#    Num CPU:          {{.NumCPU}}
{{- if .Deploy }}
#    Deploy:           {{.Deploy}}
{{- end }}
{{- if .Metadata }}
#    Metadata:
{{- end }}
{{- range $k, $v := .Metadata }}
#      {{$k}}: {{$v}}
{{- end }}
##################################################################
`

// Information 版本信息以及一些布署信息
type Information struct {
	Version       string
	GitCommit     string
	GitFullCommit string
	GitTag        string
	BuildDate     string
	GoVersion     string
	GOOS          string
	GOARCH        string
	GoMaxProcs    int
	NumCPU        int
	Deploy        string
	Metadata      map[string]string
}

type Option func(info *Information)

// WithDeploy 设置布署信息
func WithDeploy(deploy string) Option {
	return func(info *Information) {
		info.Deploy = deploy
	}
}

// WithMetadata 设置其它基本信息
func WithMetadata(m map[string]string) Option {
	return func(info *Information) {
		for k, v := range m {
			if k != "" && v != "" {
				info.Metadata[k] = v
			}
		}
	}
}

// Println 打印版本信息以及一些布署信息至os.Stdout
func Println(opts ...Option) {
	Fprintln(os.Stdout, opts...)
}

func Fprintln(w io.Writer, opts ...Option) {
	info := Information{
		Version,
		GitCommit,
		GitFullCommit,
		GitTag,
		BuildDate,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.GOMAXPROCS(0),
		runtime.NumCPU(),
		"",
		make(map[string]string),
	}
	for _, opt := range opts {
		opt(&info)
	}

	template.Must(template.New("version").Parse(versionTpl)).
		Execute(w, info) // nolint: errcheck
}
