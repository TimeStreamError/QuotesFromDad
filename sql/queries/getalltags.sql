-- name: GetAllTags :many
SELECT name, id
FROM tags
ORDER BY name ASC;
