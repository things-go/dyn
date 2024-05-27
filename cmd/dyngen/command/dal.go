package command

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/things-go/ens"
	"github.com/things-go/ens/utils"
)

type dalOpt struct {
	source
	OutputDir       string
	PackageName     string // 包名
	ModelImportPath string // required, model导入路径
	RepoImportPath  string // required, repository导入路径
	DalImportPath   string // required, dal导入路径
	CustomTemplate  string // 自定义模板
	ens.Option
}

type dalCmd struct {
	cmd *cobra.Command
	dalOpt
}

func newDakCmd() *dalCmd {
	root := &dalCmd{}

	cmd := &cobra.Command{
		Use:     "dal",
		Short:   "Generate dal from database",
		Example: "ormat dal",
		RunE: func(*cobra.Command, []string) error {
			if root.CustomTemplate == "builtin-rapier" && root.RepoImportPath == "" {
				return errors.New("使用builtin-rapier时repository导入路径, 不能为空")
			}
			schemaes, err := getSchema(&root.source)
			if err != nil {
				return err
			}
			daltpl, err := GetUsedTemplate(root.CustomTemplate)
			if err != nil {
				return err
			}
			buf := bytes.Buffer{}
			packageName := cmp.Or(root.PackageName, utils.GetPkgName(root.OutputDir))
			queryImportPath := strings.Join([]string{root.DalImportPath, "query"}, "/")
			dalOptionFilename := joinFilename(root.OutputDir, "dal_option", ".go")
			_, err = os.Stat(dalOptionFilename)
			if !(err == nil || os.IsExist(err)) {
				err = dalOptionTpl.Execute(&buf, Dal{Package: packageName})
				if err != nil {
					return err
				}
				err = WriteFile(dalOptionFilename, buf.Bytes())
				if err != nil {
					return fmt.Errorf("dal_option: %v", err)
				}
				slog.Info("👉 " + dalOptionFilename)
			} else {
				slog.Warn("🐛 dal_option.go already exists, skipping")
			}
			dal := Dal{
				Package:     packageName,
				Imports:     []string{root.ModelImportPath, queryImportPath, root.RepoImportPath},
				ModelPrefix: utils.PkgName(root.ModelImportPath) + ".",
				QueryPrefix: "query.",
				RepoPrefix:  utils.PkgName(root.RepoImportPath) + ".",
				Entity:      nil,
			}
			dalQuery := Dal{
				Package:     "query",
				Imports:     []string{},
				ModelPrefix: utils.PkgName(root.ModelImportPath) + ".",
				QueryPrefix: "",
				RepoPrefix:  "",
				Entity:      nil,
			}

			for _, entity := range schemaes.Entities {
				dalFilename := joinFilename(root.OutputDir, entity.Name, ".go")
				_, err = os.Stat(dalFilename)
				if err == nil || os.IsExist(err) {
					slog.Warn("🐛 " + entity.Name + " already exists, skipping")
					continue
				}
				dal.Entity = entity
				buf.Reset()
				err = daltpl.Execute(&buf, dal)
				if err != nil {
					return err
				}

				err = WriteFile(dalFilename, buf.Bytes())
				if err != nil {
					return fmt.Errorf("%v: %v", entity.Name, err)
				}

				buf.Reset()
				dalQuery.Entity = entity
				err = dalQueryTpl.Execute(&buf, dalQuery)
				if err != nil {
					return err
				}
				dalQueryFilename := joinFilename(filepath.Join(root.OutputDir, "query"), entity.Name, ".go")
				err = WriteFile(dalQueryFilename, buf.Bytes())
				if err != nil {
					return err
				}
				slog.Info("👉 " + dalFilename)
				slog.Info("👉 " + dalQueryFilename)
			}

			slog.Info("😄 generate success !!!")
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

	cmd.Flags().StringVarP(&root.OutputDir, "out", "o", "./dal", "out directory")
	cmd.Flags().StringVar(&root.PackageName, "package", "", "package name")
	cmd.Flags().StringVar(&root.CustomTemplate, "template", "builtin-rapier", "use custom template except [builtin-rapier, builtin-gorm]")
	cmd.Flags().StringVar(&root.ModelImportPath, "modelImportPath", "", "model导入路径")
	cmd.Flags().StringVar(&root.DalImportPath, "dalImportPath", "", "dal导入路径")
	cmd.Flags().StringVar(&root.RepoImportPath, "repoImportPath", "", "repository导入路径")

	cmd.Flags().BoolVar(&root.EnableInt, "enableInt", false, "使能int8,uint8,int16,uint16,int32,uint32输出为int,uint")
	cmd.Flags().BoolVar(&root.EnableBoolInt, "enableBoolInt", false, "使能bool输出int")
	cmd.Flags().BoolVar(&root.DisableNullToPoint, "disableNullToPoint", false, "禁用字段为null时输出指针类型,将输出为sql.Nullxx")
	cmd.Flags().StringSliceVar(&root.EscapeName, "escapeName", nil, "escape name list")

	cmd.MarkFlagsOneRequired("url", "input")
	cmd.MarkFlagRequired("modelImportPath")
	cmd.MarkFlagRequired("dalImportPath")
	root.cmd = cmd
	return root
}
