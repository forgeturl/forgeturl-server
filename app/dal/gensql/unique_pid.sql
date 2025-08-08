CREATE TABLE unique_pid
(
    `id`            BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'auto increment ID',
    `pid`           VARCHAR(16) DEFAULT '' NOT NULL COMMENT 'unique page id',
    `uid`           BIGINT(20) DEFAULT 0 NOT NULL COMMENT 'user id',
    `created_at`    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `updated_at`    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_pid` (`pid`)
);

