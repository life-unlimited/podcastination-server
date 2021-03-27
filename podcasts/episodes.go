package podcasts

import "time"

type Episode struct {
	Title       string
	Subtitle    string
	Date        time.Time
	Author      string
	Description string
	Season      int
	Episode     int
	MP3Location string
}
