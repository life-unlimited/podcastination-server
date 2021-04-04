package main

import (
	"flag"
	"life-unlimited/podcastination/app"
	"life-unlimited/podcastination/config"
	"log"
)

func main() {
	// Flags.
	configPath := flag.String("config", "config.json", "Path to the config file")
	flag.Parse()
	// Read config.
	podcastinationConfig, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatalf("could not read config: %v", err)
	}
	printDirs(podcastinationConfig)
	// Create the app.
	podcastination := app.NewApp(podcastinationConfig)
	// Boot.
	if err := podcastination.Boot(); err != nil {
		panic(err)
	}
}

func printDirs(podcastinationConfig config.PodcastinationConfig) {
	log.Printf("using pull dir %s", podcastinationConfig.PullDir)
	log.Printf("using podcast dir %s", podcastinationConfig.PodcastDir)
}
