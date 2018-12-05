package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/qcloud2018/go-demo/config"
	"github.com/qcloud2018/go-demo/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var conf config.Config

var RootCmd = &cobra.Command{
	Use:   "go-demo",
	Short: "demo application",
	Long:  "a pipeline demo web application",
}

func SetViperConfig(cfgFile string, name string, envPrefix string) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		viper.SetConfigName(name)
		viper.AddConfigPath(".")
		viper.AddConfigPath(fmt.Sprintf("%s/.%s/", home, name))
		viper.AddConfigPath(fmt.Sprintf("/etc/%s/", name))

		viper.AutomaticEnv()
		viper.SetEnvPrefix(envPrefix)
	}

	// don't use viper in centra instead of Config struct
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func main() {
	var cfgFile string
	cobra.OnInitialize(func() {
		var err = SetViperConfig(cfgFile, "go-demo", "GO_DEMO")
		if err != nil {
			panic(err)
		}

		if err := viper.Unmarshal(&conf); err != nil {
			panic(err)
		}

		if err := logger.InitZapLogger(conf.Env, conf.LogLevel, conf.LogFormat); err != nil {
			panic(err)
		}
	})

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "main config file ")

	if err := RootCmd.Execute(); err != nil {
		zap.L().Fatal("command exit error", zap.Error(err))
	}
}
