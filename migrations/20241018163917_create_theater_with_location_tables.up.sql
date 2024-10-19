CREATE TABLE theaters (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE INDEX idx_theater_name ON theaters(name);

CREATE TABLE theater_locations (
    id UUID PRIMARY KEY,
    theater_id UUID NOT NULL REFERENCES theaters(id) ON DELETE CASCADE,
    city_id UUID NOT NULL REFERENCES cities(id) ON DELETE CASCADE,
    address VARCHAR(255) NOT NULL,
    postal_code VARCHAR(20),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8)
);

CREATE INDEX idx_theater_location_city ON theater_locations(city_id);
CREATE INDEX idx_theater_location_lat_lng ON theater_locations(latitude, longitude);