-- +goose Up
CREATE TABLE posts(
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	udpated_at TIMESTAMP NOT NULL,
	title TEXT NOT NULL,
	url TEXT NOT NULL,
	description TEXT NOT NULL,
	published_at TIMESTAMP,
	feed_id UUID NOT NULL REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;
