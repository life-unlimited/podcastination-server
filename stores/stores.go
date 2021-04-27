package stores

import (
	"database/sql"
	"log"
)

type Stores struct {
	Podcasts PodcastStore
	Owners   OwnerStore
	Seasons  SeasonStore
	Episodes EpisodeStore
}

func CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Fatalf("could not close rows: %v", err)
	}
}
