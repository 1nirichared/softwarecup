-- 删除所有表（按依赖关系顺序）
USE teaching_platform;

-- 删除外键约束和表
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS learning_progress;
DROP TABLE IF EXISTS knowledge_base;
DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS chat_sessions;
DROP TABLE IF EXISTS exercise_records;
DROP TABLE IF EXISTS student_answers;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS exercises;
DROP TABLE IF EXISTS course_materials;
DROP TABLE IF EXISTS knowledge;
DROP TABLE IF EXISTS chapters;
DROP TABLE IF EXISTS courses;
DROP TABLE IF EXISTS user_profiles;
DROP TABLE IF EXISTS users;

SET FOREIGN_KEY_CHECKS = 1; 