CREATE TABLE user_page
(
    `id`            BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'auto increment ID',
    `uid`  BIGINT(20) DEFAULT 0 NOT NULL COMMENT 'user id',
    `pid`       VARCHAR(16) DEFAULT '' NOT NULL COMMENT 'origin page id，start with O',
    `sort`      INT(11) DEFAULT 0 NOT NULL COMMENT 'sort order of page',

    `deleted_at`    DATETIME                NULL,
    `created_at`    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `updated_at`    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_uid_pid` (`uid`, `pid`)
);

