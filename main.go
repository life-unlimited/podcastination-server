package main

import (
	"flag"
	"github.com/life-unlimited/podcastination-server/app"
	"github.com/life-unlimited/podcastination-server/config"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	log.Println("starting...")
	if err := podcastination.Boot(); err != nil {
		panic(err)
	}
	log.Println("up and running!")
	// Await term signal.
	awaitTerminateSignal()
	log.Println("shutting down...")
	// Shutdown
	if err := podcastination.Shutdown(); err != nil {
		log.Printf("could not shutdown podcastination: %v", err)
	}
	log.Println("good bye!")
}

func awaitTerminateSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
}

func printDirs(podcastinationConfig config.PodcastinationConfig) {
	log.Printf("using pull dir %s", podcastinationConfig.PullDir)
	log.Printf("using podcast dir %s", podcastinationConfig.PodcastDir)
}
