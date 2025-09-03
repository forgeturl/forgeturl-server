# ForgetURL Server API 测试

基于 `requests` 和 `pytest-allure` 库的 HTTP 接口测试套件，用于测试 ForgetURL 服务器的 Space API 功能。

## 环境要求

- Python 3.8+
- 服务器运行在 `http://127.0.0.1:80`

## 安装依赖

```bash
cd /Users/lxy/Desktop/Git/forgeturl-server/app/tests
pip install -r requirements.txt
```

## 运行测试

### 运行所有测试
```bash
pytest
```

### 运行指定标记的测试
```bash
# 只运行需要登录态的测试
pytest -m login

# 只运行非登录态的测试  
pytest -m no_login

# 运行空间相关测试
pytest -m space

# 运行页面相关测试
pytest -m page
```

### 生成 Allure 报告
```bash
# 运行测试并生成 allure 数据
pytest --alluredir=allure-results

# 生成并打开 allure 报告
allure serve allure-results
```

### 生成 HTML 报告
```bash
pytest --html=report.html --self-contained-html
```

## 测试用例说明

测试用例基于 README.md 中定义的 11 个场景：

1. **用户信息获取**
   - 非登录态获取基本用户信息
   - 登录态获取详细用户信息

2. **页面访问**
   - 非登录态获取公开页面
   - 登录态获取私有页面

3. **我的空间**
   - 非登录态无法访问个人空间
   - 登录态可以访问个人空间

4. **页面管理**
   - 创建页面（登录态成功/非登录态失败）
   - 更新页面（登录态成功/非登录态失败）
   - 删除页面（登录态成功/非登录态失败）

5. **页面链接管理**
   - 创建页面链接（登录态成功/非登录态失败）
   - 移除页面链接（登录态成功/非登录态失败）

6. **页面顺序管理**
   - 调整页面顺序（登录态成功/非登录态失败）

## 测试标记

- `@pytest.mark.login`: 需要登录态的测试
- `@pytest.mark.no_login`: 非登录态的测试  
- `@pytest.mark.space`: 空间相关的测试
- `@pytest.mark.page`: 页面相关的测试

## 测试配置

- **测试Token**: `X-Token: test`
- **测试用户ID**: `1`
- **服务器地址**: `http://127.0.0.1:80`

## 文件结构

```
tests/
├── conftest.py           # pytest 配置和 fixtures
├── test_space_api.py     # 主要测试用例
├── requirements.txt      # 依赖包列表
├── pytest.ini          # pytest 配置文件
└── README.md            # 本文件
```

## 注意事项

1. 确保服务器在运行测试前已启动
2. 测试用例会创建、修改和删除测试数据
3. 部分测试用例之间存在依赖关系（如先创建页面再删除）
4. 建议在测试环境中运行，避免影响生产数据
