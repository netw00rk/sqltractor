// Package main is the CLI.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"

	"github.com/netw00rk/sqltractor/tractor"
	"github.com/netw00rk/sqltractor/tractor/direction"
	"github.com/netw00rk/sqltractor/tractor/file"

	_ "github.com/netw00rk/sqltractor/driver/cassandra"
	_ "github.com/netw00rk/sqltractor/driver/mysql"
	_ "github.com/netw00rk/sqltractor/driver/postgres"
	_ "github.com/netw00rk/sqltractor/driver/sqlite3"
)

var url = flag.String("url", os.Getenv("MIGRATE_URL"), "")
var path = flag.String("path", "", "")

func main() {
	flag.Parse()
	command := flag.Arg(0)
	if command == "" || command == "help" {
		printHelpCmd()
		os.Exit(0)
	}

	if *path == "" {
		var err error
		if *path, err = os.Getwd(); err != nil {
			fmt.Println("Please specify path")
			os.Exit(1)
		}
	}

	tractor, err := tractor.NewTractor(*url, *path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch command {
	case "migrate":
		relativeN, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			fmt.Println("Unable to parse param <n>.")
			os.Exit(1)
		}

		timerStart := time.Now()
		for r := range tractor.MigrateAsync(relativeN) {
			if r.Error != nil {
				printFile(r.File, err)
				os.Exit(1)
			}
			printFile(r.File, nil)
		}
		printTimer(timerStart)

	case "goto":
		toVersion, err := strconv.Atoi(flag.Arg(1))
		if err != nil || toVersion < 0 {
			fmt.Println("Unable to parse param <v>.")
			os.Exit(1)
		}

		currentVersion, err := tractor.Version()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		relativeN := toVersion - int(currentVersion)

		timerStart := time.Now()
		for r := range tractor.MigrateAsync(relativeN) {
			if r.Error != nil {
				printFile(r.File, err)
				os.Exit(1)
			}
			printFile(r.File, nil)
		}
		printTimer(timerStart)

	case "up":
		timerStart := time.Now()
		for r := range tractor.UpAsync() {
			if r.Error != nil {
				printFile(r.File, err)
				os.Exit(1)
			}
			printFile(r.File, nil)
		}
		printTimer(timerStart)

	case "down":
		timerStart := time.Now()
		for r := range tractor.DownAsync() {
			if r.Error != nil {
				printFile(r.File, err)
				os.Exit(1)
			}
			printFile(r.File, nil)
		}
		printTimer(timerStart)

	case "version":
		version, err := tractor.Version()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(version)
	}
}

func printFile(f *file.File, err error) {
	if err != nil {
		c := color.New(color.FgRed)
		c.Println(err.Error(), "\n")
		return
	}

	c := color.New(color.FgBlue)
	if f.Direction == direction.Up {
		c.Print(">")
	} else if f.Direction == direction.Down {
		c.Print("<")
	}
	fmt.Printf(" %s\n", f.FileName)
}

func printTimer(start time.Time) {
	diff := time.Now().Sub(start).Seconds()
	if diff > 60 {
		fmt.Printf("\n%.4f minutes\n", diff/60)
	} else {
		fmt.Printf("\n%.4f seconds\n", diff)
	}
}

func printHelpCmd() {
	os.Stderr.WriteString(
		`usage: sqltractor [-path=<path>] -url=<url> <command> [<args>]

Commands:
   create <name>  Create a new migration
   up             Apply all -up- migrations
   down           Apply all -down- migrations
   reset          Down followed by Up
   redo           Roll back most recent migration, then apply it again
   version        Show current migration version
   migrate <n>    Apply migrations -n|+n
   goto <v>       Migrate to version v
   help           Show this help

'-path' defaults to current working directory.
`)
}
