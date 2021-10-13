package databases

import (
	"context"
	"os"
	"testing"

	"github.com/anandkumarpatel/main/models"
	"github.com/stretchr/testify/require"
)

func Test_Integration_Postgres(t *testing.T) {
	ctx := context.Background()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		t.Skipf("Skipping integration test because DB_URL was not defined")
	}

	pg, err := NewPostgres(ctx, dbURL)
	require.NoError(t, err)

	defer pg.Close()
	db := pg.db

	s1 := &models.Service{Name: "aaa", Description: "aaa bbb ccc", Versions: []string{"1", "2"}}
	s2 := &models.Service{Name: "abc", Description: "ddd eee fff", Versions: []string{"1"}}
	s3 := &models.Service{Name: "zzz", Description: "ggg hhh iii", Versions: []string{"1", "2", "9"}}
	s4 := &models.Service{Name: "aaaa", Description: "jjj", Versions: []string{"1", "2"}}
	testServices := []*models.Service{s3, s2, s1}

	t.Run("ListServiesWithFilter", func(t *testing.T) {
		tests := []struct {
			name   string
			filter models.ServiceFilter
			db     []*models.Service
			want   []*models.Service
		}{
			{
				name:   "return all",
				filter: models.ServiceFilter{},
				db:     testServices,
				want:   testServices,
			},
			{
				name: "return filtered a",
				filter: models.ServiceFilter{
					Search: "a",
				},
				db:   testServices,
				want: []*models.Service{s2, s1},
			},
			{
				name: "return nothing if no match",
				filter: models.ServiceFilter{
					Search: "nothing",
				},
				db:   testServices,
				want: []*models.Service{},
			},
			{
				name: "sort assending",
				filter: models.ServiceFilter{
					Sort: "asc",
				},
				db:   testServices,
				want: []*models.Service{s1, s2, s3},
			},
			{
				name: "give correct page",
				filter: models.ServiceFilter{
					PageSize:   1,
					PageNumber: 1,
				},
				db:   testServices,
				want: []*models.Service{s2},
			},
		}

		for _, tt := range tests {
			db.Exec(ctx, "TRUNCATE services;")
			for _, v := range tt.db {
				db.Exec(ctx, "insert into services(name, description, versions) values($1,$2,$3)", v.Name, v.Description, v.Versions)
			}

			t.Run(tt.name, func(t *testing.T) {
				got, err := pg.ListServiesWithFilter(ctx, &tt.filter)
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("GetServiceByName", func(t *testing.T) {
		tests := []struct {
			name       string
			targetName string
			db         []*models.Service
			want       *models.Service
		}{
			{
				name:       "return correct service",
				targetName: "aaa",
				db:         []*models.Service{s3, s2, s1, s4},
				want:       s1,
			},
			{
				name:       "return nothing",
				targetName: "fake",
				db:         []*models.Service{s3, s2, s1, s4},
				want:       nil,
			},
		}

		for _, tt := range tests {
			db.Exec(ctx, "TRUNCATE services;")
			for _, v := range tt.db {
				db.Exec(ctx, "insert into services(name, description, versions) values($1,$2,$3)", v.Name, v.Description, v.Versions)
			}

			t.Run(tt.name, func(t *testing.T) {
				got, err := pg.GetServiceByName(ctx, tt.targetName)
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			})
		}
	})
}

func Test_createListQuery(t *testing.T) {
	tests := []struct {
		name   string
		filter *models.ServiceFilter
		want   string
	}{
		{"base query", &models.ServiceFilter{}, "SELECT * FROM services ORDER BY name DESC LIMIT ALL OFFSET 0"},
		{"desc sort", &models.ServiceFilter{Sort: "desc"}, "SELECT * FROM services ORDER BY name DESC LIMIT ALL OFFSET 0"},
		{"asc sort", &models.ServiceFilter{Sort: "asc"}, "SELECT * FROM services ORDER BY name ASC LIMIT ALL OFFSET 0"},
		{"set limit", &models.ServiceFilter{PageSize: 10}, "SELECT * FROM services ORDER BY name DESC LIMIT 10 OFFSET 0"},
		{"set page", &models.ServiceFilter{PageNumber: 17}, "SELECT * FROM services ORDER BY name DESC LIMIT ALL OFFSET 0"},
		{"set page", &models.ServiceFilter{PageSize: 17, PageNumber: 38}, "SELECT * FROM services ORDER BY name DESC LIMIT 17 OFFSET 646"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createListQuery(tt.filter)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_containsSearch(t *testing.T) {
	testService := &models.Service{Name: "aaa", Description: "aaa bbb ccc", Versions: []string{"1", "2"}}

	tests := []struct {
		name    string
		service *models.Service
		filter  *models.ServiceFilter
		want    bool
	}{
		{"no filter", testService, &models.ServiceFilter{}, true},
		{"name filter hit", testService, &models.ServiceFilter{Search: "a"}, true},
		{"description filter hit", testService, &models.ServiceFilter{Search: "bbb"}, true},
		{"filter miss", testService, &models.ServiceFilter{Search: "z"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsSearch(tt.service, tt.filter)
			require.Equal(t, tt.want, got)
		})
	}
}
