package main

import (
	"fmt"
	"log"
	"os"

	"github.com/el-damiano/bootdev-gator/internal/config"
)

type state struct {
	Config *config.Config
}

func main() {

	configFile, err := config.Read()
	if err != nil {
		log.Fatalf("error reading file %v\n", err)
	}
	fmt.Printf("read config %+v\n", configFile)

	mainState := &state{
		Config: &configFile,
	}
	_ = mainState

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
}
