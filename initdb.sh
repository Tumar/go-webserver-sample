#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE TABLE books (
        isbn char(18) NOT NULL,
        title varchar(255) NOT NULL,
        author varchar(255) NOT NULL,
        price decimal(5,2) NOT NULL
    );

    INSERT INTO books (isbn, title, author, price) VALUES
    ('978-1-50326-196-9', 'Emma', 'Jayne Austen', 9.44),
    ('978-1-50525-560-7', 'The Time Machine', 'H. G. Wells', 5.99),
    ('978-1-50337-964-0', 'The Prince', 'Niccolò Machiavelli', 6.99);

    ALTER TABLE books ADD PRIMARY KEY (isbn);
EOSQL