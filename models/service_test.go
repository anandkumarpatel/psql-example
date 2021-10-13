package models

import "testing"

func TestServiceFilter_Validate(t *testing.T) {
	type fields struct {
		Search     string
		Sort       string
		PageSize   int
		PageNumber int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"should pass", fields{}, false},
		{"valid sort", fields{Sort: "asc"}, false},
		{"valid sort", fields{Sort: "desc"}, false},
		{"valid search", fields{Search: "anandkumarpatel"}, false},
		{"invalid sort", fields{Sort: "wrong"}, true},
		{"invalid search", fields{Search: "drop;table"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServiceFilter{
				Search:     tt.fields.Search,
				Sort:       tt.fields.Sort,
				PageSize:   tt.fields.PageSize,
				PageNumber: tt.fields.PageNumber,
			}
			if err := s.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ServiceFilter.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
