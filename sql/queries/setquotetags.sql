-- name: SetQuoteTags :exec
INSERT INTO quotes_tags_link (quote_id, tag_id)
VALUES (
    ?1,
    ?2
);

