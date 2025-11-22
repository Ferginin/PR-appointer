-- Backend команда (10 человек)
INSERT INTO users (username, is_active) VALUES
    ('alice', true),
    ('bob', true),
    ('charlie', true),
    ('david', true),
    ('emma', true),
    ('frank', true),
    ('grace', true),
    ('henry', true),
    ('ivan', true),
    ('julia', true)
ON CONFLICT (username) DO NOTHING;

-- Frontend команда (10 человек)
INSERT INTO users (username, is_active) VALUES
    ('isabella', true),
    ('jack', true),
    ('kate', true),
    ('liam', true),
    ('mia', true),
    ('noah', true),
    ('olivia', true),
    ('peter', true),
    ('quinn', true),
    ('rachel', true)
ON CONFLICT (username) DO NOTHING;

-- DevOps команда (10 человек)
INSERT INTO users (username, is_active) VALUES
    ('sam', true),
    ('tina', true),
    ('uma', true),
    ('victor', true),
    ('wendy', true),
    ('xavier', true),
    ('yara', true),
    ('zack', true),
    ('anna', true),
    ('brian', true)
ON CONFLICT (username) DO NOTHING;

-- Mobile команда (10 человек)
INSERT INTO users (username, is_active) VALUES
    ('clara', true),
    ('daniel', true),
    ('elena', true),
    ('finn', true),
    ('georgia', true),
    ('hugo', true),
    ('iris', true),
    ('james', true),
    ('kelly', true),
    ('leo', true)
ON CONFLICT (username) DO NOTHING;

-- QA команда (10 человек)
INSERT INTO users (username, is_active) VALUES
    ('maya', true),
    ('nathan', true),
    ('oscar', true),
    ('paula', true),
    ('quincy', true),
    ('rose', true),
    ('steve', true),
    ('tara', true),
    ('ulysses', false),  -- неактивный
    ('vera', false)     -- неактивный
ON CONFLICT (username) DO NOTHING;

-- Создание тестовых команд
INSERT INTO teams (name) VALUES
    ('Backend Team'),
    ('Frontend Team'),
    ('DevOps Team'),
    ('Mobile Team'),
    ('QA Team')
ON CONFLICT (name) DO NOTHING;

-- Добавление пользователей в команды
-- Backend Team - пользователи 1-10
INSERT INTO team_members (team_id, user_id)
SELECT
    (SELECT id FROM teams WHERE name = 'Backend Team'),
    id
FROM users
WHERE username IN ('alice', 'bob', 'charlie', 'david', 'emma', 'frank', 'grace', 'henry', 'ivan', 'julia')
ON CONFLICT (team_id, user_id) DO NOTHING;

-- Frontend Team - пользователи 11-20
INSERT INTO team_members (team_id, user_id)
SELECT
    (SELECT id FROM teams WHERE name = 'Frontend Team'),
    id
FROM users
WHERE username IN ('isabella', 'jack', 'kate', 'liam', 'mia', 'noah', 'olivia', 'peter', 'quinn', 'rachel')
ON CONFLICT (team_id, user_id) DO NOTHING;

-- DevOps Team - пользователи 21-30
INSERT INTO team_members (team_id, user_id)
SELECT
    (SELECT id FROM teams WHERE name = 'DevOps Team'),
    id
FROM users
WHERE username IN ('sam', 'tina', 'uma', 'victor', 'wendy', 'xavier', 'yara', 'zack', 'anna', 'brian')
ON CONFLICT (team_id, user_id) DO NOTHING;

-- Mobile Team - пользователи 31-40
INSERT INTO team_members (team_id, user_id)
SELECT
    (SELECT id FROM teams WHERE name = 'Mobile Team'),
    id
FROM users
WHERE username IN ('clara', 'daniel', 'elena', 'finn', 'georgia', 'hugo', 'iris', 'james', 'kelly', 'leo')
ON CONFLICT (team_id, user_id) DO NOTHING;

-- QA Team - пользователи 41-50
INSERT INTO team_members (team_id, user_id)
SELECT
    (SELECT id FROM teams WHERE name = 'QA Team'),
    id
FROM users
WHERE username IN ('maya', 'nathan', 'oscar', 'paula', 'quincy', 'rose', 'steve', 'tara', 'ulysses', 'vera')
ON CONFLICT (team_id, user_id) DO NOTHING;

INSERT INTO pull_requests (id, title, author_id, status) VALUES
    (1001, 'Add user authentication API', (SELECT id FROM users WHERE username = 'alice'), 'MERGED'),
    (1002, 'Implement JWT token refresh', (SELECT id FROM users WHERE username = 'bob'), 'OPEN'),
    (1003, 'Add database migration for teams', (SELECT id FROM users WHERE username = 'charlie'), 'OPEN'),
    (1004, 'Fix memory leak in reviewer assignment', (SELECT id FROM users WHERE username = 'david'), 'MERGED'),
    (1005, 'Optimize SQL queries for large datasets', (SELECT id FROM users WHERE username = 'emma'), 'OPEN')
ON CONFLICT (id) DO NOTHING;

-- Frontend PRs (5 штук)
INSERT INTO pull_requests (id, title, author_id, status) VALUES
    (2001, 'Redesign user profile page', (SELECT id FROM users WHERE username = 'isabella'), 'MERGED'),
    (2002, 'Add dark mode support', (SELECT id FROM users WHERE username = 'jack'), 'OPEN'),
    (2003, 'Fix responsive layout on mobile', (SELECT id FROM users WHERE username = 'kate'), 'OPEN'),
    (2004, 'Implement infinite scroll for PR list', (SELECT id FROM users WHERE username = 'liam'), 'MERGED'),
    (2005, 'Add loading states to async operations', (SELECT id FROM users WHERE username = 'mia'), 'OPEN')
ON CONFLICT (id) DO NOTHING;

-- DevOps PRs (5 штук)
INSERT INTO pull_requests (id, title, author_id, status) VALUES
    (3001, 'Setup CI/CD pipeline with GitHub Actions', (SELECT id FROM users WHERE username = 'sam'), 'MERGED'),
    (3002, 'Add Docker Compose for local development', (SELECT id FROM users WHERE username = 'tina'), 'OPEN'),
    (3003, 'Configure monitoring with Prometheus', (SELECT id FROM users WHERE username = 'uma'), 'OPEN'),
    (3004, 'Add automated backups for PostgreSQL', (SELECT id FROM users WHERE username = 'victor'), 'MERGED'),
    (3005, 'Implement blue-green deployment strategy', (SELECT id FROM users WHERE username = 'wendy'), 'OPEN')
ON CONFLICT (id) DO NOTHING;

-- Mobile PRs (5 штук)
INSERT INTO pull_requests (id, title, author_id, status) VALUES
    (4001, 'Add push notifications support', (SELECT id FROM users WHERE username = 'clara'), 'MERGED'),
    (4002, 'Implement offline mode with local cache', (SELECT id FROM users WHERE username = 'daniel'), 'OPEN'),
    (4003, 'Fix crash on app launch for Android 14', (SELECT id FROM users WHERE username = 'elena'), 'OPEN'),
    (4004, 'Add biometric authentication', (SELECT id FROM users WHERE username = 'finn'), 'MERGED'),
    (4005, 'Optimize battery usage in background', (SELECT id FROM users WHERE username = 'georgia'), 'OPEN')
ON CONFLICT (id) DO NOTHING;

-- QA PRs (5 штук)
INSERT INTO pull_requests (id, title, author_id, status) VALUES
    (5001, 'Add end-to-end tests with Playwright', (SELECT id FROM users WHERE username = 'maya'), 'MERGED'),
    (5002, 'Implement load testing with k6', (SELECT id FROM users WHERE username = 'nathan'), 'OPEN'),
    (5003, 'Add smoke tests for critical flows', (SELECT id FROM users WHERE username = 'oscar'), 'OPEN'),
    (5004, 'Setup test data factory for integration tests', (SELECT id FROM users WHERE username = 'paula'), 'MERGED'),
    (5005, 'Add visual regression testing', (SELECT id FROM users WHERE username = 'quincy'), 'OPEN')
ON CONFLICT (id) DO NOTHING;

-- Backend PRs reviewers (выбираем из той же команды, исключая автора)
-- PR 1001 (alice) -> bob, charlie
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (1001, (SELECT id FROM users WHERE username = 'bob')),
    (1001, (SELECT id FROM users WHERE username = 'charlie'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 1002 (bob) -> alice, david
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (1002, (SELECT id FROM users WHERE username = 'alice')),
    (1002, (SELECT id FROM users WHERE username = 'david'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 1003 (charlie) -> emma, frank
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (1003, (SELECT id FROM users WHERE username = 'emma')),
    (1003, (SELECT id FROM users WHERE username = 'frank'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 1004 (david) -> grace, henry
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (1004, (SELECT id FROM users WHERE username = 'grace')),
    (1004, (SELECT id FROM users WHERE username = 'henry'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 1005 (emma) -> ivan, julia
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (1005, (SELECT id FROM users WHERE username = 'ivan')),
    (1005, (SELECT id FROM users WHERE username = 'julia'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- Frontend PRs reviewers
-- PR 2001 (isabella) -> jack, kate
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (2001, (SELECT id FROM users WHERE username = 'jack')),
    (2001, (SELECT id FROM users WHERE username = 'kate'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 2002 (jack) -> isabella, liam
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (2002, (SELECT id FROM users WHERE username = 'isabella')),
    (2002, (SELECT id FROM users WHERE username = 'liam'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 2003 (kate) -> mia, noah
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (2003, (SELECT id FROM users WHERE username = 'mia')),
    (2003, (SELECT id FROM users WHERE username = 'noah'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 2004 (liam) -> olivia, peter
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (2004, (SELECT id FROM users WHERE username = 'olivia')),
    (2004, (SELECT id FROM users WHERE username = 'peter'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 2005 (mia) -> quinn, rachel
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (2005, (SELECT id FROM users WHERE username = 'quinn')),
    (2005, (SELECT id FROM users WHERE username = 'rachel'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- DevOps PRs reviewers
-- PR 3001 (sam) -> tina, uma
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (3001, (SELECT id FROM users WHERE username = 'tina')),
    (3001, (SELECT id FROM users WHERE username = 'uma'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 3002 (tina) -> sam, victor
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (3002, (SELECT id FROM users WHERE username = 'sam')),
    (3002, (SELECT id FROM users WHERE username = 'victor'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 3003 (uma) -> wendy, xavier
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (3003, (SELECT id FROM users WHERE username = 'wendy')),
    (3003, (SELECT id FROM users WHERE username = 'xavier'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 3004 (victor) -> yara, zack
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (3004, (SELECT id FROM users WHERE username = 'yara')),
    (3004, (SELECT id FROM users WHERE username = 'zack'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 3005 (wendy) -> anna, brian
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (3005, (SELECT id FROM users WHERE username = 'anna')),
    (3005, (SELECT id FROM users WHERE username = 'brian'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- Mobile PRs reviewers
-- PR 4001 (clara) -> daniel, elena
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (4001, (SELECT id FROM users WHERE username = 'daniel')),
    (4001, (SELECT id FROM users WHERE username = 'elena'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 4002 (daniel) -> clara, finn
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (4002, (SELECT id FROM users WHERE username = 'clara')),
    (4002, (SELECT id FROM users WHERE username = 'finn'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 4003 (elena) -> georgia, hugo
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (4003, (SELECT id FROM users WHERE username = 'georgia')),
    (4003, (SELECT id FROM users WHERE username = 'hugo'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 4004 (finn) -> iris, james
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (4004, (SELECT id FROM users WHERE username = 'iris')),
    (4004, (SELECT id FROM users WHERE username = 'james'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 4005 (georgia) -> kelly, leo
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (4005, (SELECT id FROM users WHERE username = 'kelly')),
    (4005, (SELECT id FROM users WHERE username = 'leo'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- QA PRs reviewers
-- PR 5001 (maya) -> nathan, oscar
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (5001, (SELECT id FROM users WHERE username = 'nathan')),
    (5001, (SELECT id FROM users WHERE username = 'oscar'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 5002 (nathan) -> maya, paula
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (5002, (SELECT id FROM users WHERE username = 'maya')),
    (5002, (SELECT id FROM users WHERE username = 'paula'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 5003 (oscar) -> quincy, rose
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (5003, (SELECT id FROM users WHERE username = 'quincy')),
    (5003, (SELECT id FROM users WHERE username = 'rose'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 5004 (paula) -> steve, tara
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (5004, (SELECT id FROM users WHERE username = 'steve')),
    (5004, (SELECT id FROM users WHERE username = 'tara'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;

-- PR 5005 (quincy) -> maya, nathan
INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES
    (5005, (SELECT id FROM users WHERE username = 'maya')),
    (5005, (SELECT id FROM users WHERE username = 'nathan'))
ON CONFLICT (pr_id, reviewer_id) DO NOTHING;