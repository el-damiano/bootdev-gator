package main

import (
	"fmt"
	"log"

	"github.com/el-damiano/bootdev-gator/internal/config"
)

func main() {
	configFile, err := config.Read()
	if err != nil {
		log.Fatalf("error reading file %v\n", err)
	}
	fmt.Printf("read config %+v\n", configFile)

	err = configFile.SetUser("damian")
	if err != nil {
		log.Fatalf("error setting user %v\n", err)
	}

	configFile, err = config.Read()
	if err != nil {
		log.Fatalf("error reading file %v\n", err)
	}

	fmt.Printf("new config set %+v\n", configFile)
}
