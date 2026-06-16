CREATE TABLE IF NOT EXISTS books (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       VARCHAR(255)  NOT NULL,
    author      VARCHAR(255)  NOT NULL,
    genre       VARCHAR(100)  NOT NULL,
    description TEXT          DEFAULT '',
    cover_url   VARCHAR(500)  DEFAULT '',
    file_url    VARCHAR(500)  DEFAULT '',
    pages       INT           DEFAULT 0,
    year        INT           DEFAULT 0,
    uploaded_by UUID          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMP     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP     NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_books_title  ON books(title);
CREATE INDEX idx_books_author ON books(author);
CREATE INDEX idx_books_genre  ON books(genre);