CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL, -- hashed password
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS category (
    category_id SERIAL PRIMARY KEY,
    name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS tags (
    tag_id SERIAL PRIMARY KEY,
    name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS events (
    event_uuid UUID PRIMARY KEY,
    creation_time TIMESTAMP DEFAULT NOW(),
    starting_time TIMESTAMP NOT NULL,
    ending_time TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    category_id INTEGER REFERENCES category(category_id),
    status TEXT NOT NULL,
    capacity INTEGER, -- how much they expecting to fit, NULL if its unlimimted
    image_url TEXT,
    organizer_id INTEGER REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS tickets (
    ticket_uuid UUID PRIMARY KEY,
    event_uuid UUID REFERENCES events(event_uuid) ON DELETE CASCADE,
    name TEXT NOT NULL,
    price INTEGER NOT NULL,
    currency CHAR(3) DEFAULT 'KZT' CHECK (currency ~ '^[A-Z]{3}$'),
    quantity INTEGER CHECK (quantity >= 0),
    sold INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS event_tags (
    event_uuid UUID REFERENCES events(event_uuid) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(tag_id) ON DELETE CASCADE,
    PRIMARY KEY (event_uuid, tag_id)
);

CREATE TABLE IF NOT EXISTS locations (
    location_id SERIAL PRIMARY KEY,
    event_uuid UUID REFERENCES events(event_uuid) ON DELETE CASCADE,
    name TEXT,
    address TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL
);

CREATE TABLE IF NOT EXISTS user_tickets (
    user_ticket_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    ticket_uuid UUID REFERENCES tickets(ticket_uuid) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1,
    purchase_time TIMESTAMP DEFAULT NOW(),
    status TEXT DEFAULT 'active' -- could be active, cancelled, refunded
);

-- Foreign key indexes
CREATE INDEX idx_events_category_id ON events(category_id);
CREATE INDEX idx_events_organizer_id ON events(organizer_id);
CREATE INDEX idx_tickets_event_uuid ON tickets(event_uuid);
CREATE INDEX idx_event_tags_event_uuid ON event_tags(event_uuid);
CREATE INDEX idx_event_tags_tag_id ON event_tags(tag_id);
CREATE INDEX idx_locations_event_uuid ON locations(event_uuid);


CREATE INDEX idx_event_ending_time ON events(ending_time);
CREATE INDEX idx_events_status ON events(status);   -- if you filter by status often

-- Basic categories
INSERT INTO category (name) VALUES (
    'Music', 'Tech', 'Sport', 'Other'
);