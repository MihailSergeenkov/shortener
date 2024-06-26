BEGIN TRANSACTION;

DROP INDEX original_url_index;
DROP INDEX short_url_index;
DROP TABLE urls;

COMMIT;
