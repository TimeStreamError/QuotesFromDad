-- +goose up
CREATE TABLE quotes_tags_link (
    quote_id TEXT NOT NULL,
    tag_id TEXT NOT NULL,
    UNIQUE (quote_id, tag_id),
    FOREIGN KEY (quote_id) REFERENCES quotes(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- +goose down
DROP TABLE quotes_tags_link;