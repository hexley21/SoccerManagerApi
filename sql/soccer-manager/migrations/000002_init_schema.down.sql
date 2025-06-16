DROP TRIGGER IF EXISTS players_total_players_trg ON players;

DROP FUNCTION IF EXISTS trg_players_total_players();

DROP INDEX IF EXISTS idx_player_translations_player_id;

DROP TABLE IF EXISTS transfer_records;
DROP TABLE IF EXISTS transfer_lists;

DROP TABLE IF EXISTS player_translations;
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS position_translations;
DROP TABLE IF EXISTS positions;