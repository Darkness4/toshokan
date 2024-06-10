PRAGMA foreign_keys = ON;
BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS tags ( -- atom:summary
  id INTEGER PRIMARY KEY,
  namespace TEXT NOT NULL, -- used for filtering
  value TEXT -- used for sorting
);

CREATE TABLE IF NOT EXISTS archives (
  id BLOB NOT NULL PRIMARY KEY, -- atom:id
  title TEXT NOT NULL, -- atom:title, do not use dc:title

  -- OPDS: https://specs.opds.io/opds-1.2#5-opds-catalog-entry-documents
  -- atom: https://www.ietf.org/rfc/rfc4287.txt, https://github.com/gorilla/feeds/blob/main/feed.go
  -- dc: https://www.dublincore.org/specifications/dublin-core/dcmi-terms/

  -- OPDS Catalog Entry Documents are Atom Entry documents.

  date_added INTEGER NOT NULL, -- atom:published, Format: RFC3339
  date_updated INTEGER NOT NULL, -- atom:updated, Format: RFC3339,
  author TEXT NOT NULL, -- atom:author, do not use dc:creator
  series TEXT NOT NULL, -- atom:rights
  date_issued TEXT NOT NULL, -- dc:issued, Format: ISO 8601-1
  language TEXT NOT NULL, -- dc:language
  publisher TEXT NOT NULL, -- dc:publisher (circle)
  issued TEXT NOT NULL, -- dc:issued (event)
  -- is_new BOOLEAN NOT NULL, -- atom:category Commented temporarily

  file_path TEXT NOT NULL,
  thumbnail_path TEXT
);

CREATE TABLE archives_tags ( -- many-to-many relationship
  archive_id INTEGER NOT NULL,
  tag_id BLOB NOT NULL,
  PRIMARY KEY (archive_id, tag_id),
  FOREIGN KEY (archive_id) REFERENCES archives(id),
  FOREIGN KEY (tag_id) REFERENCES tags(id)
);

-- Remove tags that are not used by any archive
CREATE TRIGGER delete_unused_tags
AFTER DELETE ON archives_tags
BEGIN
    DELETE FROM tags
    WHERE id NOT IN (SELECT tag_id FROM archives_tags);
END;
COMMIT;
