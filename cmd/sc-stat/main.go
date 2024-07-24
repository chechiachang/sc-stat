package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chechiachang/sc-stat/pkg/fetcher"
	"github.com/chechiachang/sc-stat/pkg/github"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "datafetcher",
	Short: "data fetcher cronjob runner",
	Long:  "fetch public data from various sources",

	SilenceUsage: false,

	RunE: func(cmd *cobra.Command, args []string) error {
		return run(context.Background())
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "log level")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(baseCtx context.Context) error {
	ctx, cancel := context.WithCancel(baseCtx)
	defer cancel()
	return runServer(ctx)
}

func runServer(ctx context.Context) error {
	log.Info("Starting server")
	cronjob := cron.New()

	// data fetcher
	cronjob.AddFunc("@every 1m", fetcher.Yilan)
	// data commit
	date := time.Now()
	commitService := github.CommitService{
		commitMessage: fmt.Sprintf("chore: upload data on %s", date.Format("2006-1-2")),
		sourceOwner:   "chechiachang",
		sourceRepo:    "sc-stat-data",
		commitBranch:  "main",
		baseBranch:    "main",
		authorName:    "sc-stat-automation",
		authorEmail:   "chechiachang999@gmail.com",
		privateKey:    "",
	}
	cronjob.AddFunc("@every 10m", commitService.CommitPush)

	cronjob.Start()
	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping server")
			return nil
		}
	}
	defer cronjob.Stop()

	return nil
}
