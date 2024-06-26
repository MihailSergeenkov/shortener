BEGIN TRANSACTION;

ALTER TABLE urls ADD COLUMN user_id VARCHAR(200);
CREATE INDEX urls_user_id_index ON urls(user_id);

COMMIT;
