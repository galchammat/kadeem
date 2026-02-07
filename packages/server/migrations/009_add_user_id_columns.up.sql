-- Create junction table for user → riot account tracking
CREATE TABLE user_tracked_accounts (
    user_id UUID NOT NULL,
    account_id INTEGER NOT NULL REFERENCES league_of_legends_accounts(id) ON DELETE CASCADE,
    tracked_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, account_id)
);

CREATE INDEX idx_user_tracked_accounts_user_id ON user_tracked_accounts(user_id);
CREATE INDEX idx_user_tracked_accounts_account_id ON user_tracked_accounts(account_id);

-- Create junction table for user → streamer tracking
CREATE TABLE user_tracked_streamers (
    user_id UUID NOT NULL,
    streamer_id INTEGER NOT NULL REFERENCES streamers(id) ON DELETE CASCADE,
    tracked_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, streamer_id)
);

CREATE INDEX idx_user_tracked_streamers_user_id ON user_tracked_streamers(user_id);
CREATE INDEX idx_user_tracked_streamers_streamer_id ON user_tracked_streamers(streamer_id);
