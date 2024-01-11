package ticker

import (
	hlog "github.com/stellar/go/support/log"
	"github.com/xdefrag/ticker/internal/gql"
	"github.com/xdefrag/ticker/internal/tickerdb"
)

func StartGraphQLServer(s *tickerdb.TickerSession, l *hlog.Entry, port string) {
	graphql := gql.New(s, l)

	graphql.Serve(port)
}
