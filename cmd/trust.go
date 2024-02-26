package cmd

import (
	"context"

	"github.com/lib/pq"
	"github.com/spf13/cobra"
	ticker "github.com/xdefrag/ticker/internal"
	"github.com/xdefrag/ticker/internal/tickerdb"
)

var (
	TrustAPIURL   string
	TrustPriority int64
)

func init() {
	rootCmd.AddCommand(cmdTrust)

	cmdTrust.Flags().StringVarP(
		&TrustAPIURL,
		"api-url",
		"u",
		"https://eurmtl.me/remote/good_assets",
		"Set API URL for trusted assets",
	)

	cmdTrust.Flags().Int64VarP(
		&TrustPriority,
		"priority",
		"p",
		0,
		"Set priority of trust source",
	)
}

var cmdTrust = &cobra.Command{
	Use:   "trust [url]",
	Short: "Generate pool of trusted assets from api url with priority",
	Run: func(cmd *cobra.Command, args []string) {
		dbInfo, err := pq.ParseURL(DatabaseURL)
		if err != nil {
			Logger.Fatal("could not parse db-url:", err)
		}

		session, err := tickerdb.CreateSession("postgres", dbInfo)
		if err != nil {
			Logger.Fatal("could not connect to db:", err)
		}
		defer session.DB.Close()

		ctx := context.Background()
		if err := ticker.GenerateTrusts(ctx, &session, Logger, TrustAPIURL, TrustPriority); err != nil {
			Logger.Fatal("could not generate trusts:", err)
		}
	},
}
