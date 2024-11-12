-- +goose up
CREATE TABLE IF NOT EXISTS tags ( -- atom:summary
  id BIGSERIAL PRIMARY KEY,
  namespace TEXT NOT NULL, -- used for filtering
  value TEXT NOT NULL, -- used for sorting
  UNIQUE (namespace, value)
);

-- Known namespaces:
-- language, which is used by dc:language
-- publisher, which is used by dc:publisher
-- artist, which is used by atom:author

CREATE TABLE IF NOT EXISTS archives (
  id BIGSERIAL NOT NULL PRIMARY KEY, -- atom:id
  title TEXT NOT NULL, -- atom:title, do not use dc:title

  -- OPDS: https://specs.opds.io/opds-1.2#5-opds-catalog-entry-documents
  -- atom: https://www.ietf.org/rfc/rfc4287.txt, https://github.com/gorilla/feeds/blob/main/feed.go
  -- dc: https://www.dublincore.org/specifications/dublin-core/dcmi-terms/

  -- OPDS Catalog Entry Documents are Atom Entry documents.

  date_added TIMESTAMPTZ NOT NULL, -- atom:published, Format: RFC3339
  date_updated TIMESTAMPTZ NOT NULL, -- atom:updated, Format: RFC3339,
  series TEXT NOT NULL, -- atom:rights
  date_issued DATE NOT NULL, -- dc:issued, Format: ISO 8601-1
  -- is_new BOOLEAN NOT NULL, -- atom:category Commented temporarily

  file_path TEXT UNIQUE NOT NULL,
  thumbnail_path TEXT
);

CREATE TABLE IF NOT EXISTS archives_tags ( -- many-to-many relationship
  archive_id BIGINT NOT NULL,
  tag_id BIGINT NOT NULL,
  PRIMARY KEY (archive_id, tag_id),
  FOREIGN KEY (archive_id) REFERENCES archives(id),
  FOREIGN KEY (tag_id) REFERENCES tags(id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_archives_title ON archives(title);
CREATE INDEX IF NOT EXISTS idx_archives_series ON archives(series);
CREATE INDEX IF NOT EXISTS idx_archives_issued ON archives (date_issued);
CREATE INDEX IF NOT EXISTS idx_tags_namespace ON tags (namespace);
CREATE INDEX IF NOT EXISTS idx_tags_value ON tags (value);

-- +goose down
DROP INDEX IF EXISTS idx_archives_title;
DROP INDEX IF EXISTS idx_archives_series;
DROP INDEX IF EXISTS idx_archives_issued;
DROP INDEX IF EXISTS idx_tags_namespace;
DROP INDEX IF EXISTS idx_tags_value;
DROP TABLE IF EXISTS archives_tags;
DROP TABLE IF EXISTS archives;
DROP TABLE IF EXISTS tags;
