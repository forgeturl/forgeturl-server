create table space(
    id bigint NOT NULL AUTO_INCREMENT,
    uid bigint default 0 not null,
    name varchar(255) default '' not null,
    description text,
    page_ids varchar(2048) default '' not null ,

    created_at timestamp default current_timestamp not null ,
    updated_at timestamp default current_timestamp on update current_timestamp,
    PRIMARY KEY (`id`),
    KEY idx_uid(`uid`)
);