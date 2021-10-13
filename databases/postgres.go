package databases

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/anandkumarpatel/main/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

// Postgres implements the ServiceModel interface backed by a postgres database.
type Postgres struct {
	db *pgxpool.Pool
}

// NewPostgres creates a new instance of Postgres.
func NewPostgres(ctx context.Context, dbURL string) (*Postgres, error) {
	db, err := pgxpool.Connect(ctx, dbURL)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to db: %s", dbURL)
	}

	return &Postgres{
		db: db,
	}, nil
}

// Close the the connection to the db.
func (s *Postgres) Close() {
	s.db.Close()
}

// ListServiesWithFilter returns all services that match passed filter.
func (s *Postgres) ListServiesWithFilter(ctx context.Context, filter *models.ServiceFilter) ([]*models.Service, error) {
	query := createListQuery(filter)
	rows, _ := s.db.Query(ctx, query)

	services := []*models.Service{}
	for rows.Next() {
		service := &models.Service{}
		err := rows.Scan(&service.Name, &service.Description, &service.Versions)

		if err != nil {
			return nil, err
		}

		if containsSearch(service, filter) {
			services = append(services, service)
		}
	}

	return services, rows.Err()
}

func createListQuery(filter *models.ServiceFilter) string {
	sort := "DESC"
	if filter.Sort == "asc" {
		sort = "ASC"
	}

	limit := "ALL"
	if filter.PageSize != 0 {
		limit = strconv.Itoa(filter.PageSize)
	}

	offset := "0"
	if filter.PageNumber != 0 {
		offset = strconv.Itoa(filter.PageNumber * filter.PageSize)
	}

	return fmt.Sprintf("SELECT * FROM services ORDER BY name %s LIMIT %s OFFSET %s", sort, limit, offset)
}

func containsSearch(service *models.Service, filter *models.ServiceFilter) bool {
	if filter.Search == "" {
		return true
	}

	if strings.Contains(service.Name, filter.Search) {
		return true
	}

	if strings.Contains(service.Description, filter.Search) {
		return true
	}

	return false
}

// GetServiceByName returns the service with passed name.
func (s *Postgres) GetServiceByName(ctx context.Context, name string) (*models.Service, error) {
	services, err := s.ListServiesWithFilter(ctx, &models.ServiceFilter{Search: name})
	if err != nil {
		return nil, err
	}

	for _, v := range services {
		if v.Name == name {
			return v, nil
		}
	}

	return nil, nil
}
