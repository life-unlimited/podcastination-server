package stores

import (
	"database/sql"
	"fmt"
	"life-unlimited/podcastination/podcasts"
)

const ownerSelect = "select id, name, email, copyright from owners"

type OwnerStore struct {
	DB *sql.DB
}

// All retrieves all owners from the store.
func (s *OwnerStore) All() ([]podcasts.Owner, error) {
	rows, err := s.DB.Query(ownerSelect)
	if err != nil {
		return nil, fmt.Errorf("could not query db for owners: %v", err)
	}
	defer CloseRows(rows)

	owners, err := parseRowsAsOwners(rows)
	if err != nil {
		return nil, fmt.Errorf("could not parse owner rows: %v", err)
	}
	return owners, nil
}

// ById retrieves an owner from the store by id.
func (s *OwnerStore) ById(id int) (podcasts.Owner, error) {
	rows, err := s.DB.Query(fmt.Sprintf("%s where id = $1", ownerSelect), id)
	if err != nil {
		return podcasts.Owner{}, fmt.Errorf("could not query db for owner by id: %v", err)
	}
	defer CloseRows(rows)

	owners, err := parseRowsAsOwners(rows)
	if err != nil {
		return podcasts.Owner{}, fmt.Errorf("could not parse owner row: %v", err)
	}
	if len(owners) != 1 {
		return podcasts.Owner{}, fmt.Errorf("get owner by id returned %d results, but wanted 1", len(owners))
	}
	return owners[0], nil
}

// parseRowsAsOwners parses rows retrieved from db as owners.
func parseRowsAsOwners(rows *sql.Rows) ([]podcasts.Owner, error) {
	var (
		id        int
		name      string
		email     string
		copyright sql.NullString
	)

	var owners []podcasts.Owner
	for rows.Next() {
		err := rows.Scan(&id, &name, &email, &copyright)
		if err != nil {
			return nil, err
		}
		owners = append(owners, podcasts.Owner{
			Id:        id,
			Name:      name,
			Email:     email,
			Copyright: copyright.String,
		})
	}
	return owners, nil
}
