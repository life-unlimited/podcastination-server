package podcast_xml

import (
	"fmt"
	"github.com/life-unlimited/podcastination-server/podcasts"
	"sort"
	"strconv"
	"strings"
	"time"
)

// CreationDetails are needed in order to create a new PodcastXML.
type CreationDetails struct {
	StaticContentURL string
	Owner            podcasts.Owner
	Podcast          podcasts.Podcast
	Seasons          []podcasts.Season
	Episodes         []podcasts.Episode
}

// nestedCreationDetails represent a nested version of CreationDetails and are easier to use when creating a PodcastXML.
type nestedCreationDetails struct {
	Owner   podcasts.Owner
	Podcast podcasts.Podcast
	Seasons []nestedSeasonDetails
}

type nestedSeasonDetails struct {
	Details  podcasts.Season
	Episodes []podcasts.Episode
}

// isValid checks if the given CreationDetails are valid and all references exist.
func (details *CreationDetails) isValid() (bool, error) {
	_, err := details.nested()
	return err == nil, err
}

// nested creates a nestedCreationDetails from CreationDetails while also performing validation. Sorting is also
// performed for seasons and episodes.
func (details *CreationDetails) nested() (nestedCreationDetails, error) {
	nested := nestedCreationDetails{
		Owner:   details.Owner,
		Podcast: details.Podcast,
	}
	// Check owner.
	if details.Podcast.OwnerId != details.Owner.Id {
		return nestedCreationDetails{}, fmt.Errorf("podcast owner ref (%d) doet not match the given owner (%d)",
			details.Podcast.OwnerId, details.Owner.Id)
	}
	// Check seasons.
	for _, season := range details.Seasons {
		// Assure correct podcast reference.
		if season.PodcastId != details.Podcast.Id {
			return nestedCreationDetails{}, fmt.Errorf("season podcast ref (%d) does not match the given podcast (%d)",
				season.PodcastId, details.Podcast.Id)
		}
		// Add to nested.
		nested.Seasons = append(nested.Seasons, nestedSeasonDetails{Details: season})
	}
	duplicateSeasonNums := -1
	// Sort seasons and check for duplicate season nums.
	sort.SliceStable(nested.Seasons, func(i, j int) bool {
		vi := nested.Seasons[i].Details.Num
		vj := nested.Seasons[j].Details.Num
		if vi == vj {
			duplicateSeasonNums = vi
		}
		return vi < vj
	})
	if duplicateSeasonNums != -1 {
		return nestedCreationDetails{}, fmt.Errorf("duplicate season number: %d", duplicateSeasonNums)
	}
	// Check episodes.
	seasonCache := struct {
		Id    int
		Index int
	}{
		Id:    -1,
		Index: -1,
	}
	for _, episode := range details.Episodes {
		// Assure existing season id.
		if episode.SeasonId <= 0 {
			return nestedCreationDetails{}, fmt.Errorf("invalid season id (%d) for episode %d",
				episode.SeasonId, episode.Id)
		}
		// Check cache.
		if episode.SeasonId == seasonCache.Id {
			// Add directly.
			nested.Seasons[seasonCache.Id].Episodes = append(nested.Seasons[seasonCache.Id].Episodes, episode)
			continue
		}
		// Not in cache --> search for season.
		season := -1
		for i, s := range details.Seasons {
			if episode.SeasonId == s.Id {
				season = i
				break
			}
		}
		if season == -1 {
			// Not found.
			return nestedCreationDetails{}, fmt.Errorf("episode %d references none of given seasons", episode.Id)
		}
		// Find season in nested.
		nestedSeasonIndex := -1
		for i, s := range nested.Seasons {
			if s.Details.Id == details.Seasons[season].Id {
				nestedSeasonIndex = i
				break
			}
		}
		// Add to episodes list and cache.
		nested.Seasons[nestedSeasonIndex].Episodes = append(nested.Seasons[nestedSeasonIndex].Episodes, episode)
		seasonCache.Id = season
		seasonCache.Index = nestedSeasonIndex
	}
	// Sort episodes and check for duplicate episode nums within the same season.
	for _, season := range nested.Seasons {
		duplicateEpisodeNums := -1
		sort.SliceStable(season.Episodes, func(i, j int) bool {
			vi := season.Episodes[i].Num
			vj := season.Episodes[j].Num
			if vi == vj {
				duplicateEpisodeNums = vi
			}
			return vi < vj
		})
		if duplicateEpisodeNums != -1 {
			return nestedCreationDetails{}, fmt.Errorf("duplicate episode number within season %d: %d",
				season.Details.Id, duplicateSeasonNums)
		}
	}
	// Everything ok.
	return nested, nil
}

// createEmptyPodcastXML creates a new PodcastXML filled with default values.
func createEmptyPodcastXML() *PodcastXML {
	return &PodcastXML{
		XmlnsAtom:    "http://www.w3.org/2005/Atom",
		XmlnsContent: "http://purl.org/rss/1.0/modules/content/",
		XmlnsITunes:  "http://www.itunes.com/dtds/podcast-1.0.dtd",
		XmlnsGPlay:   "http://www.google.com/schemas/play-podcasts/1.0",
		XmlnsMedia:   "http://www.rssboard.org/media-rss",
		Version:      "2.0",
		Channel:      channel{},
	}
}

// setOwner sets the given owner details for a PodcastXML.
func (xml *PodcastXML) setOwner(owner podcasts.Owner) {
	c := xml.Channel
	c.Copyright = owner.Copyright
	c.ITunesAuthor = owner.Name
	c.ITunesOwner = iTunesOwner{
		ITunesName:  owner.Name,
		ITunesEmail: owner.Email,
	}
	xml.Channel = c
}

// setPodcastDetails sets podcast details like title and subtitle for a PodcastXML.
func (xml *PodcastXML) setPodcastDetails(podcast podcasts.Podcast, staticContentURL string) {
	c := xml.Channel
	c.Title = podcast.Title
	c.Link = podcast.Link
	c.Language = string(podcast.Language)
	c.AtomLink = atomLink{
		Href: podcast.FeedLink,
		Rel:  "self",
		Type: "application/rss+xml",
	}
	c.ITunesSubtitle = podcast.Subtitle
	c.ITunesSummary = podcast.Description
	c.ITunesKeywords = strings.Join(podcast.Keywords, ",")
	c.Description = podcast.Description
	if podcast.ImageLocation != "" {
		c.Image = image{
			Url: fmt.Sprintf("%s/%s", staticContentURL, podcast.ImageLocation),
		}
		c.ITunesImage = iTunesImage{
			Href: fmt.Sprintf("%s/%s", staticContentURL, podcast.ImageLocation),
		}
	}
	c.ITunesCategory = iTunesCategory{
		Text:          "Religion &amp; Spirituality",
		SubCategories: []iTunesCategory{{Text: "Christianity"}},
	}
	xml.Channel = c
}

// setItems sets the episodes and seasons for a PodcastXML.
func (xml *PodcastXML) setItems(seasons []nestedSeasonDetails, staticContentURL string) {
	for _, season := range seasons {
		for _, episode := range season.Episodes {
			xml.appendEpisode(episode, season.Details, staticContentURL)
		}
	}
}

// appendEpisode adds an episode to a PodcastXML.
//
// Warning: Always add episodes in the correct order!
func (xml *PodcastXML) appendEpisode(episode podcasts.Episode, season podcasts.Season, staticContentURL string) {
	var iTunesImageVal iTunesImage
	if episode.ImageLocation != "" {
		iTunesImageVal = iTunesImage{
			Href: fmt.Sprintf("%s/%s", staticContentURL, episode.ImageLocation),
		}
	}
	episodeTitle := episode.Title
	// Append subtitle to title if existing.
	if episode.Subtitle != "" {
		episodeTitle = fmt.Sprintf("%s - %s", episode.Title, episode.Subtitle)
	}
	e := item{
		Title:          episodeTitle,
		ITunesTitle:    episode.Title,
		ITunesAuthor:   episode.Author,
		ITunesSubTitle: episode.Subtitle,
		ITunesSummary:  episode.Description,
		ITunesImage:    iTunesImageVal,
		Enclosure: enclosure{
			URL:    fmt.Sprintf("%s/%s", staticContentURL, episode.MP3Location),
			Length: strconv.Itoa(episode.MP3Length),
			Type:   "audio/mpeg",
		},
		ITunesDuration:    episode.MP3Length,
		ITunesSeason:      season.Num,
		ITunesEpisode:     episode.Num,
		ITunesEpisodeType: "full",
		Guid: guid{
			IsPermaLink: false,
			Location:    fmt.Sprintf("%s/%s", staticContentURL, episode.MP3Location),
		},
		PubDate:        episode.Date.Format(time.RFC1123Z),
		ITunesExplicit: "NO", // I guess that this will always be no.
	}
	xml.Channel.Items = append(xml.Channel.Items, e)
}

// GeneratePodcastXML generates a PodcastXML from the given CreationDetails.
func GeneratePodcastXML(details CreationDetails) (PodcastXML, error) {
	nested, err := details.nested()
	if err != nil {
		return PodcastXML{}, fmt.Errorf("invalid creation details: %v", err)
	}
	xml := createEmptyPodcastXML()
	xml.setOwner(nested.Owner)
	xml.setPodcastDetails(nested.Podcast, details.StaticContentURL)
	xml.setItems(nested.Seasons, details.StaticContentURL)
	return *xml, nil
}
