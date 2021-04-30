package transfer

import (
	"fmt"
	"life-unlimited/podcastination-server/podcasts"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type EpisodeFileLocations struct {
	BaseDir       string
	MP3FileName   string
	ImageFileName string
}

// GetEpisodeFileLocations returns the file locations for the given episode and podcast.
func GetEpisodeFileLocations(episode podcasts.Episode, podcastId int) EpisodeFileLocations {
	folderName := GetEpisodeFolderName(episode, podcastId)
	cleanTitle := filepath.Clean(removeSpecialCharacters(replaceSpacesWithUnderscore(episode.Title)))
	loc := EpisodeFileLocations{
		BaseDir:       folderName,
		MP3FileName:   fmt.Sprintf("%d_%s.mp3", episode.Id, cleanTitle),
		ImageFileName: fmt.Sprintf("thumb.png"),
	}
	return loc
}

// replaceSpacesWithUnderscore returns the given string with spaces being replaced with underscores.
func replaceSpacesWithUnderscore(s string) string {
	return strings.Replace(s, " ", "_", -1)
}

// removeSpecialCharacters returns the given string without special characters.
func removeSpecialCharacters(s string) string {
	// Regex for only letters and numbers.
	reg, err := regexp.Compile("[^a-zA-Z0-9_]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(s, "")
}

func (loc EpisodeFileLocations) MP3FullPath() string {
	return filepath.Join(loc.BaseDir, loc.MP3FileName)
}

func (loc EpisodeFileLocations) ImageFullPath() string {
	if loc.ImageFileName == "" {
		return ""
	}
	return filepath.Join(loc.BaseDir, loc.ImageFileName)
}

// GetEpisodeFolderName returns the folder name created from the given episode.
func GetEpisodeFolderName(episode podcasts.Episode, podcastId int) string {
	timestamp := episode.Date.Format("20060102_150405")
	return filepath.Join(GetPodcastFolderName(podcastId), fmt.Sprintf("%s_%d", timestamp, episode.Id))
}

func GetPodcastFolderName(podcastId int) string {
	return strconv.Itoa(podcastId)
}
