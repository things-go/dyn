package command

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sourcegraph/conc"
	"github.com/spf13/cobra"
)

type ErrnoOption struct {
	Pattern         []string
	Type            []string
	Tags            []string
	DisableStringer bool
	Epk             string
}

type RootCmd struct {
	cmd   *cobra.Command
	level string
	ErrnoOption
}

func NewRootCmd() *RootCmd {
	root := &RootCmd{}
	cmd := &cobra.Command{
		Use:           "errno-gen",
		Short:         "errno-gen generate errno from enum",
		Long:          "errno-gen generate errno from enum",
		Version:       BuildVersion(),
		SilenceUsage:  false,
		SilenceErrors: false,
		Args:          cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			if !root.DisableStringer {
				tags := ""
				if len(root.Tags) > 0 {
					tags = strings.Join(root.Tags, ",")
				}
				wg := conc.WaitGroup{}
				for _, v := range root.Type {
					vv := v
					wg.Go(func() {
						args := make([]string, 0, 8)
						args = append(args, "-type", vv, "-linecomment")
						if tags != "" {
							args = append(args, "-tags", tags)
						}
						output, err := exec.Command("stringer", args...).CombinedOutput()
						if err != nil {
							slog.Error(strings.Join(args, " "), slog.String("error", string(output)))
						}
					})
				}
				wg.Wait()
			}

			srcDir := root.Pattern[0]
			fileInfo, err := os.Stat(srcDir)
			if err != nil {
				return err
			}
			if !fileInfo.IsDir() {
				if len(root.Tags) != 0 {
					slog.Error("--tags option applies only to directories, not when files are specified")
					os.Exit(1)
				}
				srcDir = filepath.Dir(srcDir)
			}
			g := &Gen{
				Pattern:   root.Pattern,
				OutputDir: srcDir,
				Type:      root.Type,
				Tags:      root.Tags,
				Version:   version,
				Epk:       root.Epk,
			}
			err = g.Generate()
			if err != nil {
				slog.Error("生成失败", slog.Any("err", err))
			}
			return nil
		},
	}
	cobra.OnInitialize(func() {
		textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   false,
			Level:       level(root.level),
			ReplaceAttr: nil,
		})
		slog.SetDefault(slog.New(textHandler))
	})

	cmd.PersistentFlags().StringVarP(&root.level, "level", "l", "info", "log level(debug,info,warn,error)")
	cmd.Flags().StringSliceVarP(&root.Pattern, "pattern", "p", []string{"."}, "the list of files or a directory.")
	cmd.Flags().StringSliceVarP(&root.Type, "type", "t", nil, "the list type of enum names; must be set")
	cmd.Flags().StringSliceVar(&root.Tags, "tags", nil, "comma-separated list of build tags to apply")
	cmd.Flags().BoolVarP(&root.DisableStringer, "disable-stringer", "d", false, "disable use `stringer` command.")
	cmd.Flags().StringVarP(&root.Epk, "epk", "e", "github.com/things-go/dyn/errorx", "errors package import path")

	root.cmd = cmd
	return root
}

// Execute adds all child commands to the root command and sets flags appropriately.
func (r *RootCmd) Execute() error {
	return r.cmd.Execute()
}

func level(s string) slog.Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
