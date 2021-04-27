package transfer

import (
	"fmt"
	"life-unlimited/podcastination-server/podcasts"
	"path/filepath"
	"strconv"
	"strings"
)

type EpisodeFileLocations struct {
	BaseDir       string
	MP3FileName   string
	ImageFileName string
}

func GetEpisodeFileLocations(episode podcasts.Episode, podcastId int) EpisodeFileLocations {
	folderName := GetEpisodeFolderName(episode, podcastId)
	cleanTitle := filepath.Clean(episode.Title)
	loc := EpisodeFileLocations{
		BaseDir:       folderName,
		MP3FileName:   fmt.Sprintf("%d_%s.mp3", episode.Id, strings.Replace(cleanTitle, " ", "_", -1)),
		ImageFileName: fmt.Sprintf("thumb.png"),
	}
	return loc
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
