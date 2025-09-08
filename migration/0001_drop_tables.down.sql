DROP INDEX IF EXISTS idx_events_status;
DROP INDEX IF EXISTS idx_locations_event_uuid;
DROP INDEX IF EXISTS idx_event_tags_tag_id;
DROP INDEX IF EXISTS idx_event_tags_event_uuid;
DROP INDEX IF EXISTS idx_tickets_event_uuid;
DROP INDEX IF EXISTS idx_events_organizer_id;
DROP INDEX IF EXISTS idx_events_category_id;

DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS event_tags;
DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS users;