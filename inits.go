package main

import (
	"flag"
	"log"
	"os"

	"github.com/coffemanfp/beppin-server/config"
	"github.com/coffemanfp/beppin-server/database"
)

// Flags
var (
	readConfigFile bool
	readEnvVars    bool
)

var filesToUpload = make(chan *os.File)

func initSettings() {
	config.SetDefaultSettings()

	if readEnvVars {
		config.SetSettingsByEnv()
	}

	if readConfigFile {
		err := config.SetSettingsByFile()
		if err != nil {
			log.Fatalln("failed to configure file settings: ", err)
		}
	}
}

func initDatabase() {
	_, err := database.NewDefault()
	if err != nil {
		log.Fatalln("failed to start the database: ", err)
	}
}

func initFlags() {
	flag.BoolVar(&readConfigFile, "read-config-file", true, "Checks if read the config file or not.")
	flag.BoolVar(&readEnvVars, "read-env-vars", true, "Checks if read the environment vars or not.")

	flag.Parse()
}
