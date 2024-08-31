package lua

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Darkness4/toshokan/scan/plugins"
	"github.com/Darkness4/toshokan/scan/plugins/lua/engine"
	"github.com/shamaton/msgpack/v2"
)

var _ plugins.Plugin = (*Plugin)(nil)

type Plugin struct {
	// Path is the path to the Lua script.
	path    string
	name    string
	version string
}

func NewPlugin(path string) *Plugin {
	// Check if the file is a Lua script.
	if !engine.IsLuaScript(path) {
		panic(fmt.Sprintf("file %s is not a Lua script", path))
	}

	name := filepath.Base(path)
	// Remove the extension. File should end with .lua.
	name = name[:len(name)-4]
	// Remove the version. File should end with -X.Y.Z.lua.
	name, version, _ := lastCut(name, "-")
	return &Plugin{
		path:    path,
		name:    name,
		version: version,
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
	result, err := engine.ExecuteFromPath(p.path, archivePath)
	if err != nil {
		return plugins.Metadata{}, err
	}
	// Parse the result
	var metadata plugins.Metadata
	if err := msgpack.Unmarshal([]byte(result), &metadata); err != nil {
		return plugins.Metadata{}, fmt.Errorf("failed to unmarshal message pack data: %v", err)
	}
	return metadata, nil
}

func lastCut(s, sep string) (before, after string, found bool) {
	if i := strings.LastIndex(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}
