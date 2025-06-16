CREATE TABLE positions (
  code         VARCHAR(3) PRIMARY KEY NOT NULL,
  default_name TEXT NOT NULL
);

INSERT INTO positions (code, default_name) VALUES
  ('GLK', 'Goalkeeper'),
  ('DEF', 'Defender'),
  ('MID', 'Midfielder'),
  ('ATK', 'Attacker');

CREATE TABLE position_translations (
  position_code VARCHAR(3) NOT NULL REFERENCES positions(code) ON DELETE CASCADE,
  locale        VARCHAR(2) NOT NULL REFERENCES locales(code),
  label         TEXT NOT NULL,
  PRIMARY KEY (position_code, locale)
);

CREATE TABLE users (
  id          BIGINT PRIMARY KEY NOT NULL,
  username    VARCHAR(28) NOT NULL CHECK(LENGTH(username) >= 4),
  role        VARCHAR NOT NULL CHECK(role IN ('ADMIN', 'USER')),
  hash        VARCHAR(128) NOT NULL CHECK(LENGTH(hash) = 128),
  UNIQUE(username)
);

CREATE TABLE teams (
  id             BIGINT PRIMARY KEY NOT NULL,
  user_id        BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  country_code   VARCHAR(2) REFERENCES countries(code),
  budget         NUMERIC(15,2) NOT NULL CHECK (budget >= 0),
  total_players  INT NOT NULL DEFAULT 0,
  UNIQUE(user_id)
);

CREATE TABLE players (
  id            BIGINT PRIMARY KEY NOT NULL,
  team_id       BIGINT REFERENCES teams(id) ON DELETE SET NULL,
  country_code  VARCHAR(2) REFERENCES countries(code),
  age           INT NOT NULL CHECK (age >= 18 AND age <= 40),
  position_code VARCHAR(3) NOT NULL REFERENCES positions(code),
  price         NUMERIC(12,2) NOT NULL CHECK (price >= 0)
);

CREATE TABLE player_translations (
  id         BIGINT PRIMARY KEY NOT NULL,
  player_id  BIGINT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  locale     VARCHAR(2) NOT NULL REFERENCES locales(code),
  first_name VARCHAR(30) NOT NULL,
  last_name  VARCHAR(30) NOT NULL,
  UNIQUE (player_id, locale)
);
CREATE INDEX idx_player_translations_player_id ON player_translations(player_id);

CREATE TABLE transfer_lists (
  id             BIGINT PRIMARY KEY NOT NULL,
  player_id      BIGINT NOT NULL UNIQUE REFERENCES players(id) ON DELETE CASCADE,
  seller_team_id BIGINT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
  price          NUMERIC(12,2) NOT NULL CHECK (price >= 0),
  listed_at      TIMESTAMPTZ    NOT NULL DEFAULT now(),
  UNIQUE(player_id)
);

CREATE TABLE transfer_records (
  id             BIGINT PRIMARY KEY NOT NULL,
  player_id      BIGINT    NOT NULL REFERENCES players(id),
  seller_team_id BIGINT    NOT NULL REFERENCES teams(id),
  buyer_team_id  BIGINT    NOT NULL REFERENCES teams(id),
  sold_price     NUMERIC(12,2) NOT NULL CHECK (sold_price >= 0),
  listed_at      TIMESTAMPTZ NOT NULL,
  sold_at        TIMESTAMPTZ NOT NULL
);

CREATE OR REPLACE FUNCTION trg_players_total_players() RETURNS TRIGGER
LANGUAGE plpgsql AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    UPDATE teams
      SET total_players = total_players + 1
      WHERE id = NEW.team_id;
  
  ELSIF TG_OP = 'DELETE' THEN
    UPDATE teams
      SET total_players = total_players - 1
      WHERE id = OLD.team_id;

  ELSIF TG_OP = 'UPDATE' AND NEW.team_id <> OLD.team_id THEN
    UPDATE teams
      SET total_players = total_players - 1
      WHERE id = OLD.team_id;
    UPDATE teams
      SET total_players = total_players + 1
      WHERE id = NEW.team_id;
  END IF;

  RETURN NULL;
END;
$$;

CREATE TRIGGER players_total_players_trg
AFTER INSERT OR DELETE OR UPDATE OF team_id
ON players
FOR EACH ROW
EXECUTE FUNCTION trg_players_total_players();
