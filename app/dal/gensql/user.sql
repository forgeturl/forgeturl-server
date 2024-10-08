create table user(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    display_name varchar(128) default '' not null,
    username varchar(64) default '' not null,
    email varchar(100) not null,
    avatar          varchar(1024) default '' not null,
    status          int           default 0  not null, -- user status(normal 0,suspended 2,deleted 4)
    last_login_date timestamp                null,

    page_ids varchar(2048) default '' not null ,

    provider    varchar(32) default '' not null,
    external_id varchar(128) default '' not null,
    ip_info         varchar(255)  default '' not null,
    is_admin        tinyint(1)    default 0  not null,

    suspended_at    timestamp                null,
    deleted_at      timestamp                null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);


CREATE UNIQUE INDEX uk_username ON user(username);
CREATE UNIQUE INDEX uk_email ON user(email);
