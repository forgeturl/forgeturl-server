# 共享页面加入个人空间

2026-11-06，基线 `1b910c4`，开发分支 `feat/shared-page-space`。

## 行为与接口

复用 `user_page`，不新增表、不复制页面、不创建分享链接。空间入口精确保存原 R/E/A 分享 ID；同一页面的只读、编辑入口相互独立，权限仍由原链接决定。

所有写入均从登录态读取 uid，客户端不能指定目标用户。

| 接口 | 请求 | 响应 |
| --- | --- | --- |
| `POST /space/addSharedPage` | `{ "page_id": "E2cGtMA" }` | `page_brief`、`page_ids`、`already_added` |
| `POST /space/removeSharedPage` | `{ "page_id": "E2cGtMA" }` | `page_ids`，允许空数组 |

- 加入校验前缀、Base58 格式及链接当前是否可访问；空、非法或失效 ID 返回正常业务错误。
- 新入口追加到已有列表末尾；重复加入不写记录、不改变顺序。
- 自己页面的分享入口返回已有 O 页面摘要，`already_added=true`，不建立重复关联。
- 移除只物理删除当前 uid 与指定分享 ID 的关联；不存在也成功。失效链接可移除，不能通过该接口移除仍属于自己的页面。
- 加入使用唯一键 upsert 恢复旧软删除记录，移除后可再次加入。
- 撤销链接保留空间关联，`getMySpace` 返回 `is_not_found=true`、空权限占位。其他数据库错误继续返回。
- `getPage` 新增 `owner_page_id`，只对页面本人返回，供自己的分享页跳回个人空间。

## 排序并发

加入、移除、创建和排序均在事务内先对当前 user 行执行 `SELECT ... FOR UPDATE`，即使空间为空也有稳定的锁。不同页面并发加入会依次读取最新列表，不互相覆盖。加入校验分享时锁住 page 行，与撤销或删除协调。

`savePageIds` 新增字段：

```json
{
  "page_ids": ["Rabc", "Oabc"],
  "base_page_ids": ["Oabc", "Rabc"],
  "check_base": true
}
```

`base_page_ids` 必须等于服务器当前完整顺序，否则返回 409 冲突，客户端应刷新后重新排序。

旧客户端省略 `check_base` 时，按提交顺序排列，并在末尾保留未提交的已有入口。**遗漏 ID 不再代表移除**，以免旧列表覆盖并发加入。共享入口移除统一使用新接口，原页面删除仍使用 `deletePage`。

## 生成和验证

使用现有 `app/genapi.sh`，检查并使用 protoc 24.2、仓库要求的定制 `protoc-gen-go`、gin 1.0.3、errors/openapi 1.0.0、yq 和 Go 1.24.11。只保留本功能的生成差异；未手改生成的 Go/HTTP/Swagger。

通过：

- `go test ./...`
- 新接口、DAL、格式校验的定向 `go test -race`；MySQL SQL mock 验证用户行锁、分享行锁、精确分享 ID、幂等、追加排序、软删除恢复、物理移除、旧排序冲突、NotFound 占位及数据库错误传播。
- `go build -o /private/tmp/forgeturl-server-validation .`
- 配套前端的模拟 API 浏览器验收，详见前端 `docs/shared-page-space.md`。

新增持久化测试仅使用 SQL mock，不连接业务数据库。尚未进行真实 MySQL 并发压测、测试账号端到端联调或线上验收；Python 测试依赖实际 API/账号，本次未运行。

环境说明：shell 默认 Go 与项目版本不符，验证显式指定 Go 1.24.11 的 PATH/GOROOT。四处现有错误包装改为常量格式字符串，以通过该版本默认 vet。race 链接器产生本机 SQLite C 对象的 LC_DYSYMTAB 警告，测试退出码为 0。

## 部署

本次交付未部署。建议先发布后端，再发布配套前端；无需数据库变更。发布后用 A/B 测试账号复核持久化、共享内容双向更新、真实撤销以及数据库并发行为。
