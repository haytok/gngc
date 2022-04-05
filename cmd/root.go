package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/haytok/gngc/gngc"
	"github.com/haytok/gngc/model"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var config model.Config

var isNotified bool

var rootCmd = &cobra.Command{
	Use:           "gngc",
	Short:         "Get and Notify GitHub Contributions",
	Long:          "A simple command line application to get GitHub contributions from GraphQL API and notify them to IFTTT.",
	SilenceErrors: true,
	// SilenceUsage:  true,

	RunE: runCommand,
}

func runCommand(cmd *cobra.Command, args []string) error {
	if config.GitHub.UserName == "" && config.GitHub.Token == "" {
		return errors.New("Credentials for GitHub dont properly configured in $HOME/.gngc.toml.")
	}

	msg, err := gngc.GetGitHubContributions(config.GitHub)
	if err != nil {
		return err
	}

	// go run main.go -n
	// IFTTT にメッセージを送信した結果のレスポンスのメッセージが作成される。
	if isNotified {
		if config.IFTTT.EventName == "" && config.IFTTT.Token == "" {
			return errors.New("Credentials for IFTTT dont properly configured in $HOME/.gngc.toml.")
		}

		msg, err = gngc.NotifyIFTTT(config.IFTTT, msg)
		if err != nil {
			return err
		}
	}

	fmt.Println(msg)

	return nil
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// --config を付けない場合に呼び出される処理
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error", err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".gngc")
		viper.SetConfigType("toml")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Error", err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Fprintln(os.Stderr, "Error", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gngc.toml)")
	rootCmd.PersistentFlags().BoolVarP(&isNotified, "notify", "n", false, "Get GitHub contributions and notify them to IFTTT.")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
