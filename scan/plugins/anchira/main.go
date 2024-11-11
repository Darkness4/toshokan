package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Darkness4/toshokan/archive"
	"github.com/Darkness4/toshokan/scan/plugins"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shamaton/msgpack/v2"
	"gopkg.in/yaml.v3"
)

type Metadata struct {
	Source    string   `yaml:"Source"`
	URL       string   `yaml:"URL"`
	Title     string   `yaml:"Title"`
	Artist    []string `yaml:"Artist"`
	Circle    []string `yaml:"Circle"`
	Parody    []string `yaml:"Parody"`
	Magazine  []string `yaml:"Magazine"`
	Tags      []string `yaml:"Tags"`
	Released  uint64   `yaml:"Released"`
	Pages     uint64   `yaml:"Pages"`
	Thumbnail uint64   `yaml:"Thumbnail"`

	KoushokuMetadata KoharuMetadata `yaml:",inline"`
}

type KoharuMetadata struct {
	Title    string `yaml:"title"`
	General  string `yaml:"general"`
	Male     string `yaml:"male"`
	Female   string `yaml:"female"`
	Mixed    string `yaml:"mixed"`
	Other    string `yaml:"other"`
	Artist   string `yaml:"artist"`
	Circle   string `yaml:"circle"`
	Parody   string `yaml:"parody"`
	Magazine string `yaml:"magazine"`
	Language string `yaml:"language"`
	Source   string `yaml:"source"`
}

var version = ""

func init() {
	// Allow better for logging in production.
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
}

func Version() string {
	return version
}

func Execute(archivePath string) ([]byte, error) {
	log := log.With().Str("plugin", "anchira").Str("archivePath", archivePath).Logger()
	metadataPathInArchive, err := archive.FindFile(archivePath, "koushoku.yaml")
	if metadataPathInArchive == "" || err != nil {
		metadataPathInArchive, err = archive.FindFile(archivePath, "info.yaml")
	}
	if err != nil {
		log.Err(err).Msg("failed to find koushoku.yaml or info.yaml")
		return nil, err
	}

	if metadataPathInArchive == "" {
		log.Error().Msg("koushoku.yaml or info.yaml not found")
		return nil, fmt.Errorf("koushoku.yaml or info.yaml not found")
	}

	tmpDir, err := os.MkdirTemp("", "anchira")
	if err != nil {
		log.Error().Err(err).Msg("failed to create temporary directory")
		return nil, fmt.Errorf("failed to create temporary directory")
	}
	defer os.RemoveAll(tmpDir)

	outFilePath := filepath.Join(tmpDir, "metadata.yaml")

	if err = archive.ExtractFile(archivePath, metadataPathInArchive, outFilePath); err != nil {
		log.Error().Err(err).Str("path", metadataPathInArchive).Msg("failed to extract file")
		return nil, fmt.Errorf("failed to extract file")
	}

	metadata, err := parseMetadata(outFilePath)
	if err != nil {
		log.Error().Err(err).Str("path", outFilePath).Msg("failed to parse metadata")
		return nil, fmt.Errorf("failed to parse metadata")
	}

	return msgpack.Marshal(metadata)
}

func parseMetadata(path string) (plugins.MetadataV1, error) {
	var metadata Metadata

	file, err := os.Open(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("failed to open file")
		return plugins.MetadataV1{}, fmt.Errorf("failed to open file")
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&metadata); err != nil {
		log.Error().Err(err).Str("path", path).Msg("failed to decode YAML")
		return plugins.MetadataV1{}, fmt.Errorf("failed to decode YAML")
	}

	tags := make(map[plugins.CategoryV1]struct{})
	appendTags(tags, "", metadata.Tags...)
	appendTags(tags, "artist", metadata.Artist...)
	appendTags(tags, "series", metadata.Parody...)
	appendTags(tags, "magazine", metadata.Magazine...)

	// Handle koharu tags
	title := metadata.Title
	if title == "" {
		title = metadata.KoushokuMetadata.Title
	}
	appendTags(tags, "general", metadata.KoushokuMetadata.General)
	appendTags(tags, "male", metadata.KoushokuMetadata.Male)
	appendTags(tags, "female", metadata.KoushokuMetadata.Female)
	appendTags(tags, "mixed", metadata.KoushokuMetadata.Mixed)
	appendTags(tags, "other", metadata.KoushokuMetadata.Other)
	appendTags(tags, "artist", metadata.KoushokuMetadata.Artist)
	appendTags(tags, "circle", metadata.KoushokuMetadata.Circle)
	appendTags(tags, "parody", metadata.KoushokuMetadata.Parody)
	appendTags(tags, "magazine", metadata.KoushokuMetadata.Magazine)
	appendTags(tags, "language", metadata.KoushokuMetadata.Language)
	appendTags(tags, "source", metadata.KoushokuMetadata.Source)

	tags[plugins.CategoryV1{
		Namespace: "language",
		Value:     "english",
	}] = struct{}{}

	tags[plugins.CategoryV1{
		Namespace: "date_released",
		Value:     strconv.FormatInt(int64(metadata.Released), 10),
	}] = struct{}{}

	tags[plugins.CategoryV1{
		Namespace: "source",
		Value:     metadata.Source,
	}] = struct{}{}

	categories := make([]plugins.CategoryV1, 0, len(tags))
	for category := range tags {
		categories = append(categories, category)
	}

	return plugins.MetadataV1{
		Title:      title,
		Author:     metadata.Artist[0],
		Language:   "Japanese",
		Issued:     metadata.Released,
		Publisher:  metadata.Magazine[0],
		Source:     metadata.Source,
		Links:      []string{metadata.URL},
		Categories: categories,
	}, nil
}

func appendTags(tags map[plugins.CategoryV1]struct{}, namespace string, values ...string) {
	for _, value := range values {
		tags[plugins.CategoryV1{
			Namespace: namespace,
			Value:     value,
		}] = struct{}{}
	}
}
