CREATE TABLE IF NOT EXISTS pull_request_statuses (
    id SERIAL PRIMARY KEY,
    status VARCHAR(50)
);

INSERT INTO pull_request_statuses (status) VALUES
('open'),
('merged');

CREATE TABLE IF NOT EXISTS pull_requests (
    pull_request_id VARCHAR(50) NOT NULL PRIMARY KEY,
    pull_request_name VARCHAR(255) NOT NULL,
    status_id INT DEFAULT 1 REFERENCES pull_request_statuses(id) ON DELETE SET DEFAULT,
    is_need_more_reviewers BOOLEAN NOT NULL DEFAULT TRUE,
    author_id VARCHAR(50) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    merged_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pull_request_reviewers (
    pull_request_id VARCHAR(50) NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    reviewer_id VARCHAR(50) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT pk_pull_request_reviews PRIMARY KEY (pull_request_id, reviewer_id)
);

CREATE INDEX idx_pr_author ON pull_requests(author_id);

-- TODO: продумать индексы
