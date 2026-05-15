package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"

	"golang-train/internal/shared/config"
)

var tenantRe = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) < 3 {
		usage()
		os.Exit(2)
	}

	cmd := os.Args[1]    // common|tenant
	action := os.Args[2] // up|down|status|reset|down-to
	args := os.Args[3:]

	cfg := config.LoadFromEnv()

	migrationsRoot := getenv("MIGRATIONS_DIR", "migrations")
	migrationsRoot, _ = filepath.Abs(migrationsRoot)

	goose.SetDialect("mysql")

	switch cmd {
	case "common":
		publicSchema := getenv("DB_PUBLIC_SCHEMA", "public")
		ensureDatabase(cfg, publicSchema)
		db := mustOpenDB(cfg, publicSchema)
		defer db.Close()
		mustRun(action, db, filepath.Join(migrationsRoot, "common"), args)

	case "tenant":
		fs := flag.NewFlagSet("tenant", flag.ExitOnError)
		tenantID := fs.String("tenant", "", "tenant id (used to build schema name)")
		all := fs.Bool("all", false, "migrate all tenants listed in public.companies")
		_ = fs.Parse(args)

		prefix := getenv("DB_TENANT_SCHEMA_PREFIX", "tenant_")
		publicSchema := getenv("DB_PUBLIC_SCHEMA", "public")

		if *all {
			migrateAllTenants(cfg, migrationsRoot, publicSchema, prefix, action, fs.Args())
			return
		}

		if *tenantID == "" {
			log.Fatal("missing --tenant (or use --all)")
		}
		if !tenantRe.MatchString(*tenantID) {
			log.Fatalf("invalid tenant id: %q", *tenantID)
		}

		schema := prefix + *tenantID
		ensureDatabase(cfg, schema)

		db := mustOpenDB(cfg, schema)
		defer db.Close()
		mustRun(action, db, filepath.Join(migrationsRoot, "tenant"), fs.Args())

	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage:
  migrate common <up|down|status|reset|down-to> [args...]
  migrate tenant <up|down|status|reset|down-to> --tenant=<id> [args...]
  migrate tenant <up|down|status|reset|down-to> --all [args...]

Env:
  DB_HOST, DB_PORT, DB_USER, DB_PASSWORD
  DB_PUBLIC_SCHEMA (default: public)
  DB_TENANT_SCHEMA_PREFIX (default: tenant_)
  MIGRATIONS_DIR (default: migrations)
`)
}

func mustRun(action string, db *sql.DB, dir string, args []string) {
	if err := runMigration(action, db, dir, args); err != nil {
		log.Fatal(err)
	}
}

func runMigration(action string, db *sql.DB, dir string, args []string) error {
	switch action {
	case "up":
		return goose.Up(db, dir)
	case "down":
		return goose.Down(db, dir)
	case "status":
		return goose.Status(db, dir)
	case "reset":
		return goose.Reset(db, dir)
	case "down-to":
		if len(args) < 1 {
			return errors.New("down-to requires a version argument")
		}
		v, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid version argument: %w", err)
		}
		return goose.DownTo(db, dir, v)
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}

func migrateAllTenants(cfg config.Config, migrationsRoot, publicSchema, prefix, action string, args []string) {
	ensureDatabase(cfg, publicSchema)
	publicDB := mustOpenDB(cfg, publicSchema)
	defer publicDB.Close()

	rows, err := publicDB.Query("SELECT id FROM companies")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tenantID string
		if err := rows.Scan(&tenantID); err != nil {
			log.Fatal(err)
		}
		if !tenantRe.MatchString(tenantID) {
			log.Printf("skip invalid tenant id from companies: %q", tenantID)
			continue
		}

		schema := prefix + tenantID
		ensureDatabase(cfg, schema)

		db := mustOpenDB(cfg, schema)
		err := func() error {
			defer db.Close()
			return runMigration(action, db, filepath.Join(migrationsRoot, "tenant"), args)
		}()
		if err != nil {
			log.Fatalf("tenant %s: %v", tenantID, err)
		}
		log.Printf("tenant %s migrated", tenantID)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func mustOpenDB(cfg config.Config, schema string) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local&multiStatements=true",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		schema,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func ensureDatabase(cfg config.Config, schema string) {
	bootstrap := mustOpenDB(cfg, "mysql")
	defer bootstrap.Close()

	if _, err := bootstrap.Exec("CREATE DATABASE IF NOT EXISTS `" + schema + "` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"); err != nil {
		log.Fatal(err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
