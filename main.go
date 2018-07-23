package main

import "github.com/opb/secretly/cmd"

var Version = "unknown"

func main() {
	cmd.Version = Version
	cmd.Execute()
}
