package tests

import (
	"flag"
	"os"
	"testing"

	"github.com/OCCASS/avito-intern/internal/config"
	"github.com/OCCASS/avito-intern/internal/database"
)

var (
	cfg *config.Config
	db  *database.Database
)

func TestMain(m *testing.M) {
	var cfgPath string

	flag.StringVar(&cfgPath, "config", "", "Configuration file path.")
	flag.Parse()

	cfg = config.MustLoad(cfgPath)
	db = database.MustConnect(cfg.Database)

	os.Exit(m.Run())
}
