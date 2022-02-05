package cmd

import (
	"fmt"
	"os"

	"github.com/dilmnqvovpnmlib/gngc/gngc"
	"github.com/spf13/cobra"
)

var isNotified bool

var rootCmd = &cobra.Command{
	Use:   "gngc",
	Short: "Get GitHub contributions.",
	Long:  "Get GitHub contributions.",
	Run: func(cmd *cobra.Command, args []string) {
		msg := gngc.GetGitHubContributions()
		if isNotified {
			// go run main.go -n (IFTTT に通知を送信する。)
			gngc.NotifyIFTTT(msg)
		} else {
			// go run main.go
			fmt.Println(msg)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&isNotified, "notify", "n", false, "Get GitHub contributions and notify them to IFTTT.")

	gngc.LoadEnv()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
