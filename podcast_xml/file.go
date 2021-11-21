package podcast_xml

import "encoding/xml"

type PodcastXML struct {
	XMLName      xml.Name `xml:"rss"`
	XmlnsAtom    string   `xml:"xmlns:atom,attr"`
	XmlnsContent string   `xml:"xmlns:content,attr"`
	XmlnsITunes  string   `xml:"xmlns:itunes,attr"`
	XmlnsGPlay   string   `xml:"xmlns:googleplay,attr"`
	XmlnsMedia   string   `xml:"xmlns:media,attr"`
	Version      string   `xml:"version,attr"`
	Channel      channel  `xml:"channel"`
}

type channel struct {
	Title          string         `xml:"title"`
	Link           string         `xml:"link"`
	Language       string         `xml:"language"`
	AtomLink       atomLink       `xml:"atom:link"`
	Copyright      string         `xml:"copyright"`
	ITunesSubtitle string         `xml:"itunes:subtitle"`
	ITunesAuthor   string         `xml:"itunes:author"`
	ITunesSummary  string         `xml:"itunes:summary"`
	ITunesKeywords string         `xml:"itunes:keywords"`
	Description    string         `xml:"description"`
	ITunesOwner    iTunesOwner    `xml:"itunes:owner"`
	Image          image          `xml:"image"`
	ITunesImage    iTunesImage    `xml:"itunes:image"`
	ITunesCategory iTunesCategory `xml:"itunes:category"`
	Items          []item         `xml:"item"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type iTunesOwner struct {
	ITunesName  string `xml:"itunes:name"`
	ITunesEmail string `xml:"itunes:email"`
}

type image struct {
	Url string `xml:"url"`
}

type iTunesImage struct {
	Href string `xml:"href,attr"`
}

type iTunesCategory struct {
	Text          string           `xml:"text,attr"`
	SubCategories []iTunesCategory `xml:"itunes:category"`
}

type item struct {
	// Title holds the title as well as the ITunesSubTitle.
	Title             string      `xml:"title"`
	ITunesTitle       string      `xml:"itunes:title,omitempty"`
	ITunesAuthor      string      `xml:"itunes:author,omitempty"`
	ITunesSubTitle    string      `xml:"itunes:subtitle,omitempty"`
	ITunesSummary     string      `xml:"itunes:summary,omitempty"`
	ITunesImage       iTunesImage `xml:"itunes:image,omitempty"`
	Enclosure         enclosure   `xml:"enclosure,omitempty"`
	ITunesDuration    int         `xml:"itunes:duration"`
	ITunesSeason      int         `xml:"itunes:season"`
	ITunesEpisode     int         `xml:"itunes:episode"`
	ITunesEpisodeType string      `xml:"itunes:episodeType,omitempty"`
	Guid              guid        `xml:"guid"`
	PubDate           string      `xml:"pubDate"`
	ITunesExplicit    string      `xml:"itunes:explicit,omitempty"`
}

type enclosure struct {
	URL    string `xml:"url,attr"`
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

type guid struct {
	IsPermaLink bool   `xml:"isPermaLink"`
	Location    string `xml:",chardata"`
}
