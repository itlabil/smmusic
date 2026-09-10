CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    album VARCHAR(255),
    genre VARCHAR(100),
    duration_seconds INT,
    source_format VARCHAR(10) NOT NULL, -- 'flac' | 'mp3'
    flac_path VARCHAR(500),
    mp3_path VARCHAR(500),
    cover_path VARCHAR(500),
    transcode_status VARCHAR(20) NOT NULL DEFAULT 'done', -- 'pending' | 'processing' | 'done' | 'failed'
    uploaded_by INT REFERENCES users(id) ON DELETE SET NULL,
    search_vector TSVECTOR,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_songs_search ON songs USING GIN(search_vector);
CREATE INDEX idx_songs_artist ON songs(artist);

-- Auto-update search_vector saat insert/update
CREATE FUNCTION songs_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('simple', coalesce(NEW.title, '')), 'A') ||
        setweight(to_tsvector('simple', coalesce(NEW.artist, '')), 'B') ||
        setweight(to_tsvector('simple', coalesce(NEW.album, '')), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER songs_search_vector_trigger
    BEFORE INSERT OR UPDATE ON songs
    FOR EACH ROW EXECUTE FUNCTION songs_search_vector_update();