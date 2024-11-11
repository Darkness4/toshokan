package golang

import (
	"path/filepath"
	"plugin"

	"github.com/Darkness4/toshokan/scan/plugins"
	"github.com/shamaton/msgpack/v2"
)

var _ plugins.PluginV1 = (*Plugin)(nil)

type Plugin struct {
	// Path is the path to the Go plugin.
	path   string
	name   string
	plugin *plugin.Plugin

	versionFunc func() string
	executeFunc func(archivePath string) ([]byte, error)
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
	executeFunc, ok := executeFunc.(func(archivePath string) ([]byte, error))
	if !ok {
		panic("Execute function is not of type func(archivePath string) ([]byte, error)")
	}

	versionFunc, err := plugin.Lookup("Version")
	if err != nil {
		panic(err)
	}
	versionFunc, ok = versionFunc.(func() string)
	if !ok {
		panic("Version function is not of type func() string")
	}

	name := filepath.Base(path)
	// Remove the extension. File should end with .so.
	name = name[:len(name)-3]
	return &Plugin{
		path:   path,
		name:   name,
		plugin: plugin,

		versionFunc: versionFunc.(func() string),
		executeFunc: executeFunc.(func(archivePath string) ([]byte, error)),
	}
}

// Name is the name of the plugin.
func (p Plugin) Name() string {
	return p.name
}

// Version is the version of the plugin.
func (p Plugin) Version() string {
	return p.versionFunc()
}

// Path is the path to the plugin.
func (p Plugin) Path() string {
	return p.path
}

// Execute executes the plugin and returns the metadata.
func (p Plugin) Execute(archivePath string) (plugins.MetadataV1, error) {
	b, err := p.executeFunc(archivePath)
	if err != nil {
		return plugins.MetadataV1{}, err
	}

	var metadata plugins.MetadataV1
	if err := msgpack.Unmarshal(b, &metadata); err != nil {
		return plugins.MetadataV1{}, err
	}
	return metadata, err
}
