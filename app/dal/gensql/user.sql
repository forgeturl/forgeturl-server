CREATE TABLE `user`
(
    `id`            BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'auto increment ID',
    `display_name`  VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT 'display name of user',
    `username`      VARCHAR(64) DEFAULT '' NOT NULL,
    `email`         VARCHAR(100) NOT NULL COMMENT 'email from provider',
    `avatar`        VARCHAR(1024) DEFAULT '' NOT NULL COMMENT 'user avatar url',
    `status`        TINYINT   DEFAULT 0  NOT NULL COMMENT 'user status(normal 0,suspended 2,deleted 4)',
    `last_login_date` DATETIME NOT NULL,

    `provider`      VARCHAR(32)  DEFAULT '' NOT NULL COMMENT 'login source google/facebook/weixin',
    `external_id`   VARCHAR(255) DEFAULT '' NOT NULL COMMENT 'login source unique id(gmail sub 255char//weixin unionid 28char)',
    `ip_info`       VARCHAR(255) DEFAULT '' NOT NULL,
    `is_admin`      TINYINT      DEFAULT 0  NOT NULL,

    `suspended_at`  DATETIME                NULL,
    `deleted_at`    DATETIME                NULL,
    `created_at`    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `updated_at`    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_provider_external_id` (`provider`, `external_id`)
);


# ALTER TABLE `user`
#     MODIFY  `display_name`  VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT 'display name of user';
