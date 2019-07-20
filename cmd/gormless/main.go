package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"regexp"
	"sort"
	"strings"

	"github.com/jinzhu/gorm"

	gormigrate "gopkg.in/gormigrate.v1"
)

const (
	exitOk = iota
	exitErr
)

var version string

var r = regexp.MustCompile(string(os.PathSeparator))
var validActions = [6]string{
	"version", "commit", "migrate", "revert", "rollback", "run",
}
var defaultDir = "migration"

var action = flag.String(
	"action",
	"",
	"version, commit (alias: migrate, run) or revert (alias: rollback)",
)
var migrationDirectory = flag.String(
	"migration-directory",
	defaultDir,
	"path to directory contains migration files",
)

func detectMigrationDirectory(directoryPath string) string {
	if directoryPath != defaultDir { // as an argument
		return directoryPath
	}

	val, ok := os.LookupEnv("MIGRATION_DIRECTORY")
	if !ok || val == "" {
		return defaultDir
	}
	return val
}

func pathToID(path string) string {
	parts := r.Split(path, -1)
	l := len(parts)
	if l < 1 {
		return ""
	}
	return strings.Replace(parts[l-1], ".so", "", 1)
}

func validateAction(s string) bool {
	for _, a := range validActions {
		if a == s {
			return true
		}
	}
	return false
}

func validateMigrationDirectory(d string) bool {
	stat, err := os.Stat(d)
	return err == nil && stat.IsDir()
}

func loadMigrations(pattern string) ([]*gormigrate.Migration, error) {
	var err error

	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})

	var migrations []*gormigrate.Migration
	for _, f := range files {
		plug, err := plugin.Open(f)
		if err != nil {
			return nil, err
		}
		up, err := plug.Lookup("Up")
		if err != nil {
			return nil, err
		}
		down, err := plug.Lookup("Down")
		if err != nil {
			return nil, err
		}

		id := pathToID(f)
		if len(id) < 1 {
			continue
		}
		migrations = append(migrations, &gormigrate.Migration{
			ID:       id,
			Migrate:  func(tx *gorm.DB) error { return up.(func(*gorm.DB) error)(tx) },
			Rollback: func(tx *gorm.DB) error { return down.(func(*gorm.DB) error)(tx) },
		})
	}

	return migrations, nil
}

func run(db *gorm.DB, actionName, directoryPath string) error {
	pattern := filepath.Join(directoryPath, "**", "*.so")
	migrations, err := loadMigrations(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err)
		return err
	}

	if len(migrations) < 1 {
		log.Printf("No migrations")
		return nil
	}

	migrator := gormigrate.New(db, gormigrate.DefaultOptions, migrations)

	db.LogMode(true)
	log.Printf("Migration (%s) has been started", actionName)

	switch actionName {
	case "commit", "migrate", "run":
		if err = migrator.Migrate(); err != nil {
			log.Fatalf("Cound not commit: %v", err)
			return err
		}
	case "revert", "rollback":
		if err = migrator.RollbackLast(); err != nil {
			log.Fatalf("Cound not revert: %v", err)
			return err
		}
	default:
		log.Printf("Unknown action :'(")
		return err
	}

	fmt.Println("")
	log.Printf("Migration (%s) has been finished", actionName)
	return nil
}

func main() {
	os.Exit(realMain(os.Args))
}

func realMain(args []string) int {
	flag.Parse()

	// -action
	var a string
	if len(*action) != 0 {
		a = *action
	} else if len(args) >= 2 {
		for _, arg := range args[1:] {
			if !strings.HasPrefix(arg, "-") {
				a = arg
				break
			}
		}
	}
	if len(a) == 0 {
		log.Printf("An action is required :'(")
		return exitErr
	}
	if !validateAction(a) {
		log.Printf("The action is invalid :'(")
		return exitErr
	}

	if a == "version" {
		fmt.Printf("%s version %s\n", filepath.Base(args[0]), version)
		return exitOk
	}

	// -migration-directory
	d := detectMigrationDirectory(*migrationDirectory)
	if !validateMigrationDirectory(d) {
		fmt.Fprintf(os.Stderr, "No such directory %s\n", d)
		return exitErr
	}

	db, err := connect(os.Getenv("DATABASE_URL"))
	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect database: %s\n", err)
		return exitErr
	}

	if err := run(db, a, d); err == nil {
		return exitOk
	}
	return exitErr
}
