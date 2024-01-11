package tickerdb

import (
	"net/http"
	"os"
	"strings"
	"testing"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/shurcooL/httpfs/filter"

	supportHttp "github.com/stellar/go/support/http"
	bdata "github.com/xdefrag/ticker/internal/tickerdb/migrations"
)

func TestGeneratedAssets(t *testing.T) {
	var localAssets http.FileSystem = filter.Keep(http.Dir("migrations"), func(path string, fi os.FileInfo) bool {
		return fi.IsDir() || strings.HasSuffix(path, ".sql")
	})
	generatedAssets := &assetfs.AssetFS{
		Asset:     bdata.Asset,
		AssetDir:  bdata.AssetDir,
		AssetInfo: bdata.AssetInfo,
		Prefix:    "/migrations",
	}

	if !supportHttp.EqualFileSystems(localAssets, generatedAssets, "/") {
		t.Fatalf("generated migrations does not match local migrations")
	}
}
