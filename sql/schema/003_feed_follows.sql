-- +goose Up
CREATE TABLE feed_follows (
  user_id uuid NOT NULL,
  feed_id uuid NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  UNIQUE(user_id, feed_id),
  PRIMARY KEY (user_id, feed_id),
  CONSTRAINT feed_follows_user_id
  FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
  CONSTRAINT feed_follows_feed_id
  FOREIGN KEY (feed_id)
    REFERENCES feeds(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feed_follows;