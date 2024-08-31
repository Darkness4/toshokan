package golang

import (
	"path/filepath"
	"plugin"
	"strings"

	"github.com/Darkness4/toshokan/scan/plugins"
)

var _ plugins.Plugin = (*Plugin)(nil)

type Plugin struct {
	// Path is the path to the Go plugin.
	path        string
	name        string
	version     string
	plugin      *plugin.Plugin
	executeFunc func(archivePath string) (plugins.Metadata, error)
}

func NewPlugin(path string) *Plugin {
	plugin, err := plugin.Open(path)
	if err != nil {
		panic(err)
	}

	executeFunc, err := plugin.Lookup("Execute")
	if err != nil {
		panic(err)
	}

	name := filepath.Base(path)
	// Remove the extension. File should end with .so.
	name = name[:len(name)-3]
	// Remove the version. File should end with -X.Y.Z.so.
	name, version, _ := lastCut(name, "-")
	return &Plugin{
		path:        path,
		name:        name,
		version:     version,
		plugin:      plugin,
		executeFunc: executeFunc.(func(archivePath string) (plugins.Metadata, error)),
	}
}

// Name is the name of the plugin.
func (p Plugin) Name() string {
	return p.name
}

// Version is the version of the plugin.
func (p Plugin) Version() string {
	return p.version
}

// Path is the path to the plugin.
func (p Plugin) Path() string {
	return p.path
}

// Execute executes the plugin and returns the metadata.
func (p Plugin) Execute(archivePath string) (plugins.Metadata, error) {
	return p.executeFunc(archivePath)
}

func lastCut(s, sep string) (before, after string, found bool) {
	if i := strings.LastIndex(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}
