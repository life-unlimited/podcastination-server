package podcast_xml

type PodcastXML struct {
	RSS     rss     `xml:"rss"`
	Channel channel `xml:"channel"`
}

type rss struct {
	XmlnsAtom    string `xml:"xmlns:atom,attr"`
	XmlnsContent string `xml:"xmlns:content,attr"`
	XmlnsITunes  string `xml:"xmlns:itunes,attr"`
	XmlnsGPlay   string `xml:"xmlns:googleplay,attr"`
	XmlnsMedia   string `xml:"xmlns:media,attr"`
	Version      string `xml:"version,attr"`
}

type channel struct {
	Title            string           `xml:"title"`
	Link             string           `xml:"link"`
	Language         string           `xml:"language"`
	AtomLink         atomLink         `xml:"atom:link"`
	Copyright        string           `xml:"copyright"`
	ITunesSubtitle   string           `xml:"itunes:subtitle"`
	ITunesAuthor     string           `xml:"itunes:author"`
	ITunesSummary    string           `xml:"itunes:summary"`
	ITunesKeywords   string           `xml:"itunes:keywords"`
	Description      string           `xml:"description"`
	ITunesOwner      iTunesOwner      `xml:"itunes:owner"`
	Image            image            `xml:"image"`
	ITunesImage      iTunesImage      `xml:"itunes:image"`
	ITunesCategories []iTunesCategory `xml:"itunes:category"`
	Items            []item           `xml:"item"`
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
	Text string `xml:"text,attr"`
}

type item struct {
	Title             string      `xml:"title"`
	ITunesTitle       string      `xml:"itunes:title"`
	ITunesAuthor      string      `xml:"itunes:author"`
	ITunesSubTitle    string      `xml:"itunes:subtitle"`
	ITunesSummary     string      `xml:"itunes:summary"`
	ITunesImage       iTunesImage `xml:"itunes:image"`
	Enclosure         enclosure   `xml:"enclosure"`
	ITunesDuration    int         `xml:"itunes:duration"`
	ITunesSeason      int         `xml:"itunes:season"`
	ITunesEpisode     int         `xml:"itunes:episode"`
	ITunesEpisodeType string      `xml:"itunes:episodeType"`
	Guid              guid        `xml:"guid"`
	PubDate           string      `xml:"pubDate"`
	ITunesExplicit    string      `xml:"itunes:explicit"`
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
