package stores

import (
	"database/sql"
	"fmt"
	"life-unlimited/podcastination/podcasts"
	"time"
)

const episodeSelect = "select id, title, subtitle, date, author, description, mp3_location, season_id, num, image_location from episodes"

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
	rows, err := s.DB.Query(fmt.Sprintf("%s join seasons on episodes.season_id = seasons.id where seasons.podcast_id = ?", episodeSelect), podcastId)
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
	rows, err := s.DB.Query(fmt.Sprintf("%s where seasons.id = ?", episodeSelect), seasonId)
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
	rows, err := s.DB.Query(fmt.Sprintf("%s where id = ?", episodeSelect), id)
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
		subtitle      string
		date          time.Time
		author        string
		description   string
		mp3Location   string
		seasonId      int
		num           int
		imageLocation string
	)

	var episodes []podcasts.Episode
	for rows.Next() {
		err := rows.Scan(&id, &title, &subtitle, &date, &author, &description, &mp3Location, &seasonId, &num, &imageLocation)
		if err != nil {
			return nil, err
		}
		episodes = append(episodes, podcasts.Episode{
			Id:            id,
			Title:         title,
			Subtitle:      subtitle,
			Date:          date,
			Author:        author,
			Description:   description,
			ImageLocation: imageLocation,
			MP3Location:   mp3Location,
			SeasonId:      seasonId,
			Num:           num,
		})
	}
	return episodes, nil
}
