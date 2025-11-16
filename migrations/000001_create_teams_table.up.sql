CREATE TABLE IF NOT EXISTS teams (
    team_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    team_name VARCHAR(100) UNIQUE NOT NULL
);