# DROP TABLE page;

CREATE TABLE `page`
(
    `id`        BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
    `uid`       BIGINT DEFAULT 0 NOT NULL COMMENT 'page owner of uid',
    `pid`       VARCHAR(16) DEFAULT '' NOT NULL COMMENT 'origin page id，start with O',
    `readonly_pid` VARCHAR(16) DEFAULT '' NOT NULL COMMENT 'read only page id, start with R',
    `edit_pid`  VARCHAR(16) DEFAULT '' NOT NULL COMMENT 'edit page id, start with E',
    `admin_pid` VARCHAR(16) DEFAULT '' NOT NULL COMMENT 'super admin page id, start with A',
    `title`     VARCHAR(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT '标题',
    `brief`     VARCHAR(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT '简要描述',
    `content`   MEDIUMTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '实体内容(文件夹、链接定义)',
    `version`   BIGINT(20) DEFAULT 0 NOT NULL COMMENT '版本号',

    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    PRIMARY KEY (`id`),
    UNIQUE  KEY `uk_pid` (`pid`),
    KEY `idx_readonly_pid` (`readonly_pid`), -- 需要创建id时保持唯一
    KEY `idx_edit_pid` (`edit_pid`), -- 需要创建id时保持唯一
    KEY `idx_admin_pid` (`admin_pid`), -- 需要创建id时保持唯一
    KEY `idx_created_at` (`created_at`),
    KEY `idx_updated_at` (`updated_at`)
);




# ALTER TABLE `page`
#     MODIFY `title`     VARCHAR(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT '标题';
#
ALTER TABLE `page`
     MODIFY `content`    MEDIUMTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '实体内容(文件夹、链接定义)';