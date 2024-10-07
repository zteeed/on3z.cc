\c docker;
CREATE TABLE IF NOT EXISTS short_url_maps (long_url TEXT NOT NULL, description TEXT, short_url TEXT UNIQUE NOT NULL, auth0_sub TEXT);
-- ALTER TABLE short_url_maps ADD COLUMN description TEXT;
