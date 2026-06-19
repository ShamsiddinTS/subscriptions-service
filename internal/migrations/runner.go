package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Migration struct {
	Version int
	Name    string
	Path    string
	SQL     string
}

var reUp = regexp.MustCompile(`^(\d+)_.*\.up\.sql$`)

func ResolveMigrationsDir(dir string) (string, error) {
	if filepath.IsAbs(dir) {
		if st, err := os.Stat(dir); err == nil && st.IsDir() {
			return dir, nil
		}
		return "", fmt.Errorf("migrations dir %s is not a directory", dir)
	}

	if wd, err := os.Getwd(); err == nil {
		p := filepath.Join(wd, dir)
		if st, err := os.Stat(p); err == nil && st.IsDir() {
			return p, nil
		}
	}

	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	base := filepath.Dir(exe)
	p := filepath.Join(base, dir)

	if st, err := os.Stat(p); err == nil && st.IsDir() {
		return p, nil
	}

	return "", fmt.Errorf("migrations dir %s is not a directory", dir)
}

func Run(ctx context.Context, db *sql.DB, dir string) error {
	if err := ensureMigrationsTable(ctx, db); err != nil {
		return err
	}

	migs, err := loadUpMigrations(dir)
	if err != nil {
		return err
	}

	if len(migs) == 0 {
		return fmt.Errorf("no .up.sql migrations found in %s", dir)
	}

	applied, err := appliedVersions(ctx, db)
	if err != nil {
		return err
	}

	for _, mig := range migs {
		if applied[mig.Version] {
			continue
		}

		if err := applyOne(ctx, db, mig); err != nil {
			return fmt.Errorf("apply migration %d failed: %w", mig.Version, err)
		}
	}

	return nil
}

func ensureMigrationsTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    applied_at 	TIMESTAMPTZ NOT NULL DEFAULT now()
);`)
	return err
}

func loadUpMigrations(dir string) ([]Migration, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migs []Migration

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		matches := reUp.FindStringSubmatch(name)
		if matches == nil {
			continue
		}

		version, _ := strconv.Atoi(matches[1])

		path := filepath.Join(dir, name)

		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		migs = append(migs, Migration{
			Version: version,
			Name:    name,
			Path:    path,
			SQL:     string(b),
		})
	}

	sort.Slice(migs, func(i, j int) bool {
		return migs[i].Version < migs[j].Version
	})

	return migs, nil
}

func appliedVersions(ctx context.Context, db *sql.DB) (map[int]bool, error) {
	rows, err := db.QueryContext(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	return applied, rows.Err()
}

func applyOne(ctx context.Context, db *sql.DB, mig Migration) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	sqlText := strings.TrimSpace(mig.SQL)
	if sqlText == "" {
		return fmt.Errorf("migration file is empty: %s", mig.Name)
	}

	if _, err := tx.ExecContext(ctx, sqlText); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO schema_migrations(version, name, applied_at) VALUES($1,$2,$3)`,
		mig.Version, mig.Name, time.Now(),
	); err != nil {
		return err
	}

	return tx.Commit()
}
