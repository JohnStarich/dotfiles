package weather

import (
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/oschwald/maxminddb-golang"
)

func newGeoIPHandler(tb testing.TB, now time.Time) (net.IP, http.Handler) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tb.Fatal("Unexpected MaxMindDB call:", r.URL.Path)
	})
	testMaxMindDBFilePath := filepath.Join("testdata", "GeoIP2-City-Test.mmdb")
	dbPath := fmt.Sprintf("/free/dbip-city-lite-%04d-%02d.mmdb.gz", now.Year(), now.Month())
	mux.HandleFunc(dbPath, func(w http.ResponseWriter, r *http.Request) {
		dbFile, err := os.Open(testMaxMindDBFilePath)
		if err != nil {
			panic(err)
		}
		gzipWriter := gzip.NewWriter(w)
		_, err = io.Copy(gzipWriter, dbFile)
		if err != nil {
			panic(err)
		}
		err = gzipWriter.Close()
		if err != nil {
			panic(err)
		}
	})

	db, err := maxminddb.Open(testMaxMindDBFilePath)
	if err != nil {
		tb.Fatal(err)
	}
	defer db.Close()
	networks := db.Networks()
	for networks.Next() {
		var network ipCoordinates
		ip, err := networks.Network(&network)
		if err != nil {
			tb.Fatal(err)
		}
		if network != (ipCoordinates{}) {
			return ip.IP, mux
		}
	}
	tb.Fatal("no valid IP addresses with locations in test MaxMindDB database file")
	return nil, nil
}
