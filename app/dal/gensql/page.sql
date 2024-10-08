create table page (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uid INTEGER DEFAULT 0 NOT NULL, -- 页面归属的用户
    pid varchar(32) default '' not null, -- 页面id，使用随机数生成
    readonly_pid varchar(32) default '' not null,-- comment '只读权限地址',
    edit_pid varchar(32) default '' not null, -- comment '编辑权限地址',
    admin_pid varchar(32) default '' not null, -- comment '超级权限地址(待定，暂不支持)',
    title varchar(128) default '' not null,
    content mediumtext default '' not null, -- 链接实体

    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

CREATE INDEX idx_pid ON page(pid);
CREATE INDEX idx_readonly_pid ON page(readonly_pid);
CREATE INDEX idx_edit_pid ON page(edit_pid);
