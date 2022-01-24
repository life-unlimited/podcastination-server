// Package feedgen is used for generating the podcast xml.
package feedgen

import (
	"encoding/xml"
	"fmt"
	"github.com/life-unlimited/podcastination-server/podcast_xml"
	"github.com/life-unlimited/podcastination-server/podcasts"
	"github.com/life-unlimited/podcastination-server/stores"
	"github.com/life-unlimited/podcastination-server/transfer"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
)

// RefreshFeedForPodcasts generates feeds for all podcasts.
func RefreshFeedForPodcasts(store stores.Stores, staticContentURL, podcastDir, feedFileName string) error {
	storePodcasts, err := store.Podcasts.All()
	if err != nil {
		return errors.Wrap(err, "get all podcasts from store")
	}
	owners, err := store.Owners.All()
	if err != nil {
		return errors.Wrap(err, "get all owners from store")
	}
	seasons, err := store.Seasons.All()
	if err != nil {
		return errors.Wrap(err, "get all seasons from store")
	}
	episodes, err := store.Episodes.All()
	if err != nil {
		return errors.Wrap(err, "get all episodes from store")
	}
	// For each podcast, we generate the feed.
	for _, podcast := range storePodcasts {
		creationDetails := podcast_xml.CreationDetails{
			StaticContentURL: staticContentURL,
			Podcast:          podcast,
			Seasons:          make([]podcasts.Season, 0),
			Episodes:         make([]podcasts.Episode, 0),
		}
		// Filter owners.
		found := false
		for _, owner := range owners {
			if podcast.OwnerId == owner.Id {
				creationDetails.Owner = owner
				found = true
				break
			}
		}
		if !found {
			return errors.Wrap(err, fmt.Sprintf("owner %d not found for podcast %d", podcast.OwnerId, podcast.Id))
		}
		// Filter seasons.
		knownSeasonsForPodcast := make(map[int]struct{})
		for _, season := range seasons {
			if season.PodcastId == podcast.Id {
				knownSeasonsForPodcast[season.Id] = struct{}{}
				creationDetails.Seasons = append(creationDetails.Seasons, season)
			}
		}
		// Filter episodes.
		for _, episode := range episodes {
			if _, ok := knownSeasonsForPodcast[episode.SeasonId]; ok {
				creationDetails.Episodes = append(creationDetails.Episodes, episode)
			}
		}
		// Generate.
		podcastXML, err := podcast_xml.GeneratePodcastXML(creationDetails)
		if err != nil {
			return errors.Wrap(err, "generate podcast xml")
		}
		// Marshal.
		podcastXMLRaw, err := xml.MarshalIndent(podcastXML, "", "  ")
		if err != nil {
			return errors.Wrap(err, "marshal podcast xml")
		}
		// Write.
		podcastXMLFilePath := filepath.Join(podcastDir, transfer.GetPodcastFolderName(podcast.Id), feedFileName)
		err = ioutil.WriteFile(podcastXMLFilePath, podcastXMLRaw, 0633)
		if err != nil {
			return errors.Wrap(err, "write podcast xml")
		}
	}
	return nil
}
