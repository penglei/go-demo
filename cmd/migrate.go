package main

import (
	"fmt"
	_ "github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/qcloud2018/go-demo/migration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var versionsDir string

const (
	ACTION_UP     = "up"
	ACTION_DOWN   = "down"
	ACTION_NEW    = "new"
	MAX_TITLE_LEN = 40
)

var migrateCmd = &cobra.Command{
	Use:       "migrate [<action>(up|down|new)] [flags]\n\nDefault Args: \n ",
	Short:     "Update database",
	Long:      `Update database with version files.`,
	ValidArgs: []string{ACTION_UP, ACTION_DOWN, ACTION_NEW},
	Args:      cobra.OnlyValidArgs,

	Run: func(cmd *cobra.Command, args []string) {
		var action string
		if len(args) > 0 {
			action = args[0]
		} else {
			action = ""
		}
		runMigrate(action, cmd)
	},
}

func runMigrate(action string, cmd *cobra.Command) {
	var l = zap.L()
	var m = migration.NewFileMigration(conf.Database)
	var absPath, err = filepath.Abs(versionsDir)
	if err != nil {
		l.Error("get version files directory failed",
			zap.String("directory", versionsDir), zap.String("path", absPath))
		return
	}

	m.Versions = absPath

	if action == ACTION_UP {
		// m.Upgrade(uint(steps))
		m.Upgrade()
	} else if action == ACTION_DOWN {
		// m.Downgrade(uint(steps))
		m.Downgrade()
	} else if action == ACTION_NEW {
		var now = time.Now()
		var timeStr = now.Format("20060102150405")
		title, err := cmd.Flags().GetString("title")
		var shortTitle string
		if err != nil {
			l.Error("error while parsing params title", zap.Error(err))
		}
		if len(title) > MAX_TITLE_LEN {
			shortTitle = title[:MAX_TITLE_LEN]
		} else {
			shortTitle = title
		}
		shortTitle = strings.ToLower(strings.Replace(shortTitle, " ", "-", -1))
		var up = fmt.Sprintf("%s/%s_%s.up.sql", absPath, timeStr, shortTitle)
		var down = fmt.Sprintf("%s/%s_%s.down.sql", absPath, timeStr, shortTitle)
		upFile, err := os.Create(up)
		if err != nil {
			l.Error("Error while creating migration up file", zap.String("file", up), zap.Error(err))
			return
		}
		downFile, err := os.Create(down)
		if err != nil {
			l.Error("Error while creating migration down file", zap.String("file", down), zap.Error(err))
			return
		}
		defer upFile.Close()
		defer downFile.Close()
		upFile.WriteString("-- " + title + " upgrade script")
		downFile.WriteString("-- " + title + " downgrade script")
		fmt.Printf("Migration file generated:\n  up: %s\n  down: %s\n", path.Base(up), path.Base(down))
	} else {
		var forceDown, err = cmd.Flags().GetBool("force-down")
		if err != nil {
			l.Error("Get flag failed", zap.Error(err))
		} else {
			if forceDown {
				m.ForceResetDown()
				return
			}
		}
		cmd.Help()
	}

}

func init() {
	RootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().StringVarP(&versionsDir, "versions", "v", "./migration/versions", "version files directory")
	// migrateCmd.Flags().IntVarP(&steps, "step", "s", 2, "steps")
	migrateCmd.Flags().StringP("host", "H", "", "database host address")
	migrateCmd.Flags().IntP("port", "P", 0, "database host port")
	migrateCmd.Flags().StringP("user", "u", "", "database user")
	migrateCmd.Flags().StringP("password", "p", "", "database user password")
	migrateCmd.Flags().StringP("name", "n", "", "database name")
	migrateCmd.Flags().Bool("force-down", false, "reset dirty state and downgrade the current version")
	migrateCmd.Flags().StringP("title", "t", "", "migration title")

	// 可以通过命令行覆盖参数
	viper.BindPFlag("database.host", migrateCmd.Flags().Lookup("host"))
	viper.BindPFlag("database.port", migrateCmd.Flags().Lookup("port"))
	viper.BindPFlag("database.user", migrateCmd.Flags().Lookup("user"))
	viper.BindPFlag("database.password", migrateCmd.Flags().Lookup("password"))
	viper.BindPFlag("database.database", migrateCmd.Flags().Lookup("name"))
}
