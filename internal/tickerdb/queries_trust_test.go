//go:build integration

package tickerdb

import (
	"context"
	"testing"
	"time"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testTrusts = []Trust{
	{
		Code:          "STAS",
		IssuerAccount: "GCUEVVS4KIHZM72DAHKXIRWSCN3V3Y4KX6UNNUU7PV73VQK44CNKAMNI",
		Source:        "MTL fund",
		Priority:      1,
		UpdatedAt:     time.Now(),
	},
	{
		Code:          "MTL",
		IssuerAccount: "GACKTN5DAZGWXRWB2WLM6OPBDHAMT6SJNGLJZPQMEZBUR4JUGBX2UK7V",
		Source:        "MTL fund",
		Priority:      0,
		UpdatedAt:     time.Now(),
	},
}

func TestUpdateTrusts(t *testing.T) {
	db := OpenTestDBConnection(t)
	defer db.Close()

	var session TickerSession
	session.DB = db.Open()
	ctx := context.Background()
	defer session.DB.Close()

	migrations := &migrate.FileMigrationSource{
		Dir: "./migrations",
	}
	_, err := migrate.Exec(session.DB.DB, "postgres", migrations, migrate.Up)
	require.NoError(t, err)

	err = session.UpsertTrusts(ctx, testTrusts)
	require.NoError(t, err)

	var destTrust1 Trust
	err = session.GetRaw(ctx, &destTrust1, `SELECT * FROM trusts WHERE code='STAS'`)
	require.NoError(t, err)

	assert.Equal(t, testTrusts[0].IssuerAccount, destTrust1.IssuerAccount)
	assert.Equal(t, testTrusts[0].Source, destTrust1.Source)

	var destTrust2 Trust
	err = session.GetRaw(ctx, &destTrust2, `SELECT * FROM trusts WHERE code='MTL'`)
	require.NoError(t, err)

	assert.Equal(t, testTrusts[1].IssuerAccount, destTrust2.IssuerAccount)
	assert.Equal(t, testTrusts[1].Source, destTrust2.Source)

	const newSource = "MTL Association"

	err = session.UpsertTrusts(ctx, []Trust{
		{
			Code:          "STAS",
			IssuerAccount: "GCUEVVS4KIHZM72DAHKXIRWSCN3V3Y4KX6UNNUU7PV73VQK44CNKAMNI",
			Source:        newSource,
			Priority:      1,
			UpdatedAt:     time.Now(),
		},
	})
	require.NoError(t, err)

	var destTrust3 Trust
	err = session.GetRaw(ctx, &destTrust3, `SELECT * FROM trusts WHERE code='STAS'`)
	require.NoError(t, err)

	assert.Equal(t, newSource, destTrust3.Source)
}

func TestGetTrusts(t *testing.T) {
	db := OpenTestDBConnection(t)
	defer db.Close()

	var session TickerSession
	session.DB = db.Open()
	ctx := context.Background()
	defer session.DB.Close()

	migrations := &migrate.FileMigrationSource{
		Dir: "./migrations",
	}
	_, err := migrate.Exec(session.DB.DB, "postgres", migrations, migrate.Up)
	require.NoError(t, err)

	err = session.UpsertTrusts(ctx, testTrusts)
	require.NoError(t, err)

	trusts, err := session.GetTrusts(ctx)
	require.NoError(t, err)

	assert.Len(t, trusts, 2)
}
