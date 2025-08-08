# 概述

系统的表关系核心是：

- `user` 是主体；
- `page` 是资源，由用户创建；
- `unique_pid` 是防止 ID 冲突的“ID 分配池”；
- `user_page` 是用户与页面之间的“展示关系”映射表，支持排序。

> 🎯 最终图示应体现：  
> **用户 → 创建页面 → 页面的各类 PID 必须在 unique_pid 中唯一 → 用户空间通过 user_page 管理页面顺序**



```text
+------------------+          +-----------------------+
|      user        |          |     unique_pid        |
|------------------|          |-----------------------|
| id (PK)          |<---------| uid                   |
| username         |          | pid (PK, UK)          |
| email            |          | created_at            |
| display_name     |          | updated_at            |
| ...              |          +-----------------------+
+------------------+                   ▲
     |                                 |
     | 1                               | N (逻辑引用)
     |                                 |
     v                                 |
+------------------+       +-----------------------+
|      page        |-------| (pid, readonly_pid,   |
|------------------|       |  edit_pid, admin_pid) |
| id (PK)          |       | 全部必须存在于         |
| uid (FK)         |       | unique_pid.pid 中      |
| pid (UK)         |-----------------------------+
| readonly_pid     |
| edit_pid         |
| admin_pid        |
| title            |
| content          |
| ...              |
+------------------+
     ▲
     |
     | N
     |
+------------------+
|   user_page      |
|------------------|
| id (PK)          |
| uid (FK)         |-----> user.id
| pid (FK)         |-----> page.pid
| sort             |
| deleted_at       |
| created_at       |
| updated_at       |
+------------------+


```