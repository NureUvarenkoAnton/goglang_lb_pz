CREATE TABLE IF NOT EXISTS walks (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  owner_id BIGINT REFERENCES users(id),
  walker_id BIGINT REFERENCES users(id),
  pet_id BIGINT REFERENCES pets(id),
  start_time DATETIME,
  finish_time DATETIME DEFAULT NULL,
  state ENUM('pending', 'accepted', 'declined', 'in_proccess', 'finished')
);
