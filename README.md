# 概述

## 非登录态
非登录态：
getPage 
getUserInfo(暂无场景)

# 登录相关
/auth/{name}
/callback/{name}
getUserInfo

# 登录态
getUserInfo
getMySpace
getPage

createPage
updatePage
deletePage
createPageLink
removePageLink

savePageIds

# 测试方法

```bash
# 非登录态
curl '127.0.0.1:80/space/getUserInfo' -d '{"uid": 1}' -H 'content-type: application/json'

# 登录态
curl '127.0.0.1:80/space/getUserInfo' -d '{"uid": 1}' -H 'content-type: application/json' -H 'X-Token: test'

```