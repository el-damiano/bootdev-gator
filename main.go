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
	queries *database.Queries
	config  *config.Config
}

func main() {

	configFile, err := config.Read()
	if err != nil {
		log.Fatalf("error reading file %v\n", err)
	}
	fmt.Printf("read config %+v\n", configFile)

	mainState := &state{
		config: &configFile,
	}

	commands := commandRegistry{
		reg: map[string]func(*state, command) error{},
	}

	sysArgs := os.Args
	if len(sysArgs) < 2 {
		log.Fatalf("not enough arguments were provided, wanted 2 but got %d\n", len(sysArgs))
		os.Exit(1)
	}

	commandLogin := &command{
		Name: sysArgs[1],
		Args: sysArgs[2:],
	}
	commands.register(commandLogin.Name, handlerLogin)
	err = commands.run(mainState, *commandLogin)
	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}

	configFile, err = config.Read()
	if err != nil {
		log.Fatalf("error reading file %v\n", err)
	}

	fmt.Printf("new config set %+v\n", configFile)

	db, err := sql.Open("postgres", mainState.config.DatabaseUrl)
	if err != nil {
		log.Fatalf("error opening database %v\n", err)
	}

	mainState.queries = database.New(db)
}
