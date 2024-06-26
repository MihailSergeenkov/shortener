BEGIN TRANSACTION;

DROP INDEX original_url_index;
DROP INDEX short_url_index;
CREATE UNIQUE INDEX short_url_index ON urls(short_url) WHERE is_deleted = false;
CREATE UNIQUE INDEX original_url_index ON urls(original_url) WHERE is_deleted = false;

COMMIT;
