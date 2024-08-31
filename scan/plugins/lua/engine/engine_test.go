package engine_test

import (
	"testing"

	"github.com/Darkness4/toshokan/scan/plugins/lua/engine"
	"github.com/shamaton/msgpack/v2"
)

type ExampleData struct {
	Title       string   `msgpack:"title"`
	ReleaseDate int64    `msgpack:"release_date"`
	Tags        []string `msgpack:"tags"`
}

var fixtureArchivePath = "fixtures/archive.zip"

func TestExecuteFromPath(t *testing.T) {
	tt := []struct {
		path   string
		assert func(*testing.T, string)
	}{
		{
			path: "fixtures/helloworld.lua",
		},
		{
			path: "fixtures/fetching_data.lua",
			assert: func(t *testing.T, res string) {
				var data ExampleData
				err := msgpack.Unmarshal([]byte(res), &data)
				if err != nil {
					t.Fatalf("failed to unmarshal message pack data: %v", err)
				}

				if data.Title != "Example Title" {
					t.Fatalf("unexpected title: %s", data.Title)
				}

				if data.ReleaseDate != 1682390400 {
					t.Fatalf("unexpected release date: %d", data.ReleaseDate)
				}

				if len(data.Tags) != 3 {
					t.Fatalf("unexpected tags: %v", data.Tags)
				}
			},
		},
		{
			path: "fixtures/params.lua",
			assert: func(t *testing.T, res string) {
				if res != fixtureArchivePath {
					t.Fatalf("unexpected result: %s", res)
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.path, func(t *testing.T) {
			res, err := engine.ExecuteFromPath(tc.path, fixtureArchivePath)
			if err != nil {
				t.Fatalf("failed to execute Lua script from file: %v", err)
			}

			if tc.assert != nil {
				tc.assert(t, res)
			}
		})
	}
}
