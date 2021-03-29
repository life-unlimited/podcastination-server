package podcasts

import "time"

type Episode struct {
	Title       string
	Subtitle    string
	Date        time.Time
	Author      string
	Description string
	Season      int // TODO: Separate structure
	Episode     int // TODO: To season
	MP3Location string
}
