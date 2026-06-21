-- 提现申请表
CREATE TABLE IF NOT EXISTS withdrawal_requests (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount INT NOT NULL,
    account_name VARCHAR(100),
    account_no VARCHAR(100),
    bank_name VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    remark TEXT,
    processed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_withdrawals_user (user_id),
    INDEX idx_withdrawals_status (status),
    INDEX idx_withdrawals_created_at (created_at),
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
