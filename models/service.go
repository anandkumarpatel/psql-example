package models

import (
	"context"
	"fmt"
	"regexp"
)

// Service is the model for a service.
type Service struct {
	Name        string
	Description string
	Versions    []string
}

// ServiceModel is the interface for the `Service` model.
type ServiceModel interface {
	ListServiesWithFilter(ctx context.Context, filter *ServiceFilter) ([]*Service, error)
	GetServiceByName(ctx context.Context, name string) (*Service, error)
}

// ServiceFilter holds filtering request by the user.
type ServiceFilter struct {
	Search string

	// Sort applies to the Name. it can be asc, desc, or "".
	Sort string

	// The following params are for pagination.
	PageSize   int
	PageNumber int
}

var IsAlphaNum = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

// Validate checks the filter to ensure it is safe and correct.
func (s *ServiceFilter) Validate() error {
	if s.Search != "" && !IsAlphaNum(s.Search) {
		return fmt.Errorf("search contains invalid charactor. Only [a-zA-Z0-9] are allowed: %s", s.Search)
	}

	if s.Sort != "" && (s.Sort != "asc" && s.Sort != "desc") {
		return fmt.Errorf("sort contains invalid value. Only asc, desc, or '' is allowed: %s", s.Sort)
	}

	return nil
}
