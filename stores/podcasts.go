package stores

import (
	"database/sql"
	"fmt"
	"life-unlimited/podcastination/podcasts"
	"strings"
)

type PodcastStore struct {
	DB *sql.DB
}

const podcastSelect = "select id, title, subtitle, language, owner_id, description, keywords, link, image_location, type, key from podcasts"

// All retrieves all podcasts from the store.
func (s *PodcastStore) All() ([]podcasts.Podcast, error) {
	rows, err := s.DB.Query(fmt.Sprintf("%s;", podcastSelect))
	if err != nil {
		return nil, fmt.Errorf("could not query db for podcasts: %v", err)
	}
	defer CloseRows(rows)

	pcs, err := parseRowsAsPodcasts(rows)
	if err != nil {
		return nil, fmt.Errorf("error while parsing podcast rows: %v", err)
	}
	return pcs, nil
}

// ById retrieves a podcast from the store with the given id.
func (s *PodcastStore) ById(id int) (*podcasts.Podcast, error) {
	rows, err := s.DB.Query(fmt.Sprintf("%s where id = ?", podcastSelect), id)
	if err != nil {
		return nil, fmt.Errorf("could not query db for podcast by id %v: %v", id, err)
	}
	defer CloseRows(rows)

	pcs, err := parseRowsAsPodcasts(rows)
	if err != nil {
		return nil, fmt.Errorf("error while parsing podcast row: %v", err)
	}
	if len(pcs) != 1 {
		return nil, fmt.Errorf("get podcast by id from DB returned %v results, but wanted 1", len(pcs))
	}
	return &pcs[1], nil
}

// ByKey retrieves a podcast from the store with the given key.
func (s *PodcastStore) ByKey(key string) (*podcasts.Podcast, error) {
	rows, err := s.DB.Query(fmt.Sprintf("%s where key = ?", podcastSelect), key)
	if err != nil {
		return nil, fmt.Errorf("could not query db for podcast by key %s: %v", key, err)
	}
	defer CloseRows(rows)

	pcs, err := parseRowsAsPodcasts(rows)
	if err != nil {
		return nil, fmt.Errorf("error while parsing podcast row: %v", err)
	}
	if len(pcs) != 1 {
		return nil, fmt.Errorf("get podcast by id from DB returned %v results, but wanted 1", len(pcs))
	}
	return &pcs[1], nil
}

// parseRowsAsPodcasts parses rows retrieved from db as podcasts.
func parseRowsAsPodcasts(rows *sql.Rows) ([]podcasts.Podcast, error) {
	var (
		id            int
		title         string
		subtitle      string
		language      podcasts.Language
		ownerId       int
		description   string
		keywords      string
		link          string
		imageLocation string
		podcastType   podcasts.PodcastType
		key           string
	)

	var pcs []podcasts.Podcast
	for rows.Next() {
		err := rows.Scan(&id, &title, &subtitle, &language, &ownerId, &description, &keywords, &link, &imageLocation, &podcastType)
		if err != nil {
			return nil, err
		}
		pcs = append(pcs, podcasts.Podcast{
			Id:            id,
			Title:         title,
			Subtitle:      subtitle,
			Language:      language,
			OwnerId:       ownerId,
			Description:   description,
			Keywords:      strings.Split(keywords, ","),
			Link:          link,
			ImageLocation: imageLocation,
			PodcastType:   podcastType,
			Key:           key,
		})
	}
	return pcs, nil
}
