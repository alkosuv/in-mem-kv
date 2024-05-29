package compute

import (
	"reflect"
	"testing"
)

func TestParsingQuery(t *testing.T) {
	tests := []struct {
		name    string
		req     string
		want    []string
		wantErr bool
	}{
		{
			name: "test 1: SET request",
			req:  `SET weather_2_pm cold_moscow_weather`,
			want: []string{"SET", "weather_2_pm", "cold_moscow_weather"},
		},
		{
			name: "test 2: GET request",
			req:  `GET /etc/nginx/config.yaml`,
			want: []string{"GET", "/etc/nginx/config.yaml"},
		},
		{
			name: "test 3: DEL request",
			req:  `DEL user_****`,
			want: []string{"DEL", "user_****"},
		},
		{
			name: "test 4: SET request with JSON value",
			req:  `SET key "{ \"note\": {\"from\":\"Jessy\",\"to\": \"Joe\", \"title\": \"Water bill\", \"body\": \"Do not forget the water bill this week!\" } }"`,
			want: []string{"SET", "key", `{ \"note\": {\"from\":\"Jessy\",\"to\": \"Joe\", \"title\": \"Water bill\", \"body\": \"Do not forget the water bill this week!\" } }`},
		},
		{
			name: "test 5: SET request with JSON value",
			req:  `SET key '{ "note": {"from": "Jessy", "to": "Joe", "title": "Water bill", "body": "Do not forget the water bill this week!" } }'`,
			want: []string{"SET", "key", `{ "note": {"from": "Jessy", "to": "Joe", "title": "Water bill", "body": "Do not forget the water bill this week!" } }`},
		},
		{
			name: "test 6: SET request with JSON value",
			req:  `SET key '{ "note": {"from": "Jessy", "to": "Joe", "title": "Water bill", "body": "Do not forget the water bill this week!" } }'`,
			want: []string{"SET", "key", `{ "note": {"from": "Jessy", "to": "Joe", "title": "Water bill", "body": "Do not forget the water bill this week!" } }`},
		},
		{
			name: "test 7: SET request with new line",
			req:  "SET key 'value'\n",
			want: []string{"SET", "key", "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsingQuery(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Processing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Processing() = %v, want %v", got, tt.want)
			}
		})
	}
}
