-- PPPoE Users Table
-- Stores all customer PPPoE accounts for each router

CREATE TABLE IF NOT EXISTS pppoe_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    router_id INT NOT NULL,
    profile VARCHAR(100) DEFAULT 'default',
    is_enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY unique_user_router (username, router_id),
    INDEX idx_router (router_id),
    INDEX idx_enabled (is_enabled),
    INDEX idx_username (username)
);

-- Optional: Add foreign key to routers table
-- ALTER TABLE pppoe_users ADD CONSTRAINT fk_router
--   FOREIGN KEY (router_id) REFERENCES routers(id) ON DELETE CASCADE;
