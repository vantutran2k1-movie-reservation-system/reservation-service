DROP TABLE IF EXISTS shows;
DROP TYPE IF EXISTS show_status;

CREATE TYPE show_status AS ENUM ('ACTIVE', 'CANCELLED', 'COMPLETED', 'EXPIRED', 'SCHEDULED', 'ON-HOLD');
CREATE TABLE IF NOT EXISTS shows (
    id UUID PRIMARY KEY,
    movie_id UUID NOT NULL,
    theater_id UUID NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status show_status NOT NULL DEFAULT 'SCHEDULED',
    created_at TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    CONSTRAINT fk_movie FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE CASCADE,
    CONSTRAINT fk_theater FOREIGN KEY (theater_id) REFERENCES theaters (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_show_movie_id ON shows (movie_id);
CREATE INDEX IF NOT EXISTS idx_show_theater_id ON shows (theater_id);
CREATE INDEX IF NOT EXISTS idx_show_time_range ON shows (start_time, end_time);