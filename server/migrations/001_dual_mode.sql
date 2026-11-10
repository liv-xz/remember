-- Migration 001: 双模式架构 - 数据库表结构变更
-- 执行前建议备份数据库
-- 执行方式: mysql -u root -p RememberMe < server/migrations/001_dual_mode.sql

-- 1. graves 表新增 grave_mode 字段
ALTER TABLE graves ADD COLUMN grave_mode TINYINT NOT NULL DEFAULT 1 COMMENT '1=纪念模式 2=族谱模式';
-- 所有现有数据默认 grave_mode=1（纪念模式），后续可手动升级

-- 2. deceaseds 表新增 origin_grave_id 字段
ALTER TABLE deceaseds ADD COLUMN origin_grave_id BIGINT UNSIGNED NULL COMMENT '族谱模式专用：嫁入者来源墓园ID';

-- 3. 重命名 user_deceased_relations → kinship_relations
RENAME TABLE user_deceased_relations TO kinship_relations;

-- 4. 重命名 deceased_relations → kinship_deceased_relations
RENAME TABLE deceased_relations TO kinship_deceased_relations;

-- 5. 新建 memorial_relations 表
CREATE TABLE memorial_relations (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at          DATETIME(3) NOT NULL,
    updated_at          DATETIME(3) NOT NULL,
    deleted_at          DATETIME(3) NULL,
    grave_id            BIGINT UNSIGNED NOT NULL,
    user_id             BIGINT UNSIGNED NOT NULL,
    deceased_id         BIGINT UNSIGNED NOT NULL,
    relation_type       TINYINT NOT NULL COMMENT '1=儿子 2=女儿 3=配偶 4=兄弟 5=姐妹 6=孙辈 7=父亲 8=母亲 9=祖父 10=祖母 11=长辈 99=追随者',
    confirmation_status TINYINT NOT NULL DEFAULT 0 COMMENT '0=待确认 1=已确认 2=已拒绝',
    confirmed_by        BIGINT UNSIGNED NULL,
    note                VARCHAR(255) NOT NULL DEFAULT '',
    INDEX idx_mr_grave (grave_id),
    INDEX idx_mr_deceased (deceased_id),
    UNIQUE INDEX idx_mr_user_deceased_grave (user_id, deceased_id, grave_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
