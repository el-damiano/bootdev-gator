package main

import (
	"fmt"

	"github.com/el-damiano/bootdev-gator/internal/config"
)

func main() {
	configFile, err := config.Read()
	if err != nil {
		fmt.Printf("error reading file %s\n", err)
	}

	err = configFile.SetUser("damian")
	if err != nil {
		fmt.Printf("error setting user %s\n", err)
	}

	configFile, err = config.Read()
	if err != nil {
		fmt.Printf("error reading file %s\n", err)
	}

	fmt.Printf("new config set %+v\n", configFile)
}
