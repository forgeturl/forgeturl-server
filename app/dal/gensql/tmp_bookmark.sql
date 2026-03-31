CREATE TABLE `tmp_bookmark`
(
    `id`             BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
    `uid`            BIGINT(20) NOT NULL COMMENT 'user id',
    `title`          VARCHAR(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT 'bookmark title',
    `title_hash`     CHAR(64)  DEFAULT '' NOT NULL COMMENT 'sha256 of title(lower+trim)',
    `url`            VARCHAR(2048) DEFAULT '' NOT NULL COMMENT 'raw url',
    `normalized_url` VARCHAR(2048) DEFAULT '' NOT NULL COMMENT 'normalized url',
    `url_hash`       CHAR(64)  DEFAULT '' NOT NULL COMMENT 'sha256 of normalized url',
    `status`         TINYINT   DEFAULT 1  NOT NULL COMMENT '1 active, 0 deleted',
    `created_at`     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `updated_at`     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_uid_title_url_status` (`uid`, `title_hash`, `url_hash`, `status`),
    KEY `idx_uid_status` (`uid`, `status`),
    KEY `idx_url_hash` (`url_hash`)
);
