\c docker;
CREATE TABLE IF NOT EXISTS short_url_maps (long_url TEXT UNIQUE NOT NULL, short_url TEXT UNIQUE NOT NULL);