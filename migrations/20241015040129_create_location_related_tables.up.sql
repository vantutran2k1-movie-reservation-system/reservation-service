CREATE TABLE countries (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    code CHAR(2) NOT NULL UNIQUE
);

CREATE INDEX idx_country_name ON countries(name);

CREATE TABLE states (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(10),
    country_id UUID NOT NULL REFERENCES countries(id) ON DELETE CASCADE
);

CREATE INDEX idx_state_name ON states(name);
CREATE INDEX idx_state_country ON states(country_id);

CREATE TABLE cities (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    state_id UUID NOT NULL REFERENCES states(id) ON DELETE CASCADE
);

CREATE INDEX idx_city_name ON cities(name);
CREATE INDEX idx_city_state ON cities(state_id);
