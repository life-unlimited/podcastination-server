package stores

import (
	"database/sql"
	"log"
)

func CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Fatalf("could not close rows: %v", err)
	}
}
