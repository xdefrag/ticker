package tickerdb

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (s *TickerSession) UpsertTrusts(ctx context.Context, trusts []Trust) error {
	if len(trusts) == 0 {
		return nil
	}

	qb := sq.Insert("trusts").
		Columns(getDBFieldTags(Trust{}, true)...)

	for _, trust := range trusts {
		qb = qb.Values(getDBFieldValues(trust, true)...)
	}

	qb = qb.Suffix(`ON CONFLICT (code, issuer_account) DO UPDATE SET
	priority = EXCLUDED.priority,
	source = EXCLUDED.source,
	updated_at = now()`)

	_, err := s.Exec(ctx, qb)

	return err
}

func (s *TickerSession) GetTrusts(ctx context.Context) ([]Trust, error) {
	var trusts []Trust
	err := s.Select(ctx, &trusts, sq.Expr(`SELECT * FROM trusts`))
	return trusts, err
}
