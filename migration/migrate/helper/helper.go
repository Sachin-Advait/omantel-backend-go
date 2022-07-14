package helper

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils"
)

const (
	createMigrationTableQuery = "CREATE TABLE IF NOT EXISTS migrations (version character varying(64) PRIMARY KEY, updated_at TIME DEFAULT current_time NOT NULL);"
	migrationInsertQuery      = "insert into migrations(version, updated_at) VALUES($1, $2) RETURNING version, updated_at"
	migrationReadQuery        = "select version, updated_at from migrations order by updated_at asc"
)

type MigrationRecords struct {
	Version   string
	UpdatedAt time.Time
}

//startMigration runs all pending migrations
func StartMigration(dbConn *sql.DB, migrationPath string, migrationFiles ...string) {

	_, err := createMigrationTable(dbConn)
	if err != nil {
		utils.LogObj.Fatal(err.Error(), "Migration: failed to create migration table")
	}

	migrateDataBase(dbConn, migrationPath, migrationFiles...)
}

// migrateDB finds the last run migration, and run all those after it in order
// We use the migration table to do this
func migrateDataBase(dbConn *sql.DB, migrationPath string, migrationFiles ...string) {

	files := make([]string, 0)

	switch {
	case len(migrationFiles) > 0:
		files = migrationFiles
		break
	default:
		// Get a list of migration files
		migrationFiles, err := filepath.Glob(migrationPath)
		if err != nil {
			utils.LogObj.Fatal("path: "+migrationPath, err.Error(), "Migration: failed to get all files files")
		}

		files = migrationFiles
	}

	// Sort the list alphabetically
	sort.Slice(files, func(i, j int) bool {
		return stringToNum(files[i]) < stringToNum(files[j])
	})

	existingFiles, err := readMetadata(dbConn)
	if err != nil {
		utils.LogObj.Fatal(err.Error(), "path: "+migrationPath, "Migration: read metadata failed")
	}

	for _, file := range files {

		fileName := filepath.Base(file)

		if contains(fileName, existingFiles) {
			utils.LogObj.Info("fileName: "+fileName, "Migration: file already migrated")
			continue
		}

		fileData, err := ioutil.ReadFile(file)
		if err != nil {
			utils.LogObj.Fatal("file: "+file, "err: "+err.Error(), "Migration: failed to read file content")
		}

		trx, err := dbConn.BeginTx(context.Background(), nil)
		if err != nil {
			utils.LogObj.Fatal("file: "+file, "err: "+err.Error(), "Migration: failed to begin transaction")
		}

		_, err = trx.Exec(string(fileData))
		if err != nil {
			trx.Rollback()
			utils.LogObj.Error("file: "+file, "err: "+err.Error(), "Migration: failed to execute file content")
			return
		}

		rslt, err := trx.Exec(migrationInsertQuery, fileName, time.Now())
		if err != nil {
			trx.Rollback()
			utils.LogObj.Error("file: "+file, "err: "+err.Error(), "Migration: failed to update migration metadata")
			return
		}

		id, err := rslt.LastInsertId()
		utils.LogObj.Info("file: "+file, "id:"+fmt.Sprintf("%d", id), "err: "+err.Error(), "Migration: inserted migration metadata")

		err = trx.Commit()
		if err != nil {
			utils.LogObj.Fatal("file: "+file, "err: "+err.Error(), "Migration: failed to commit file content")
			return
		}

		utils.LogObj.Info("file: "+file, "Migration: done")
	}

	utils.LogObj.Info("Migration: successful")
}

// readMetadata reads the metadata from the migration table
func readMetadata(dbConn *sql.DB) ([]*MigrationRecords, error) {

	rows, err := dbConn.Query(migrationReadQuery)
	if err != nil {
		return nil, err
	}

	records := make([]*MigrationRecords, 0)

	for rows.Next() {

		var (
			version   string
			updatedAt time.Time
		)

		if err := rows.Scan(&version, &updatedAt); err != nil {
			return nil, err
		}

		rcrd := &MigrationRecords{
			Version:   version,
			UpdatedAt: updatedAt,
		}

		records = append(records, rcrd)
	}

	return records, nil
}

// createMigrationTable created a migration table if not found
func createMigrationTable(dbConn *sql.DB) (sql.Result, error) {

	rslt, err := dbConn.Exec(createMigrationTableQuery)
	if err != nil {
		return nil, err
	}

	return rslt, nil
}

func stringToNum(s string) int {
	fileName1 := strings.Split(filepath.Base(s), ".sql")[0]
	num := 0
	for i := range fileName1 {
		if fileName1[i] <= '9' && fileName1[i] >= '0' {
			num = num*10 + int(fileName1[i]-'0')
		} else {
			break
		}
	}
	return num
}

// contains checks whether an array of strings contains a string
func contains(s string, a []*MigrationRecords) bool {
	for _, k := range a {
		if s == k.Version {
			return true
		}
	}
	return false
}
