# 02 - 课程模块

## 2.1 课程数据模型

### courses 表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | - |
| title | VARCHAR(500) | 课程标题（主标题） |
| subtitle | VARCHAR(500) | 副标题 |
| slug | VARCHAR(200) UNIQUE | URL 标识 |
| category_id | INT FK | 分类 |
| quality_label | VARCHAR(50) | 质量标注 |
| cover_image | VARCHAR(500) | 主图 URL（按参考网站尺寸） |
| content | TEXT | 商品详情（详情图 HTML） |
| price | INT DEFAULT 0 | 普通价格（金币） |
| vip_price | INT DEFAULT 0 | 会员价格（VIP免费，前端显示为 0） |
| status | VARCHAR(20) DEFAULT 'draft' | draft/published/pending/trashed |
| view_count | INT DEFAULT 0 | - |
| like_count | INT DEFAULT 0 | 收藏数 |
| download_count | INT DEFAULT 0 | - |
| comment_count | INT DEFAULT 0 | - |
| author_id | INT FK | 发布者 |
| published_at | TIMESTAMP | - |
| created_at | TIMESTAMP | - |
| updated_at | TIMESTAMP | - |

### course_tags 关联表
| 字段 | 类型 |
|------|------|
| course_id | BIGINT FK |
| tag_id | INT FK |
| PRIMARY KEY | (course_id, tag_id) |

### course_resources 资源表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | - |
| course_id | BIGINT FK | - |
| name | VARCHAR(200) | 资源名 |
| url | VARCHAR(1000) | 下载地址 |
| password | VARCHAR(100) | 提取码 |
| sort_order | INT DEFAULT 0 | - |

### categories 表（一级/二级分类）
| 字段 | 类型 |
|------|------|
| id | SERIAL PK |
| name | VARCHAR(50) |
| slug | VARCHAR(50) UNIQUE |
| description | TEXT |
| parent_id | INT FK (self) |
| sort_order | INT DEFAULT 0 |
| course_count | INT DEFAULT 0 |
| is_active | BOOLEAN DEFAULT TRUE |

### 分类（后台可管理，金刚区与课堂二级导航对应）
| 名称 | Slug | 说明 |
|------|------|------|
| AIGC课堂 | aigc | 金刚区入口 |
| Blender课堂 | blender | 金刚区入口 |
| 3dmax课程 | 3dmax | - |
| AE课程 | ae | - |
| ai课程 | ai | - |
| C4D课程 | c4d | - |
| UI课程 | ui | - |
| zbrush课程 | zbrush | - |
| 室内设计 | interior | - |
| 平面设计 | graphic | - |
| 手绘课程 | drawing | - |
| 摄影课程 | photography | - |
| 电商设计 | ecommerce | - |
| 视频课程 | video | - |

> 金刚区导航与分类对应，主副标题后台可添加。后台新增导航模块时，前端链接对应内容模块，内容后台可调整，模块框架不变。

### tags 表
| 字段 | 类型 |
|------|------|
| id | SERIAL PK |
| name | VARCHAR(50) |
| slug | VARCHAR(50) UNIQUE |
| course_count | INT DEFAULT 0 |

### 质量标注枚举
```
画质高清有课件笔刷
画质高清有课件
画质高清只有视频
画质高清有素材
画质高清有大部分课件和笔刷
画质不错有课件笔刷
画质不错只有视频
画质不错有素材
```

---

## 2.2 课程列表页

### 路由
- `/courses` — 全部课堂
- `/courses/category/{slug}` — 分类筛选（AIGC课堂/Blender课堂等）
- `/courses/tag/{slug}` — 标签筛选
- `/new` — 最新发布

### 页面结构
```
body
├── header.header
├── div.main
│   └── div.container
│       ├── div.breadcrumbs（仅分类/标签页）
│       └── div#posts.posts.grids
│           └── div.post.grid × 16   # 课程卡片
│               ├── div.img
│               │   └── a > img.thumb（主图，按参考网站尺寸）
│               └── div.con
│                   ├── div.cat > a  # 分类标签
│                   ├── h3 > a       # 标题
│                   ├── div.excerpt   # 摘要
│                   └── div.grid-meta
│                       ├── span.time  # 相对时间
│                       └── span.price # 价格
├── div.pagination
├── footer.footer
└── div.footer-fixed-nav
```

### 列表布局
- 按发布时间显示8个商品（上四下四布局）
- 每页16个课程卡片
- "查看更多"链接到对应的全部课程页面

### API
```
GET /api/v1/courses
Query: page=1, pageSize=16, categoryId, tagId, sort=latest|popular|downloads, keyword
Response: {
  "code": 0,
  "data": {
    "items": [{
      "id": 25286,
      "title": "课程主标题",
      "subtitle": "课程副标题",
      "slug": "xxx",
      "cover": "主图URL",
      "category": { "id": 11, "name": "手绘课程", "slug": "drawing" },
      "price": 2,
      "vip_price": 0,
      "relative_time": "12小时前",
      "published_at": "2026-06-02T12:00:00Z"
    }],
    "total": 3200,
    "page": 1,
    "pageSize": 16,
    "totalPages": 200
  }
}
```

---

## 2.3 课程详情页

### 路由
`/courses/{slug}`

### 页面结构
```
body.single
├── header.header
├── div.main
│   └── div.container
│       ├── div.breadcrumbs
│       ├── div.content (主内容区，左侧)
│       │   └── article
│       │       ├── div.article-header
│       │       │   ├── h1  # 主标题
│       │       │   ├── p.subtitle  # 副标题
│       │       │   └── div.article-meta（日期、分类）
│       │       ├── div.article-content
│       │       │   ├── p > img  # 详情图（多张）
│       │       │   └── div.download-box
│       │       │       ├── "资源下载"
│       │       │       ├── 下载价格 X 金币
│       │       │       ├── 终身VIP免费
│       │       │       └── [立即购买] / [免费下载]
│       │       ├── div.article-act（收藏）
│       │       ├── div.article-tags
│       │       ├── div.article-shares
│       │       ├── nav.article-nav（上一篇/下一篇）
│       │       ├── div.related-posts（猜你喜欢 × 3）
│       │       └── div#comments（评论区）
│       └── aside.sidebar (右侧)
│           ├── div.widget-price（价格+购买按钮）
│           └── div.widget-tags（热门标签）
├── footer.footer
└── div.footer-fixed-nav
```

### 商品详情结构
```
主标题：课程名称
副标题：课程副描述
商品价格：
  - 普通价格：X 金币
  - 会员价格：注册成为会员，前端页面价格显示为 0
商品详情：
  - 主标题
  - 副标题
  - 详情图（多张预览图）
```

### 下载区状态
```
未登录:
  "请先登录"

已登录非VIP未购买:
  下载价格 X 金币
  终身VIP免费
  [立即购买]

VIP:
  终身VIP免费
  [免费下载]（价格显示为 0）

已购买:
  [下载]
  提取码: xxxx
```

### API
```
GET /api/v1/courses/{slug}
Response: {
  "code": 0,
  "data": {
    "id": 18277,
    "title": "主标题",
    "subtitle": "副标题",
    "slug": "xxx",
    "cover": "主图URL",
    "content": "详情图HTML",
    "price": 2,
    "vip_price": 0,
    "quality_label": "画质高清有课件笔刷",
    "category": { "id": 11, "name": "手绘课程", "slug": "drawing" },
    "tags": [{ "id": 1, "name": "PS课程", "slug": "ps" }],
    "like_count": 53,
    "comment_count": 0,
    "published_at": "2024-09-19T00:00:00Z",
    "prev_course": { "slug": "xxx", "title": "xxx" },
    "next_course": { "slug": "xxx", "title": "xxx" },
    "related_courses": [{ ... }],
    "user_access": {
      "can_download": false,
      "has_purchased": false,
      "is_vip": false,
      "is_favorited": false
    },
    "resources": []  // 仅 can_download=true 时返回
  }
}
```

---

## 2.4 搜索

### 说明
实现本站内的搜索功能。

### 搜索位置
- Header 右侧搜索图标 → 点击展开搜索弹层
- 首页搜索栏（Banner 下方）

### API
```
GET /api/v1/search?q=关键字&categoryId=&page=1&pageSize=16
Response: 同课程列表
```

---

## 2.5 收藏

### 交互
- 点击文章底部爱心图标 → 收藏/取消
- 需登录
- 数字实时更新

### API
```
POST /api/v1/courses/{id}/like
Response: { "code": 0, "data": { "liked": true, "like_count": 54 } }
```

---

## 2.6 管理后台 - 课程管理

### 课堂分类管理
- 一级/二级分类
- 排序（拖拽排序）
- 上下架

### 课堂编辑
| 字段 | 说明 |
|------|------|
| 主图 | 尺寸按参考网站尺寸 |
| 标题 | 课程主标题 |
| 副标题 | 课程副标题 |
| 商品价格 - 普通价格 | 非VIP用户购买价格（金币） |
| 商品价格 - 会员价格 | VIP免费，前端显示 0 |
| 商品详情 - 主标题 | 详情页内的标题 |
| 商品详情 - 副标题 | 详情页内的副标题 |
| 商品详情 - 详情图 | 课程预览图（多张） |
| 分类 | 一级/二级分类选择 |
| 标签 | 多选输入 |
| 质量标注 | 下拉选择 |
| 状态 | 发布/草稿/上架/下架 |
| 资源 | 网盘链接 + 提取码（可多个） |

### API
```
GET    /api/v1/admin/courses?page=&pageSize=&categoryId=&status=&keyword=
POST   /api/v1/admin/courses
PUT    /api/v1/admin/courses/{id}
DELETE /api/v1/admin/courses/{id}
PUT    /api/v1/admin/courses/{id}/status
POST   /api/v1/admin/courses/batch
```
