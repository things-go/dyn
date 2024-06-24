package command

import (
	"errors"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/things-go/dyn/cmd/dyngen/command/api"
)

type apiServiceOpt struct {
	source
	OutputDir   string            // M, 输出路径
	PackageName string            // M, 包名
	Options     map[string]string // M, proto option
	Filename    string            // o, filename
}

type apiServiceCmd struct {
	cmd *cobra.Command
	apiServiceOpt
}

func newApiServiceCmd() *apiServiceCmd {
	root := &apiServiceCmd{}
	cmd := &cobra.Command{
		Use:     "service",
		Short:   "Generate api service from database",
		Example: "dyngen api service",
		RunE: func(cmd *cobra.Command, args []string) error {
			schemaes, err := getSchema(&root.source)
			if err != nil {
				return err
			}
			if len(schemaes.Entities) == 0 {
				return errors.New("at least one schema entity")
			}
			entity := schemaes.Entities[0].IntoProto()
			filename := joinFilename(root.OutputDir, root.Filename, ".proto")
			_, err = os.Stat(filename)
			if err == nil || os.IsExist(err) {
				slog.Warn("🐛 '" + root.Filename + "' already exists, skipping")
				return nil
			}
			codegen := api.CodeGen{
				Entity:      entity,
				PackageName: root.PackageName,
				Options:     root.Options,
			}
			data := codegen.GenService().Bytes()
			err = WriteFile(filename, data)
			if err != nil {
				return err
			}
			slog.Info("👉 " + filename)
			return nil
		},
	}
	// input file
	cmd.Flags().StringSliceVarP(&root.InputFile, "input", "i", nil, "input file")
	cmd.Flags().StringVarP(&root.Schema, "schema", "s", "file+mysql", "parser file driver, [file+mysql,file+tidb](仅input时有效)")
	// database url
	cmd.Flags().StringVarP(&root.URL, "url", "u", "", "mysql://root:123456@127.0.0.1:3306/test")
	cmd.Flags().StringSliceVarP(&root.Tables, "table", "t", nil, "only out custom table")
	cmd.Flags().StringSliceVarP(&root.Exclude, "exclude", "e", nil, "exclude table pattern")

	cmd.Flags().StringVarP(&root.OutputDir, "out", "o", "./proto", "out directory")
	cmd.Flags().StringVar(&root.PackageName, "package", "", "proto package name")
	cmd.Flags().StringVar(&root.Filename, "filename", "service", "proto filename")
	cmd.Flags().StringToStringVar(&root.Options, "options", nil, "proto options")

	cmd.MarkFlagsOneRequired("url", "input")
	cmd.MarkFlagRequired("package")
	cmd.MarkFlagRequired("options")
	root.cmd = cmd
	return root
}
