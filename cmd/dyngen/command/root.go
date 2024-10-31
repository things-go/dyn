package command

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/things-go/dyn/cmd/internal/meta"
	_ "github.com/things-go/ens/driver/mysql"
)

type RootCmd struct {
	cmd   *cobra.Command
	level string
}

func NewRootCmd() *RootCmd {
	root := &RootCmd{}
	cmd := &cobra.Command{
		Use:           "ormat",
		Short:         "gorm reflect tools",
		Long:          "database/sql to golang struct",
		Version:       meta.BuildVersion(),
		SilenceUsage:  false,
		SilenceErrors: false,
		Args:          cobra.NoArgs,
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
	cmd.AddCommand(
		newDalCmd().cmd,
		newApiCmd().cmd,
	)
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
