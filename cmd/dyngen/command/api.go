package command

import (
	"errors"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/things-go/dyn/cmd/dyngen/command/api"
)

type apiOpt struct {
	source
	OutputDir                 string            // M, 输出路径
	PackageName               string            // M, 包名
	Options                   map[string]string // required, proto option
	Style                     string            // 字段代码风格, snakeCase, smallCamelCase, pascalCase
	DisableBool               bool              // 禁用bool,使用int32
	DisableTimestamp          bool              // 禁用google.protobuf.Timestamp,使用int64
	EnableOpenapiv2Annotation bool              // 启用int64的openapiv2注解
}

type apiCmd struct {
	cmd *cobra.Command
	apiOpt
}

func newApiCmd() *apiCmd {
	root := &apiCmd{}
	cmd := &cobra.Command{
		Use:     "api",
		Short:   "Generate api from database",
		Example: "dyngen api",
		RunE: func(cmd *cobra.Command, args []string) error {
			schemaes, err := getSchema(&root.source)
			if err != nil {
				return err
			}
			if len(schemaes.Entities) == 0 {
				return errors.New("at least one schema entity")
			}
			entity := schemaes.Entities[0].IntoProto()
			filename := joinFilename(root.OutputDir, entity.TableName, ".proto")
			_, err = os.Stat(filename)
			if err == nil || os.IsExist(err) {
				slog.Warn("🐛 '" + entity.TableName + "' already exists, skipping")
				return nil
			}
			codegen := api.CodeGen{
				Entity:                    entity,
				PackageName:               root.PackageName,
				Options:                   root.Options,
				Style:                     root.Style,
				DisableBool:               root.DisableBool,
				DisableTimestamp:          root.DisableTimestamp,
				EnableOpenapiv2Annotation: root.EnableOpenapiv2Annotation,
			}
			data := codegen.Gen().Bytes()
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
	cmd.Flags().StringToStringVar(&root.Options, "options", nil, "proto options")
	cmd.Flags().StringVar(&root.Style, "style", "smallCamelCase", "字段代码风格, snakeCase, smallCamelCase, pascalCase")

	cmd.Flags().BoolVar(&root.DisableBool, "disableBool", false, "禁用bool,使用int32")
	cmd.Flags().BoolVar(&root.DisableTimestamp, "disableTimestamp", false, "禁用google.protobuf.Timestamp,使用int64")
	cmd.Flags().BoolVar(&root.EnableOpenapiv2Annotation, "EnableOpenapiv2Annotation", false, "启用int64的openapiv2注解")

	cmd.MarkFlagsOneRequired("url", "input")
	cmd.MarkFlagRequired("package")
	cmd.MarkFlagRequired("options")
	root.cmd = cmd
	return root
}
