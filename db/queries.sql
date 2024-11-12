-- name: DeleteUnusedTags :exec
DELETE FROM tags WHERE id NOT IN (SELECT tag_id FROM archives_tags);

-- name: CreateTag :exec
INSERT INTO tags (namespace, value)
VALUES (@namespace, @value)
ON CONFLICT
DO NOTHING;

-- name: CreateArchive :exec
INSERT INTO archives (
  title,
  date_added,
  date_updated,
  series,
  date_issued,
  file_path,
  thumbnail_path
) values (
  @title,
  @date_added,
  @date_updated,
  @series,
  @date_issued,
  @file_path,
  @thumbnail_path
);
