-- name: EditQuote :one
SELECT *
FROM quotes
WHERE id=?1;