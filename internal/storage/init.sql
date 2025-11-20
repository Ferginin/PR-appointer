-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Таблица команд
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Таблица связи пользователей и команд (многие ко многим)
CREATE TABLE IF NOT EXISTS team_members (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(team_id, user_id)
    );

-- Таблица Pull Request'ов
CREATE TABLE IF NOT EXISTS pull_requests (
    id SERIAL PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status VARCHAR(20) NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Таблица назначенных ревьюверов на PR
CREATE TABLE IF NOT EXISTS pr_reviewers (
    id SERIAL PRIMARY KEY,
    pr_id INTEGER NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(pr_id, reviewer_id)
    );



-- Индекс для быстрого поиска PR по автору
CREATE INDEX IF NOT EXISTS idx_pr_author ON pull_requests(author_id);

-- Индекс для быстрого поиска PR по статусу
CREATE INDEX IF NOT EXISTS idx_pr_status ON pull_requests(status);

-- Индекс для быстрого поиска ревьюверов по PR
CREATE INDEX IF NOT EXISTS idx_pr_reviewers_pr ON pr_reviewers(pr_id);

-- Индекс для быстрого поиска PR по ревьюверу
CREATE INDEX IF NOT EXISTS idx_pr_reviewers_reviewer ON pr_reviewers(reviewer_id);

-- Индекс для быстрого поиска участников команды
CREATE INDEX IF NOT EXISTS idx_team_members_team ON team_members(team_id);
CREATE INDEX IF NOT EXISTS idx_team_members_user ON team_members(user_id);

-- Индекс для быстрого поиска активных пользователей
CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active) WHERE is_active = TRUE;



-- Создание тестовых пользователей
INSERT INTO users (username, is_active) VALUES
    ('alice', TRUE),
    ('bob', TRUE),
    ('charlie', TRUE),
    ('david', TRUE),
    ('eve', FALSE),
    ('danil', TRUE)
    ON CONFLICT (username) DO NOTHING;

-- Создание тестовых команд
INSERT INTO teams (name) VALUES
    ('Backend Team'),
    ('Frontend Team')
    ON CONFLICT (name) DO NOTHING;

-- Добавление пользователей в команды
INSERT INTO team_members (team_id, user_id)
SELECT t.id, u.id
FROM teams t, users u
WHERE t.name = 'Backend Team' AND u.username IN ('alice', 'bob', 'charlie', 'eve')
    ON CONFLICT (team_id, user_id) DO NOTHING;

INSERT INTO team_members (team_id, user_id)
SELECT t.id, u.id
FROM teams t, users u
WHERE t.name = 'Frontend Team' AND u.username IN ('david', 'danil')
    ON CONFLICT (team_id, user_id) DO NOTHING;