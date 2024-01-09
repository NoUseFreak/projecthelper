package command

import (
	"net/url"
	"testing"
)

func Test_makeURL(t *testing.T) {
	tests := []struct {
		name      string
		renameMap map[string]string
		URL       string
		want      string
		wantErr   bool
	}{
		{
			name:      "nomap",
			renameMap: map[string]string{},
			URL:       "ssh://git@github.com/bla/bla",
			want:      "ssh://git@github.com/bla/bla",
			wantErr:   false,
		},
		{
			name: "notinma",
			renameMap: map[string]string{
				"bitbucket.org/bla/bla": "bb-personal",
			},
			URL:     "ssh://git@github.com/bla/bla",
			want:    "ssh://git@github.com/bla/bla",
			wantErr: false,
		},
		{
			name: "replace",
			renameMap: map[string]string{
				"github.com/bla": "gh-personal",
			},
			URL:     "ssh://git@github.com/bla/bla",
			want:    "ssh://git@gh-personal/bla/bla",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, _ := url.Parse(tt.URL)
			got, gotErr := makeURL(url, tt.renameMap)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("makeURL() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("makeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
