create table page (
    id bigint NOT NULL AUTO_INCREMENT,
    readonly_pid varchar(32) default '' not null comment '只读权限地址',
    edit_pid varchar(32) default '' not null comment '编辑权限地址',
    admin_pid varchar(32) default '' not null comment '超级权限地址',
    title varchar(128) default '' not null,
    content text,

    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp,
    PRIMARY KEY (`id`),
    KEY idx_readonly_pid(`readonly_pid`),
    KEY idx_edit_pid(`edit_pid`)
)