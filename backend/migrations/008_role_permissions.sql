-- 动态角色权限配置表
CREATE TABLE IF NOT EXISTS role_permissions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    role VARCHAR(30) NOT NULL,
    permission VARCHAR(50) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_role_permission (role, permission)
);

-- 初始化默认权限数据（对齐当前硬编码逻辑）
-- super_admin 拥有全部权限，不需要存储（代码中直接放行）

-- admin 角色
INSERT INTO role_permissions (role, permission) VALUES
('admin', 'dashboard'),
('admin', 'courses'),
('admin', 'categories'),
('admin', 'tags'),
('admin', 'users'),
('admin', 'orders'),
('admin', 'guestbook'),
('admin', 'comments'),
('admin', 'home_config'),
('admin', 'data'),
('admin', 'tickets'),
('admin', 'settings'),
('admin', 'admins');

-- operator 角色（运营）
INSERT INTO role_permissions (role, permission) VALUES
('operator', 'dashboard'),
('operator', 'courses'),
('operator', 'categories'),
('operator', 'tags'),
('operator', 'orders'),
('operator', 'home_config'),
('operator', 'data');

-- support 角色（客服）
INSERT INTO role_permissions (role, permission) VALUES
('support', 'tickets'),
('support', 'guestbook'),
('support', 'comments');
