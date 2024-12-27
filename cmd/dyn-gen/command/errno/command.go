package errno

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sourcegraph/conc"
	"github.com/spf13/cobra"
	"github.com/things-go/dyn/cmd/internal/meta"
)

type ErrnoOption struct {
	Pattern         []string
	Type            []string
	Tags            []string
	DisableStringer bool
	Epk             string
}

type ErrnoCmd struct {
	Cmd *cobra.Command
	ErrnoOption
}

func NewErrnoCmd() *ErrnoCmd {
	root := &ErrnoCmd{}
	cmd := &cobra.Command{
		Use:           "errno",
		Short:         "generate errno from enum",
		Long:          "generate errno from enum",
		Version:       meta.BuildVersion(),
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
				Version:   meta.Version,
				Epk:       root.Epk,
			}
			err = g.Generate()
			if err != nil {
				slog.Error("生成失败", slog.Any("err", err))
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&root.Pattern, "pattern", "p", []string{"."}, "the list of files or a directory.")
	cmd.Flags().StringSliceVarP(&root.Type, "type", "t", nil, "the list type of enum names; must be set")
	cmd.Flags().StringSliceVar(&root.Tags, "tags", nil, "comma-separated list of build tags to apply")
	cmd.Flags().BoolVarP(&root.DisableStringer, "disable-stringer", "d", false, "disable use `stringer` command.")
	cmd.Flags().StringVarP(&root.Epk, "epk", "e", "github.com/things-go/dyn/errorx", "errors package import path")

	root.Cmd = cmd
	return root
}
