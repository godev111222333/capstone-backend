package service

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/godev111222333/capstone-backend/src/misc"
	"github.com/godev111222333/capstone-backend/src/store"
)

const (
	ShouldResetDatabase = true
)

var (
	TestDb       *store.DbStore
	TestS3Store  *store.S3Store
	TestConfig   *misc.GlobalConfig
	TestFeConfig *misc.FEConfig
)

func TestMain(m *testing.M) {
	cfg, err := misc.LoadConfig("../../config.yaml")
	if err != nil {
		panic(err)
	}

	feCfg, err := misc.LoadFEConfig("../../fe-config.yaml")
	if err != nil {
		panic(err)
	}

	TestFeConfig = feCfg

	TestConfig = cfg
	TestS3Store = store.NewS3Store(cfg.AWS)
	dbConfig := cfg.Database
	initTestDb(dbConfig)
	code := m.Run()
	os.Exit(code)
}

func initTestDb(cfg *misc.DatabaseConfig) {
	if ShouldResetDatabase {
		if err := ResetDb(cfg); err != nil {
			panic(err)
		}
	}

	var err error
	TestDb, err = store.NewDbStore(cfg)
	if err != nil {
		panic(err)
	}
}

func ResetDb(cfg *misc.DatabaseConfig) error {
	dbString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DbUsername, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName,
	)
	downCmd := exec.Command("migrate", "-path", "../../migration", "-database", dbString, "-verbose", "down")
	downCmd.Stdout = os.Stdout
	downCmd.Stderr = os.Stderr
	downStdIn, err := downCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	if err := downCmd.Start(); err != nil {
		return err
	}

	// send "y" cmd
	if _, err := io.WriteString(downStdIn, "y\n"); err != nil {
		return err
	}

	if err := downCmd.Wait(); err != nil {
		return err
	}

	upCmd := exec.Command("migrate", "-path", "../../migration", "-database", dbString, "-verbose", "up")
	if err := upCmd.Run(); err != nil {
		return err
	}

	return nil
}
