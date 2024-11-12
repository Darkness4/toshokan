package main

import (
	"testing"

	_ "embed"

	"github.com/Darkness4/toshokan/scan/plugins"
	"github.com/stretchr/testify/require"
)

func TestParseMetadata(t *testing.T) {
	tt := []struct {
		filePath string
		expected plugins.MetadataV1
	}{
		{
			filePath: "fixtures/anchira.yaml",
			expected: plugins.MetadataV1{
				Title:          "Title",
				Issued:         1657469676,
				ThumbnailIndex: 1,
				Categories: []plugins.CategoryV1{
					{
						Namespace: "",
						Value:     "A",
					},
					{
						Namespace: "",
						Value:     "B",
					},
					{
						Namespace: "",
						Value:     "C",
					},
					{
						Namespace: "pages",
						Value:     "42",
					},
					{
						Namespace: "artist",
						Value:     "Artist",
					},
					{
						Namespace: "date_released",
						Value:     "1657469676",
					},
					{
						Namespace: "language",
						Value:     "english",
					},
				},
			},
		},
		{
			filePath: "fixtures/koharu1.yaml",
			expected: plugins.MetadataV1{
				Title:          "Title",
				Issued:         1722812410,
				ThumbnailIndex: 1,
				Categories: []plugins.CategoryV1{
					{
						Namespace: "",
						Value:     "A",
					},
					{
						Namespace: "",
						Value:     "B",
					},
					{
						Namespace: "",
						Value:     "C",
					},
					{
						Namespace: "pages",
						Value:     "18",
					},
					{
						Namespace: "artist",
						Value:     "Artist",
					},
					{
						Namespace: "magazine",
						Value:     "Magazine",
					},
					{
						Namespace: "publisher",
						Value:     "Publisher",
					},
					{
						Namespace: "parody",
						Value:     "Original Work",
					},
					{
						Namespace: "date_released",
						Value:     "1722812410",
					},
					{
						Namespace: "language",
						Value:     "english",
					},
				},
			},
		},
		{
			filePath: "fixtures/koharu2.yaml",
			expected: plugins.MetadataV1{
				Title:          "title",
				Issued:         0,
				ThumbnailIndex: 0,
				Categories: []plugins.CategoryV1{
					{
						Namespace: "male",
						Value:     "a",
					},
					{
						Namespace: "male",
						Value:     "b",
					},
					{
						Namespace: "male",
						Value:     "c",
					},
					{
						Namespace: "female",
						Value:     "a",
					},
					{
						Namespace: "female",
						Value:     "b",
					},
					{
						Namespace: "female",
						Value:     "c",
					},
					{
						Namespace: "mixed",
						Value:     "a",
					},
					{
						Namespace: "mixed",
						Value:     "b",
					},
					{
						Namespace: "mixed",
						Value:     "c",
					},
					{
						Namespace: "other",
						Value:     "a",
					},
					{
						Namespace: "artist",
						Value:     "artist",
					},
					{
						Namespace: "circle",
						Value:     "circle",
					},
					{
						Namespace: "parody",
						Value:     "original",
					},
					{
						Namespace: "language",
						Value:     "english",
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.filePath, func(t *testing.T) {
			meta, err := parseMetadata(tc.filePath)
			if err != nil {
				t.Fatal(err)
			}

			require.Equal(t, tc.expected.Title, meta.Title)
			require.Equal(t, tc.expected.Issued, meta.Issued)
			require.Equal(t, tc.expected.ThumbnailIndex, meta.ThumbnailIndex)
			require.Equal(t, len(tc.expected.Categories), len(meta.Categories))
			for _, category := range tc.expected.Categories {
				require.Contains(t, meta.Categories, category)
			}
		})
	}
}
