package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/el-damiano/bootdev-gator/internal/config"
	"github.com/el-damiano/bootdev-gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {

	configFile, err := config.Read()
	if err != nil {
		log.Fatalf("error reading file %v\n", err)
	}
	fmt.Printf("read config %+v\n", configFile)

	db, err := sql.Open("postgres", configFile.DatabaseUrl)
	if err != nil {
		log.Fatalf("error opening database %v\n", err)
	}

	dbQueries := database.New(db)
	mainState := &state{
		config: &configFile,
		db:     dbQueries,
	}

	commands := commandRegistry{
		reg: map[string]func(*state, command) error{},
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)

	sysArgs := os.Args
	if len(sysArgs) < 2 {
		log.Fatalf("not enough arguments were provided, wanted 2 but got %d\n", len(sysArgs))
		os.Exit(1)
	}

	commandName := sysArgs[1]
	commandArgs := sysArgs[2:]
	err = commands.run(mainState, command{Name: commandName, Args: commandArgs})
	if err != nil {
		log.Fatal(err)
	}

	configFile, err = config.Read()
	if err != nil {
		log.Fatalf("error reading file %v\n", err)
	}
	fmt.Printf("new config set %+v\n", configFile)
}
