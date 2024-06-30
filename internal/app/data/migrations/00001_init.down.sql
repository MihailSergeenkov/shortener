BEGIN TRANSACTION;

ALTER TABLE urls DROP COLUMN is_deleted;
DROP INDEX urls_user_id_index;
ALTER TABLE urls DROP COLUMN user_id;
DROP INDEX original_url_index;
DROP INDEX short_url_index;
DROP TABLE urls;

COMMIT;
