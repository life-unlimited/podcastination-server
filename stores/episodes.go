package stores

import (
	"database/sql"
	"fmt"
	"life-unlimited/podcastination/podcasts"
	"time"
)

const episodeSelect = "select id, title, subtitle, date, author, description, mp3_location, season_id, num, image_location, yt_url, mp3_length, is_available from episodes"

type EpisodeStore struct {
	DB *sql.DB
}

// All retrieves all episodes from the store.
func (s *EpisodeStore) All() ([]podcasts.Episode, error) {
	rows, err := s.DB.Query(fmt.Sprintf("%s;", episodeSelect))
	if err != nil {
		return nil, fmt.Errorf("could not query db for episodes: %v", err)
	}
	defer CloseRows(rows)

	episodes, err := parseRowsAsEpisodes(rows)
	if err != nil {
		return nil, fmt.Errorf("could not parse episode rows: %v", err)
	}
	return episodes, nil
}

// ByPodcast retrieves all episodes from the store that belong to the given podcast.
func (s *EpisodeStore) ByPodcast(podcastId int) ([]podcasts.Episode, error) {
	rows, err := s.DB.Query(fmt.Sprintf("%s join seasons on episodes.season_id = seasons.id where seasons.podcast_id = $1", episodeSelect), podcastId)
	if err != nil {
		return nil, fmt.Errorf("could not query db for episodes by podcast: %v", err)
	}
	defer CloseRows(rows)

	episodes, err := parseRowsAsEpisodes(rows)
	if err != nil {
		return nil, fmt.Errorf("could not parse episode rows: %v", err)
	}
	return episodes, nil
}

// BySeason retrieves all episodes from the store that belong to the given season.
func (s *EpisodeStore) BySeason(seasonId int) ([]podcasts.Episode, error) {
	rows, err := s.DB.Query(fmt.Sprintf("%s where seasons.id = $1", episodeSelect), seasonId)
	if err != nil {
		return nil, fmt.Errorf("could not query db for episodes by season: %v", err)
	}
	defer CloseRows(rows)

	episodes, err := parseRowsAsEpisodes(rows)
	if err != nil {
		return nil, fmt.Errorf("could not parse episode rows: %v", err)
	}
	return episodes, nil
}

// ById retrieves an episode from the store by id.
func (s *EpisodeStore) ById(id int) (*podcasts.Episode, error) {
	rows, err := s.DB.Query(fmt.Sprintf("%s where id = $1", episodeSelect), id)
	if err != nil {
		return nil, fmt.Errorf("could not query db for episode by id: %v", err)
	}
	defer CloseRows(rows)

	episodes, err := parseRowsAsEpisodes(rows)
	if err != nil {
		return nil, fmt.Errorf("could not parse episode row: %v", err)
	}
	if len(episodes) != 1 {
		return nil, fmt.Errorf("get episode by id returned %d results, but wanted 1", len(episodes))
	}
	return &episodes[0], nil
}

// parseRowsAsEpisodes parses rows retrieved from db as episodes.
func parseRowsAsEpisodes(rows *sql.Rows) ([]podcasts.Episode, error) {
	var (
		id            int
		title         string
		subtitle      sql.NullString
		date          time.Time
		author        sql.NullString
		description   sql.NullString
		mp3Location   sql.NullString
		mp3Length     int
		seasonId      int
		num           int
		imageLocation sql.NullString
		ytURL         sql.NullString
		isAvailable   bool
	)

	var episodes []podcasts.Episode
	for rows.Next() {
		err := rows.Scan(&id, &title, &subtitle, &date, &author, &description, &mp3Location, &seasonId, &num,
			&imageLocation, &ytURL, &mp3Length, &isAvailable)
		if err != nil {
			return nil, err
		}
		episodes = append(episodes, podcasts.Episode{
			Id:            id,
			Title:         title,
			Subtitle:      subtitle.String,
			Date:          date,
			Author:        author.String,
			Description:   description.String,
			ImageLocation: imageLocation.String,
			MP3Location:   mp3Location.String,
			YouTubeURL:    ytURL.String,
			SeasonId:      seasonId,
			Num:           num,
			MP3Length:     mp3Length,
			IsAvailable:   isAvailable,
		})
	}
	return episodes, nil
}

const episodeInsert = `INSERT INTO episodes (title, subtitle, date, author, description, mp3_location, season_id, num,
                      image_location, yt_url, mp3_length, is_available)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id`

// Create inserts a new episode into db and returns the episode with the assigned id.
func (s *EpisodeStore) Create(e podcasts.Episode) (podcasts.Episode, error) {
	var id int
	err := s.DB.QueryRow(episodeInsert, e.Title, e.Subtitle, e.Date, e.Author, e.Description, e.MP3Location, e.SeasonId,
		e.Num, e.ImageLocation, e.YouTubeURL, e.MP3Length, e.IsAvailable).Scan(&id)
	if err != nil {
		return podcasts.Episode{}, fmt.Errorf("could not insert episode into db: %v", err)
	}
	res := e
	res.Id = id
	return res, nil
}

const episodeUpdate = `UPDATE episodes
SET title=$1, subtitle=$2, date=$3, author=$4, description=$5, mp3_location=$6, season_id=$7, num=$8,
    image_location=$9, yt_url=$10, mp3_length=$11, is_available=$12
WHERE id=$13`

// Update updates an episode in the db based on its id.
func (s *EpisodeStore) Update(e podcasts.Episode) error {
	id := -1
	err := s.DB.QueryRow(episodeUpdate, e.Title, e.Subtitle, e.Date, e.Author, e.Description, e.MP3Location, e.SeasonId,
		e.Num, e.ImageLocation, e.YouTubeURL, e.MP3Length, e.IsAvailable, e.Id).Scan(&id)
	if err != nil {
		return fmt.Errorf("could not update episode in db: %v", err)
	}
	if id == -1 {
		return fmt.Errorf("could not update episode in db: episode %d not found", e.Id)
	}
	return nil
}
