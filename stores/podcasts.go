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

const podcastSelect = "select id, title, subtitle, language, owner_id, description, keywords, link, image_location, type, key, feed_link from podcasts"

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
func (s *PodcastStore) ById(id int) (podcasts.Podcast, error) {
	rows, err := s.DB.Query(fmt.Sprintf("%s where id = $1", podcastSelect), id)
	if err != nil {
		return podcasts.Podcast{}, fmt.Errorf("could not query db for podcast by id %v: %v", id, err)
	}
	defer CloseRows(rows)

	pcs, err := parseRowsAsPodcasts(rows)
	if err != nil {
		return podcasts.Podcast{}, fmt.Errorf("error while parsing podcast row: %v", err)
	}
	if len(pcs) != 1 {
		return podcasts.Podcast{}, fmt.Errorf("get podcast by id from db returned %v results, but wanted 1", len(pcs))
	}
	return pcs[0], nil
}

// ByKey retrieves a podcast from the store with the given key.
func (s *PodcastStore) ByKey(key string) (podcasts.Podcast, error) {
	rows, err := s.DB.Query(fmt.Sprintf("%s where key = $1", podcastSelect), key)
	if err != nil {
		return podcasts.Podcast{}, fmt.Errorf("could not query db for podcast by key %s: %v", key, err)
	}
	defer CloseRows(rows)

	pcs, err := parseRowsAsPodcasts(rows)
	if err != nil {
		return podcasts.Podcast{}, fmt.Errorf("error while parsing podcast row: %v", err)
	}
	if len(pcs) != 1 {
		return podcasts.Podcast{}, fmt.Errorf("get podcast by key from db returned %v results, but wanted 1", len(pcs))
	}
	return pcs[0], nil
}

// parseRowsAsPodcasts parses rows retrieved from db as podcasts.
func parseRowsAsPodcasts(rows *sql.Rows) ([]podcasts.Podcast, error) {
	var (
		id            int
		title         string
		subtitle      sql.NullString
		language      podcasts.Language
		ownerId       int
		description   sql.NullString
		keywords      sql.NullString
		link          sql.NullString
		imageLocation sql.NullString
		podcastType   sql.NullString
		key           sql.NullString
		feedLink      string
	)

	pcs := make([]podcasts.Podcast, 0)
	for rows.Next() {
		err := rows.Scan(&id, &title, &subtitle, &language, &ownerId, &description, &keywords, &link, &imageLocation,
			&podcastType, &key, &feedLink)
		if err != nil {
			return nil, err
		}
		pcs = append(pcs, podcasts.Podcast{
			Id:            id,
			Title:         title,
			Subtitle:      subtitle.String,
			Language:      language,
			OwnerId:       ownerId,
			Description:   description.String,
			Keywords:      strings.Split(keywords.String, ","),
			Link:          link.String,
			ImageLocation: imageLocation.String,
			PodcastType:   podcasts.PodcastType(podcastType.String),
			Key:           key.String,
			FeedLink:      feedLink,
		})
	}
	return pcs, nil
}
