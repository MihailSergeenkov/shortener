BEGIN TRANSACTION;

CREATE TABLE urls(
	id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	short_url VARCHAR(200) NOT NULL,
	original_url VARCHAR(300) NOT NULL
);

ALTER TABLE urls ADD COLUMN user_id VARCHAR(200);
CREATE INDEX urls_user_id_index ON urls(user_id);
ALTER TABLE urls ADD COLUMN is_deleted boolean DEFAULT false;
CREATE UNIQUE INDEX short_url_index ON urls(short_url) WHERE is_deleted = false;
CREATE UNIQUE INDEX original_url_index ON urls(original_url) WHERE is_deleted = false;

COMMIT;
