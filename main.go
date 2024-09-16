package main

import (
	"fmt"
	"os"

	"github.com/yawn/spottty/command"
)

func main() {

	if err := command.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

}
