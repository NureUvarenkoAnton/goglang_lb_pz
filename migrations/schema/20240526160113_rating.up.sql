CREATE TABLE IF NOT EXISTS ratings (
  rater_id BIGINT REFERENCES users(id),
  ratee_id BIGINT REFERENCES users(id),
  value INTEGER,
  PRIMARY KEY (rater_id, ratee_id)
);
