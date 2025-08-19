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

## 非登录态
```bash
# 获取用户信息
curl '127.0.0.1:80/space/getUserInfo' -d '{"uid": 1}' -H 'content-type: application/json'

# 获取页面数据
curl '127.0.0.1:80/space/getPage' -d '{"page_id": "test_page_id"}' -H 'content-type: application/json'
```

## 登录态
```bash
# 获取用户信息
curl '127.0.0.1:80/space/getUserInfo' -d '{"uid": 1}' -H 'content-type: application/json' -H 'X-Token: test'

# 拉取我的空间
curl '127.0.0.1:80/space/getMySpace' -d '{"uid": 1}' -H 'content-type: application/json' -H 'X-Token: test'

# 获取页面数据
curl '127.0.0.1:80/space/getPage' -d '{"page_id": "test_page_id"}' -H 'content-type: application/json' -H 'X-Token: test'

# 创建页面
curl '127.0.0.1:80/space/createPage' -d '{
  "title": "我的页面",
  "brief": "这是一个测试页面",
  "collections": [
    {
      "links": [
        {
          "title": "示例链接",
          "url": "https://example.com",
          "tags": ["示例", "测试"],
          "photo_url": "https://example.com/photo.jpg",
          "sub_links": [
            {
              "sub_title": "子链接",
              "sub_url": "https://sub.example.com"
            }
          ]
        }
      ]
    }
  ]
}' -H 'content-type: application/json' -H 'X-Token: test'

# 更新页面
curl '127.0.0.1:80/space/updatePage' -d '{
  "page_id": "O3sFmpq",
  "title": "更新后的页面标题",
  "brief": "更新后的页面描述",
  "collections": [
    {
      "links": [
        {
          "title": "更新后的链接",
          "url": "https://updated.example.com",
          "tags": ["更新", "测试"]
        }
      ]
    }
  ],
  "version": 0,
  "mask": 7
}' -H 'content-type: application/json' -H 'X-Token: test'

# 删除页面
curl '127.0.0.1:80/space/deletePage' -d '{"page_id": "test_page_id"}' -H 'content-type: application/json' -H 'X-Token: test'

# 创建页面链接
curl '127.0.0.1:80/space/createPageLink' -d '{
  "page_id": "test_page_id",
  "page_type": "readonly"
}' -H 'content-type: application/json' -H 'X-Token: test'

# 移除页面链接
curl '127.0.0.1:80/space/removePageLink' -d '{"page_id": "test_page_id"}' -H 'content-type: application/json' -H 'X-Token: test'

# 保存页面ID顺序
curl '127.0.0.1:80/space/savePageIds' -d '{
  "uid": 1,
  "page_ids": ["page1", "page2", "page3"]
}' -H 'content-type: application/json' -H 'X-Token: test'

# 创建临时页面（已废弃）
curl '127.0.0.1:80/space/createTmpPage' -d '{"user_uuid": "test_user_uuid"}' -H 'content-type: application/json' -H 'X-Token: test'
```

## 接口说明

### 页面类型 (page_type)
- `readonly`: 只读链接
- `edit`: 可编辑链接  
- `admin`: 超级权限链接

### 更新掩码 (mask)
- `0x01`: 更新标题 (title)
- `0x02`: 更新描述 (brief)
- `0x04`: 更新集合 (collections)
- `0x07`: 更新所有字段

### 页面配置 (PageConf)
- `read_only`: 是否只读
- `can_edit`: 是否可编辑
- `can_delete`: 是否可删除