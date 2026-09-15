-- name: PutEditedQuote :one
UPDATE quotes
SET
    updated_at = ?4,
    quote = ?2,
    author = ?3
WHERE id = ?1
RETURNING *;