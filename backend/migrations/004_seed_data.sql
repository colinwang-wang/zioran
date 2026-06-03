-- VIP套餐
INSERT IGNORE INTO vip_packages (id, name, price, original_price, duration_days, is_active, sort_order) 
VALUES (1, '终身VIP', 99, 699, NULL, 1, 1);

-- 分类
INSERT IGNORE INTO categories (id, name, slug, sort_order, is_active) VALUES
(1, 'AIGC课堂', 'aigc', 1, 1),
(2, 'Blender课堂', 'blender', 2, 1),
(3, 'C4D课程', 'c4d', 3, 1),
(4, '手绘课程', 'drawing', 4, 1),
(5, 'AE课程', 'ae', 5, 1),
(6, 'UI课程', 'ui', 6, 1),
(7, '摄影课程', 'photography', 7, 1),
(8, '室内设计', 'interior', 8, 1),
(9, '平面设计', 'graphic', 9, 1),
(10, '电商设计', 'ecommerce', 10, 1),
(11, '3dmax课程', '3dmax', 11, 1),
(12, 'zbrush课程', 'zbrush', 12, 1),
(13, 'ai课程', 'ai', 13, 1),
(14, '视频课程', 'video', 14, 1);
