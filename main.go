package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

type config struct {
	Directory  string `json:"directory"`
	DbAddress  string `json:"dbAddress"`
	DbPort     string `json:"dbPort"`
	DbUsername string `json:"dbUsername"`
	DbPassword string `json:"dbPassword"`
	DbDatabase string `json:"dbDatabase"`
	DbSchema   string `json:"dbSchema"`
}

func isValidConnectionString(connectionDetails string) bool {
	if strings.Count(connectionDetails, ":") < 2 {
		return false
	}
	if strings.Count(connectionDetails, "@") < 1 {
		return false
	}
	if strings.Count(connectionDetails, "/") < 1 {
		return false
	}
	return true
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: sqldatadump [--directory=<directory>] <username>:<password>@<address>:<port>/<database>\n")
	}

	cfg := &config{}

	flag.StringVar(&cfg.Directory, "directory", "", "root directory to export to")

	flag.Parse()

	// expecting string of form <username>:<password>@<address>:<port>/<database>
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(-1)
	}

	if cfg.Directory == "" {
		flag.Usage()
		os.Exit(-1)
	}

	newpath := filepath.Join(cfg.Directory)
	err := os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	connectionDetails := args[0]
	if !isValidConnectionString(connectionDetails) {
		flag.Usage()
		os.Exit(-1)
	}

	cfg.DbUsername = connectionDetails[:strings.Index(connectionDetails, ":")]
	cfg.DbPassword = connectionDetails[strings.Index(connectionDetails, ":")+1 : strings.Index(connectionDetails, "@")]
	cfg.DbAddress = connectionDetails[strings.Index(connectionDetails, "@")+1 : strings.LastIndex(connectionDetails, ":")]
	cfg.DbPort = connectionDetails[strings.LastIndex(connectionDetails, ":")+1 : strings.LastIndex(connectionDetails, "/")]
	cfg.DbDatabase = connectionDetails[strings.Index(connectionDetails, "/")+1:]

	mssqldb, err := sqlx.Connect("sqlserver", fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;port=%s", cfg.DbAddress, cfg.DbUsername, cfg.DbPassword, cfg.DbDatabase, cfg.DbPort))
	if err != nil {
		log.Fatal(err)
	}

	filesToExecute := []string{}

	err = filepath.Walk(cfg.Directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if filepath.Ext(info.Name()) == ".sql" {
					filesToExecute = append(filesToExecute, path)
				}
			}
			return nil
		},
	)
	if err != nil {
		log.Println(err)
	}

	for _, f := range filesToExecute {
		log.Default().Printf("processing %s", f)
		s := time.Now()
		contents, err := os.ReadFile(f)
		if err != nil {
			log.Println(err)
		}

		log.Default().Printf("read file %s in %s", f, time.Since(s))
		s1 := time.Now()

		res, err := mssqldb.Exec(string(contents))
		if err != nil {
			log.Println(err)
			log.Default().Printf("aborted file %s in %s", f, time.Since(s1))
			continue
		}
		res.RowsAffected()

		log.Default().Printf("executed file %s in %s", f, time.Since(s1))
		s2 := time.Now()

		log.Default().Printf("completed file %s in %s", f, time.Since(s2))
	}
}
