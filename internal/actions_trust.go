package ticker

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	hlog "github.com/stellar/go/support/log"
	"github.com/xdefrag/ticker/internal/tickerdb"
)

type TrustResponse struct {
	Accounts []struct {
		Account string
		Assets  []struct {
			Asset string
		}
	}
}

func GenerateTrusts(
	ctx context.Context,
	s *tickerdb.TickerSession,
	l *hlog.Entry,
	apiURL string,
	priority int64,
) error {
	res, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	tres := TrustResponse{}
	if err := json.NewDecoder(res.Body).Decode(&tres); err != nil {
		return err
	}

	trs := make([]tickerdb.Trust, 0)

	for _, account := range tres.Accounts {
		for _, asset := range account.Assets {
			trs = append(trs, tickerdb.Trust{
				Code:          asset.Asset,
				IssuerAccount: account.Account,
				Source:        apiURL,
				Priority:      priority,
				UpdatedAt:     time.Now(),
			})
		}
	}

	return s.UpsertTrusts(ctx, trs)
}
