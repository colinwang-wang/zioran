# 07 - 数据库设计

## 7.1 ER 图（核心表关系）

```
users ──1:N──→ orders ──→ courses
  │               │
  │               └──→ vip_packages
  │
  ├──1:1──→ coin_accounts
  ├──1:N──→ coin_transactions
  ├──1:1──→ user_vip
  ├──1:N──→ guestbook_messages
  ├──1:N──→ guestbook_likes
  ├──1:N──→ comments
  ├──1:N──→ tickets（工单）
  ├──1:N──→ favorites
  └──1:N──→ download_logs

courses ──N:M──→ tags (course_tags)
courses ──N:1──→ categories
```

---

## 7.2 完整建表 SQL

```sql
-- ============================================
-- 用户相关
-- ============================================

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(100),
    password_hash VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500),
    bio TEXT,
    role VARCHAR(20) NOT NULL DEFAULT 'user',  -- user, vip, admin
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, disabled
    wechat_openid VARCHAR(100),  -- 微信登录绑定
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_wechat ON users(wechat_openid) WHERE wechat_openid IS NOT NULL;

-- ============================================
-- 分类（一级/二级，支持上下架）
-- ============================================

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    parent_id INT REFERENCES categories(id),
    sort_order INT NOT NULL DEFAULT 0,
    course_count INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,  -- 上下架
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ============================================
-- 标签
-- ============================================

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(50) NOT NULL UNIQUE,
    course_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ============================================
-- 课程
-- ============================================

CREATE TABLE courses (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(500) NOT NULL,          -- 主标题
    subtitle VARCHAR(500),                 -- 副标题
    slug VARCHAR(200) NOT NULL UNIQUE,
    category_id INT NOT NULL REFERENCES categories(id),
    quality_label VARCHAR(50),
    cover_image VARCHAR(500),              -- 主图
    content TEXT,                          -- 商品详情（详情图HTML）
    detail_title VARCHAR(500),             -- 商品详情-主标题
    detail_subtitle VARCHAR(500),          -- 商品详情-副标题
    price INT NOT NULL DEFAULT 0,          -- 普通价格（金币）
    vip_price INT NOT NULL DEFAULT 0,      -- 会员价格（VIP免费显示0）
    status VARCHAR(20) NOT NULL DEFAULT 'draft', -- draft, published, trashed
    view_count INT NOT NULL DEFAULT 0,
    like_count INT NOT NULL DEFAULT 0,
    download_count INT NOT NULL DEFAULT 0,
    comment_count INT NOT NULL DEFAULT 0,
    author_id INT REFERENCES users(id),
    published_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_courses_category ON courses(category_id);
CREATE INDEX idx_courses_status ON courses(status);
CREATE INDEX idx_courses_published ON courses(published_at DESC) WHERE status = 'published';
CREATE INDEX idx_courses_slug ON courses(slug);

-- 课程-标签关联
CREATE TABLE course_tags (
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (course_id, tag_id)
);

-- 课程资源文件
CREATE TABLE course_resources (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    name VARCHAR(200),
    url VARCHAR(1000) NOT NULL,
    password VARCHAR(100),
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ============================================
-- 金币系统
-- ============================================

CREATE TABLE coin_accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) UNIQUE,
    balance INT NOT NULL DEFAULT 0,
    total_earned INT NOT NULL DEFAULT 0,
    total_spent INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE coin_transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    type VARCHAR(20) NOT NULL,    -- charge, spend, refund, gift
    amount INT NOT NULL,
    balance_after INT NOT NULL,
    description VARCHAR(200),
    order_id BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_coin_tx_user ON coin_transactions(user_id);

-- ============================================
-- VIP
-- ============================================

CREATE TABLE vip_packages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    price INT NOT NULL,           -- 当前价格（99金币）
    original_price INT,           -- 原价（699金币）
    duration_days INT,            -- NULL = 永久
    benefits JSONB,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE user_vip (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    package_id INT NOT NULL REFERENCES vip_packages(id),
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,         -- NULL = 永久
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_vip_user ON user_vip(user_id);

-- ============================================
-- 订单
-- ============================================

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(32) NOT NULL UNIQUE,
    user_id INT NOT NULL REFERENCES users(id),
    type VARCHAR(20) NOT NULL,    -- course, vip, coin
    target_id INT,
    target_name VARCHAR(200),
    amount INT NOT NULL,
    pay_method VARCHAR(20),       -- wechat, alipay, coin
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    paid_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_no ON orders(order_no);

-- 购买记录
CREATE TABLE purchases (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    course_id BIGINT NOT NULL REFERENCES courses(id),
    order_id BIGINT REFERENCES orders(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, course_id)
);

CREATE INDEX idx_purchases_user ON purchases(user_id);

-- ============================================
-- 下载日志
-- ============================================

CREATE TABLE download_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    course_id BIGINT NOT NULL REFERENCES courses(id),
    ip VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ============================================
-- 收藏
-- ============================================

CREATE TABLE favorites (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    course_id BIGINT NOT NULL REFERENCES courses(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, course_id)
);

-- ============================================
-- 留言板
-- ============================================

CREATE TABLE guestbook_messages (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    like_count INT NOT NULL DEFAULT 0,
    parent_id BIGINT REFERENCES guestbook_messages(id),
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(20) NOT NULL DEFAULT 'published',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_gb_status ON guestbook_messages(status);
CREATE INDEX idx_gb_created ON guestbook_messages(created_at DESC);

CREATE TABLE guestbook_likes (
    id BIGSERIAL PRIMARY KEY,
    message_id BIGINT NOT NULL REFERENCES guestbook_messages(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(message_id, user_id)
);

-- ============================================
-- 评论
-- ============================================

CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    target_type VARCHAR(20) NOT NULL, -- course
    target_id INT NOT NULL,
    content TEXT NOT NULL,
    parent_id BIGINT REFERENCES comments(id),
    status VARCHAR(20) NOT NULL DEFAULT 'published', -- published, hidden, pending
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comments_target ON comments(target_type, target_id);

-- ============================================
-- 工单系统
-- ============================================

CREATE TABLE tickets (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, processing, replied, closed
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE ticket_replies (
    id BIGSERIAL PRIMARY KEY,
    ticket_id BIGINT NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tickets_user ON tickets(user_id);
CREATE INDEX idx_tickets_status ON tickets(status);

-- ============================================
-- 系统配置
-- ============================================

CREATE TABLE settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ============================================
-- 管理员操作日志
-- ============================================

CREATE TABLE admin_logs (
    id BIGSERIAL PRIMARY KEY,
    admin_id INT NOT NULL REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    target_type VARCHAR(50),
    target_id INT,
    detail JSONB,
    ip VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ============================================
-- 支付异常日志
-- ============================================

CREATE TABLE payment_logs (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT REFERENCES orders(id),
    type VARCHAR(20) NOT NULL,    -- wechat, alipay
    status VARCHAR(20) NOT NULL,  -- success, fail, exception
    request_data JSONB,
    response_data JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ============================================
-- 金刚区导航配置
-- ============================================

CREATE TABLE nav_items (
    id SERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    subtitle VARCHAR(100),
    icon VARCHAR(500),
    link VARCHAR(500) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ============================================
-- Banner 配置
-- ============================================

CREATE TABLE banners (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100),
    image_url VARCHAR(500) NOT NULL,
    link VARCHAR(500),
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

---

## 7.3 索引策略

| 表 | 索引 | 用途 |
|----|------|------|
| courses | status + published_at | 列表查询 |
| courses | category_id | 分类筛选 |
| courses | slug | 详情页查询 |
| orders | user_id + status | 用户订单查询 |
| orders | order_no | 订单号查询 |
| guestbook_messages | status + created_at | 留言列表 |
| comments | target_type + target_id | 评论列表 |
| purchases | user_id + course_id | 购买检查 |
| tickets | user_id + status | 工单查询 |

---

## 7.4 缓存策略

| 数据 | 缓存 Key | TTL | 说明 |
|------|----------|-----|------|
| 首页课程列表 | `home:courses:latest` | 5min | 最新8条 |
| 首页分类Tab | `home:courses:category:{id}` | 5min | 每分类8条 |
| 课程详情 | `course:{slug}` | 10min | 课程完整信息 |
| 分类列表 | `categories:all` | 1h | 分类 |
| 标签云 | `tags:hot` | 1h | 热门标签 |
| VIP套餐 | `vip:packages` | 1h | 套餐列表 |
| 用户VIP状态 | `user:{id}:vip` | 5min | 是否VIP |
| 用户金币 | `user:{id}:coins` | 1min | 余额 |
| 验证码 | `captcha:{key}` | 5min | 图片验证码 |
| 金刚区 | `nav:items` | 1h | 金刚区配置 |
| Banner | `banners:active` | 1h | Banner列表 |
```
