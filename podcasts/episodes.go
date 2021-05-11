package podcasts

import "time"

type Episode struct {
	Id            int       `json:"id"`
	Title         string    `json:"title"`
	Subtitle      string    `json:"subtitle"`
	Date          time.Time `json:"date"`
	Author        string    `json:"author"`
	Description   string    `json:"description"`
	ImageLocation string    `json:"image_location"`
	PDFLocation   string    `json:"pdf_location"`
	MP3Location   string    `json:"mp3_location"`
	MP3Length     int       `json:"mp3_length"`
	SeasonId      int       `json:"season_id"`
	Num           int       `json:"num"`
	YouTubeURL    string    `json:"yt_url"`
	IsAvailable   bool      `json:"is_available"`
}
