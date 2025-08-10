package main

import (
	"flag"
	"fmt"
	"os"
	"tuido/app"
	"tuido/config"
	"tuido/data"
	"tuido/webapi"
)

func main() {
	webapp := flag.Bool("webapi", false, "Run the web API server instead of the TUI application")
	storageType := flag.String("config.storageType", "", "Storage type (local or remote)")
	remoteUrl := flag.String("config.remoteUrl", "", "Remote URL for the API, required if storage type is remote")

	flag.Parse()

	if *webapp {
		webapi.Run()
		return
	}

	conf, err := config.Load()

	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	setConfig(storageType, remoteUrl, &conf)

	app.Run(getProvider(conf))

}

func setConfig(storageType *string, remoteUrl *string, conf *config.Config) {

	save := false

	if conf.StorageType == "" {
		conf.StorageType = "local"
		save = true
	}

	if *storageType != "" {
		conf.StorageType = *storageType

		save = true
	}

	if *remoteUrl != "" {
		conf.RemoteUrl = *remoteUrl
		save = true
	}

	if conf.StorageType == "remote" && conf.RemoteUrl == "" {
		fmt.Println("config.remoteUrl is required when using config.storageType == remote.")
		os.Exit(1)
	}

	if save {
		err := config.Save(*conf)

		if err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}
	}
}

func getProvider(conf config.Config) data.Provider {
	switch conf.StorageType {
	case "local":
		return data.SqliteProvider{}
	case "remote":
		return data.HttpProvider{Url: conf.RemoteUrl}
	default:
		return data.SqliteProvider{}
	}
}
