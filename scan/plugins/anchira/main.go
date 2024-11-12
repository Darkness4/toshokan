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

// StringOrArray represents a value that can be either a single string or an array of strings.
type StringOrArray []string

// UnmarshalYAML is a custom unmarshaller that handles both single string and array of strings.
func (s *StringOrArray) UnmarshalYAML(value *yaml.Node) error {
	// Try unmarshalling as a single string
	var single string
	if err := value.Decode(&single); err == nil {
		*s = []string{single}
		return nil
	}

	// Try unmarshalling as an array of strings
	var array []string
	if err := value.Decode(&array); err == nil {
		*s = array
		return nil
	}

	// If it doesn't match either, return an error
	return fmt.Errorf("value must be a string or an array of strings")
}

type Metadata struct {
	Title          string        `yaml:"Title"`
	Artist         StringOrArray `yaml:"Artist"`
	Circle         StringOrArray `yaml:"Circle"`
	Parody         StringOrArray `yaml:"Parody"`
	Publisher      StringOrArray `yaml:"Publisher"`
	Magazine       StringOrArray `yaml:"Magazine"`
	Tags           StringOrArray `yaml:"Tags"`
	Released       int64         `yaml:"Released"`
	Pages          uint64        `yaml:"Pages"`
	Thumbnail      int64         `yaml:"Thumbnail"`
	ThumbnailIndex int64         `yaml:"ThumbnailIndex"`

	KoushokuMetadata KoharuMetadata `yaml:",inline"`
}

// KoharuMetadata represents the legacy metadata for Koharu.
type KoharuMetadata struct {
	Title    string        `yaml:"title"`
	General  StringOrArray `yaml:"general"`
	Male     StringOrArray `yaml:"male"`
	Female   StringOrArray `yaml:"female"`
	Mixed    StringOrArray `yaml:"mixed"`
	Other    StringOrArray `yaml:"other"`
	Artist   StringOrArray `yaml:"artist"`
	Circle   StringOrArray `yaml:"circle"`
	Parody   StringOrArray `yaml:"parody"`
	Magazine StringOrArray `yaml:"magazine"`
	Language StringOrArray `yaml:"language"`
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
	metadataPathInArchive, _, err := archive.FindFile(archivePath, "koushoku.yaml")
	if metadataPathInArchive == "" || err != nil {
		metadataPathInArchive, _, err = archive.FindFile(archivePath, "info.yaml")
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

	return parseMetadataYAML(metadata)
}

func parseMetadataYAML(metadata Metadata) (plugins.MetadataV1, error) {
	tags := make(map[plugins.CategoryV1]struct{})
	appendTags(tags, "", metadata.Tags...)
	appendTags(tags, "artist", metadata.Artist...)
	appendTags(tags, "circle", metadata.Circle...)
	appendTags(tags, "parody", metadata.Parody...)
	appendTags(tags, "magazine", metadata.Magazine...)
	appendTags(tags, "publisher", metadata.Publisher...)
	if metadata.Pages > 0 {
		appendTags(tags, "pages", strconv.FormatUint(metadata.Pages, 10))
	}

	// Handle koharu tags
	title := metadata.Title
	if title == "" {
		title = metadata.KoushokuMetadata.Title
	}
	appendTags(tags, "general", metadata.KoushokuMetadata.General...)
	appendTags(tags, "male", metadata.KoushokuMetadata.Male...)
	appendTags(tags, "female", metadata.KoushokuMetadata.Female...)
	appendTags(tags, "mixed", metadata.KoushokuMetadata.Mixed...)
	appendTags(tags, "other", metadata.KoushokuMetadata.Other...)
	appendTags(tags, "artist", metadata.KoushokuMetadata.Artist...)
	appendTags(tags, "circle", metadata.KoushokuMetadata.Circle...)
	appendTags(tags, "parody", metadata.KoushokuMetadata.Parody...)
	appendTags(tags, "magazine", metadata.KoushokuMetadata.Magazine...)
	appendTags(tags, "language", metadata.KoushokuMetadata.Language...)

	if len(metadata.KoushokuMetadata.Language) == 0 {
		tags[plugins.CategoryV1{
			Namespace: "language",
			Value:     "english",
		}] = struct{}{}
	}

	if metadata.Released > 0 {
		tags[plugins.CategoryV1{
			Namespace: "date_released",
			Value:     strconv.FormatInt(int64(metadata.Released), 10),
		}] = struct{}{}
	}

	thumbnailIndex := metadata.ThumbnailIndex
	if thumbnailIndex == 0 {
		thumbnailIndex = metadata.Thumbnail
	}

	categories := make([]plugins.CategoryV1, 0, len(tags))
	for category := range tags {
		categories = append(categories, category)
	}

	return plugins.MetadataV1{
		Title:          title,
		Issued:         metadata.Released,
		Categories:     categories,
		ThumbnailIndex: thumbnailIndex,
	}, nil
}

func appendTags(tags map[plugins.CategoryV1]struct{}, namespace string, values ...string) {
	for _, value := range values {
		if value == "" {
			continue
		}
		tags[plugins.CategoryV1{
			Namespace: namespace,
			Value:     value,
		}] = struct{}{}
	}
}
