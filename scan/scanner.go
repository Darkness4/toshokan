// Package scan provides a scanner for scanning archives into metadatas.
package scan

import (
	"io/fs"

	"github.com/Darkness4/toshokan/archive"
	"github.com/Darkness4/toshokan/scan/plugins"
	"github.com/rs/zerolog/log"
)

type Scanner struct {
	fsys    fs.FS
	Plugins []plugins.PluginV1
}

// NewScanner creates a new scanner with the given plugins.
func NewScanner(plugins ...plugins.PluginV1) *Scanner {
	return &Scanner{Plugins: plugins}
}

// Init prints the version of the plugins.
func (s *Scanner) Init() {
	for _, plugin := range s.Plugins {
		log.Info().Str("plugin", plugin.Name()).Str("version", plugin.Version()).Msg("loaded plugin")
	}
}

type ScanResult struct {
	Meta plugins.MetadataV1
	Path string
}

// Scan scans the archive at the given path and returns the metadata.
func (s *Scanner) Scan() func(func(ScanResult) bool) {
	return func(yield func(ScanResult) bool) {
		for entry := range WalkDirIter(s.fsys, ".") {
			if entry.Err != nil {
				continue
			}
			if !archive.IsSupported(entry.Path) {
				continue
			}
			var meta plugins.MetadataV1
			var err error
			for _, plugin := range s.Plugins {
				meta, err = plugin.Execute(entry.Path)
				if err != nil {
					log.Err(err).Msg("failed to execute plugin")
					continue
				}

				// TODO: merge metadata
			}
			if !yield(ScanResult{
				Meta: meta,
				Path: entry.Path,
			}) {
				return
			}
		}
	}
}

type WalkDirEntry struct {
	Path    string
	Entry   fs.DirEntry
	Err     error
	skipDir *bool
}

// SkipDir causes the iteration to skip the contents
// of the entry. This will have no effect if called outside the iteration
// for this particular entry.
func (entry WalkDirEntry) SkipDir() {
	*entry.skipDir = true
}

// WalkDirIter returns an iterator that can be used to iterate over the contents
// of a directory. It uses WalkDir under the hood. To skip a directory,
// call the SkipDir method on an entry.
func WalkDirIter(fsys fs.FS, root string) func(func(WalkDirEntry) bool) {
	return func(yield func(WalkDirEntry) bool) {
		fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
			skipDir := false
			info := WalkDirEntry{
				Path:    path,
				Entry:   d,
				Err:     err,
				skipDir: &skipDir,
			}
			if !yield(info) {
				return fs.SkipAll
			}
			if skipDir {
				return fs.SkipDir
			}
			return nil
		})
	}
}
