CREATE TABLE page (
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `uid` BIGINT DEFAULT 0 NOT NULL COMMENT '页面归属的用户',
    `pid` VARCHAR(32) DEFAULT '' NOT NULL COMMENT '页面id，使用随机数生成',
    `readonly_pid` VARCHAR(32) DEFAULT '' NOT NULL COMMENT '只读权限地址',
    `edit_pid` VARCHAR(32) DEFAULT '' NOT NULL COMMENT '编辑权限地址',
    `admin_pid` VARCHAR(32) DEFAULT '' NOT NULL COMMENT '超级权限地址(待定，暂不支持)',
    `title` VARCHAR(128) DEFAULT '' not null,
    `content` MEDIUMTEXT DEFAULT '' NOT NULL COMMENT '链接实体',

    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_pid` (`pid`),
    UNIQUE KEY `uk_readonly_pid` (`readonly_pid`),
    UNIQUE KEY `uk_edit_pid` (`edit_pid`),
    UNIQUE KEY `uk_admin_pid` (`admin_pid`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_updated_at` (`updated_at`)
);

CREATE INDEX idx_pid ON page(pid);
CREATE INDEX idx_readonly_pid ON page(readonly_pid);
CREATE INDEX idx_edit_pid ON page(edit_pid);
