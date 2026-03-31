CREATE TABLE `openclaw_api_key`
(
    `id`         BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
    `uid`        BIGINT(20) NOT NULL COMMENT 'user id',
    `api_key`    VARCHAR(128) DEFAULT '' NOT NULL COMMENT 'openclaw api key plain text',
    `status`     TINYINT      DEFAULT 1  NOT NULL COMMENT '1 enabled, 0 disabled',
    `expires_at` DATETIME              NULL COMMENT 'nullable expiration',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_api_key` (`api_key`),
    KEY `idx_uid` (`uid`)
);
