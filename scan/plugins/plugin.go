package plugins

type MetadataV1 struct {
	// Title is the title of the publication.
	Title string `msgpack:"title"`
	// Author is the author/artist of the publication.
	Author string `msgpack:"author"`
	// Language is the language of the publication.
	Language string `msgpack:"language"`
	// Issued is the date of the publication.
	Issued uint64 `msgpack:"issued"`
	// Publisher is the name of the publisher.
	Publisher string `msgpack:"publisher"`
	// Source is the path to archive file.
	Source string `msgpack:"source"`
	// Links point to the sources of the archive.
	Links []string `msgpack:"links"`
	// Categories are the categories of the archive. It can also be considered as the tags.
	Categories []CategoryV1 `msgpack:"categories"`
}

type CategoryV1 struct {
	Namespace string `msgpack:"namespace"`
	Value     string `msgpack:"value"`
}

// PluginV1 is an interface that all plugins must implement.
type PluginV1 interface {
	// Name returns the name of the plugin.
	Name() string
	// Version returns the version of the plugin.
	Version() string
	// Path returns the path to the plugin.
	Path() string
	// Execute executes the plugin and returns the metadata.
	Execute(archivePath string) (MetadataV1, error)
}
