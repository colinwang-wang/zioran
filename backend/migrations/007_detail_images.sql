-- Add detail_images column to courses table (JSON array of image URLs)
ALTER TABLE courses ADD COLUMN detail_images TEXT DEFAULT NULL AFTER content;
