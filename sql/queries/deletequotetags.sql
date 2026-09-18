-- name: DeleteQuoteTags :exec
DELETE FROM quotes_tags_link
WHERE quote_id=?1;