package main

import (
	"flag"
)

var (
	FlagConfigFilePath *string
)

func initFlags() {
	FlagConfigFilePath = flag.String(
		"config",
		"./config.yml",
		"Path to the configuration file",
	)
}
