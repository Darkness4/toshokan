package plugins

type MetadataV1 struct {
	// Title is the title of the publication.
	Title string `msgpack:"title"`
	// Issued is the date of the publication.
	Issued int64 `msgpack:"issued"`
	// Categories are the categories of the archive. It can also be considered as the tags.
	Categories []CategoryV1 `msgpack:"categories"`
	// Index of the thumbnail in the archive.
	ThumbnailIndex int64 `msgpack:"thumbnail_index"`
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
