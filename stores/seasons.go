package stores

import (
	"database/sql"
	"fmt"
	"life-unlimited/podcastination/podcasts"
)

const seasonSelect = "select id, title, subtitle, description, image_location, podcast_id, num, key from seasons"

type SeasonStore struct {
	DB *sql.DB
}

// All retrieves all seasons from the store.
func (s *SeasonStore) All() ([]podcasts.Season, error) {
	rows, err := s.DB.Query(fmt.Sprintf("%s;", seasonSelect))
	if err != nil {
		return nil, fmt.Errorf("could not query db for seasons: %v", err)
	}
	defer CloseRows(rows)

	seasons, err := parseRowsAsSeasons(rows)
	if err != nil {
		return nil, fmt.Errorf("could not parse season rows: %v", err)
	}
	return seasons, nil
}

// ById retrieves a season from the store with the given id.
func (s *SeasonStore) ById(id int) (*podcasts.Season, error) {
	rows, err := s.DB.Query(fmt.Sprintf("%s where id = ?;", seasonSelect), id)
	if err != nil {
		return nil, fmt.Errorf("could not query db for season by id: %v", err)
	}
	defer CloseRows(rows)

	seasons, err := parseRowsAsSeasons(rows)
	if err != nil {
		return nil, fmt.Errorf("could not parse season row: %v", err)
	}
	if len(seasons) != 1 {
		return nil, fmt.Errorf("get season by id returned %d results, but wanted 1", len(seasons))
	}
	return &seasons[0], nil
}

// ByKey retrieves a season from the store with the given key and podcast id.
func (s *SeasonStore) ByKey(key string, podcastId int) (*podcasts.Season, error) {
	rows, err := s.DB.Query(fmt.Sprintf("%s inner join podcasts as p on podcast.id = season.id where key = ? "+
		"and p.id = ?;", seasonSelect), key, podcastId)
	if err != nil {
		return nil, fmt.Errorf("could not query db for season by key %s: %v", key, err)
	}
	defer CloseRows(rows)

	seasons, err := parseRowsAsSeasons(rows)
	if err != nil {
		return nil, fmt.Errorf("could not parse season row: %v", err)
	}
	if len(seasons) != 1 {
		return nil, fmt.Errorf("get season by key returned %d results, but wanted 1", len(seasons))
	}
	return &seasons[0], nil
}

// ByPodcast retrieves all season from the store corresponding to the given podcast.
func (s *SeasonStore) ByPodcast(podcastId int) ([]podcasts.Season, error) {
	rows, err := s.DB.Query(fmt.Sprintf("%s where podcast_id = ?;", seasonSelect), podcastId)
	if err != nil {
		return nil, fmt.Errorf("could not query db for seasons by podcast id: %v", err)
	}
	defer CloseRows(rows)

	seasons, err := parseRowsAsSeasons(rows)
	if err != nil {
		return nil, fmt.Errorf("could not parse season rows: %v", err)
	}
	return seasons, nil
}

// parseRowsAsSeasons parses rows retrieved from db as seasons.
func parseRowsAsSeasons(rows *sql.Rows) ([]podcasts.Season, error) {
	var (
		id            int
		title         string
		subtitle      string
		description   string
		imageLocation string
		podcastId     int
		num           int
		key           string
	)

	var seasons []podcasts.Season
	for rows.Next() {
		err := rows.Scan(&id, &title, &subtitle, &description, &imageLocation, &podcastId, &num)
		if err != nil {
			return nil, err
		}
		seasons = append(seasons, podcasts.Season{
			Title:         title,
			Subtitle:      subtitle,
			Description:   description,
			ImageLocation: imageLocation,
			Num:           num,
			Key:           key,
		})
	}
	return seasons, nil
}