package main

import (
	"os"

	"github.com/maragudk/honeycomb-cli/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
