package web

import (
	"net/http"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/victorjacobs/csv2ynab/config"
	"github.com/victorjacobs/csv2ynab/importer"
	"github.com/victorjacobs/csv2ynab/ynab"
	"golang.org/x/net/webdav"
)

// Serve serves a webdav server
func Serve(listenAddr string, cfg config.Config) {
	fs := NewFS(func(path string, contents []byte) {
		if strings.HasPrefix(path, "._") {
			return
		}

		if len(contents) == 0 {
			return
		}

		log.Printf("Processing: %v\n", path)

		ynabConfig, err := cfg.YNABConfigForFile(path)
		if err != nil {
			log.Errorf("Failed to get YNAB config: %v", err)
			return
		}

		if transactions, err := importer.Parse(filepath.Base(path), contents); err != nil {
			log.Errorf("Failed to parse transactions: %v", err)
		} else if err := ynab.Write(ynabConfig, transactions); err != nil {
			log.Errorf("Failed to import transactions: %v", err)
		}
	})

	handler := &webdav.Handler{
		Prefix:     "/", // Root path for the WebDAV endpoint
		FileSystem: fs,
		LockSystem: webdav.NewMemLS(),
	}

	// Start the server
	log.Printf("WebDAV server is running on %v\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, handler))
}
