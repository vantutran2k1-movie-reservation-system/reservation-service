DROP TABLE IF EXISTS seats;

CREATE TYPE seat_type AS ENUM('REGULAR', 'VIP');

CREATE TABLE seats (
    id UUID PRIMARY KEY,
    theater_id UUID NOT NULL REFERENCES theaters(id) ON DELETE CASCADE,
    row VARCHAR(1) NOT NULL,
    number INT NOT NULL,
    type seat_type NOT NULL DEFAULT 'REGULAR'
);

ALTER TABLE seats
ADD CONSTRAINT unique_seat_in_theater UNIQUE (theater_id, row, number);