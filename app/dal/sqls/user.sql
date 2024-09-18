create table user(
    id bigint NOT NULL AUTO_INCREMENT,
    username varchar(32) not null,
    password varchar(32) not null,
    email varchar(128) not null,

    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp,
    PRIMARY KEY (`id`),
    UNIQUE KEY uk_username(`username`),
    UNIQUE KEY uk_email(`email`)
);