-- name: DeleteUnusedTags :exec
DELETE FROM tags    WHERE id NOT IN (SELECT tag_id FROM archives_tags);
