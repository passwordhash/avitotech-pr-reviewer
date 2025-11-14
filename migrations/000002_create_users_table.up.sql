CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(50) NOT NULL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    team_id UUID REFERENCES teams(team_id)
);
