BEGIN TRANSACTION;

CREATE TABLE urls(
	id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	short_url VARCHAR(200) NOT NULL,
	original_url VARCHAR(300) NOT NULL
);

CREATE UNIQUE INDEX short_url_index ON urls(short_url);
CREATE UNIQUE INDEX original_url_index ON urls(original_url);

COMMIT;
