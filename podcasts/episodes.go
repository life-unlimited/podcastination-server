package podcasts

import "time"

type Episode struct {
	Id            int
	Title         string
	Subtitle      string
	Date          time.Time
	Author        string
	Description   string
	ImageLocation string
	MP3Location   string
	SeasonId      int
	Num           int
}
