# 测试用例 Trace ID 功能使用说明

## 功能概述

我们已经优化了测试用例，现在所有的 HTTP 请求都会自动打印响应头中的 `x-b3-traceid`，方便进行问题跟踪和调试。

## 功能特性

### 1. 自动 Trace ID 提取
- 支持多种常见的 trace header 格式：
  - `x-b3-traceid`
  - `X-B3-TraceId`  
  - `X-B3-TRACEID`
  - `x-trace-id`
  - `X-Trace-Id`

### 2. 两种日志模式

#### 基本模式 (默认)
使用 `api_client.post()` 方法时，会打印基本的 trace 信息：

```
============================================================
[TRACE INFO] 接口: /space/createPage
[TRACE INFO] TraceID: a1b2c3d4e5f6g7h8
[TRACE INFO] 状态码: 200
[TRACE INFO] 响应时间: 0.123s
============================================================
```

#### 详细模式
使用 `api_client.post_with_detailed_log()` 方法时，会打印完整的请求和响应信息：

```
============================================================
[REQUEST] 接口: /space/createPage
[REQUEST] 请求数据: {
  "title": "我的页面",
  "brief": "这是一个测试页面"
}
[REQUEST] 请求头: {
  "Content-Type": "application/json",
  "X-Token": "test"
}
[RESPONSE] 状态码: 200
[RESPONSE] TraceID: a1b2c3d4e5f6g7h8
[RESPONSE] 响应时间: 0.123s
[RESPONSE] 响应内容: {
  "code": 1,
  "page_id": "new_page_123"
}
============================================================
```

## 使用方法

### 在现有测试中
所有现有测试用例无需修改，已经自动支持基本的 trace ID 打印。

### 新增测试用例

#### 基本用法（推荐）
```python
def test_example(self, headers_with_login):
    # 自动打印 trace ID
    response = self.api_client.post("/space/getUserInfo", data, headers_with_login)
    # 验证响应...
```

#### 详细日志用法（调试时使用）
```python
def test_example_detailed(self, headers_with_login):
    # 打印完整的请求和响应信息
    response = self.api_client.post_with_detailed_log("/space/getUserInfo", data, headers_with_login)
    # 验证响应...
```

### 手动使用 TraceUtils
```python
from test_utils import TraceUtils

def test_manual_trace(self):
    # 发送请求
    response = requests.post(url, json=data, headers=headers)
    
    # 手动提取和打印 trace ID
    trace_id = TraceUtils.extract_trace_id(response)
    trace_info = TraceUtils.format_trace_info("/api/endpoint", response)
    print(trace_info)
    
    # 或者打印详细信息
    TraceUtils.log_request_details("/api/endpoint", data, headers, response)
```

## 运行测试

### 运行所有测试（基本日志模式）
```bash
pytest tests/test_space_api.py -v
```

### 运行特定测试（查看详细日志）
```bash
pytest tests/test_space_api.py::TestSpaceAPI::test_05b_create_page_with_login_success -v -s
```

### 启用所有日志输出
```bash
pytest tests/test_space_api.py -v -s --log-cli-level=INFO
```

## 注意事项

1. **性能考虑**：详细日志模式会打印完整的请求和响应信息，在大量测试时可能影响性能，建议仅在调试时使用。

2. **响应内容限制**：为避免日志过长，只有小于 1000 字符的响应内容会被打印。

3. **Trace ID 缺失**：如果服务端没有返回 trace ID，会显示"未找到"，这是正常情况。

4. **日志级别**：基本的 trace 信息会通过 `print()` 输出到控制台，详细的日志信息会记录到 pytest 的日志系统中。

## 示例输出

运行测试时，你会看到类似以下的输出：

```
tests/test_space_api.py::TestSpaceAPI::test_05b_create_page_with_login_success 

============================================================
[REQUEST] 接口: /space/createPage
[REQUEST] 请求数据: {
  "title": "我的页面",
  "brief": "这是一个测试页面",
  "collections": [...]
}
[REQUEST] 请求头: {
  "Content-Type": "application/json",
  "X-Token": "test"
}
[RESPONSE] 状态码: 200
[RESPONSE] TraceID: f47ac10b-58cc-4372-a567-0e02b2c3d479
[RESPONSE] 响应时间: 0.245s
============================================================

PASSED
```

这样你就可以轻松地跟踪每个 API 请求的 trace ID，方便进行问题定位和性能分析。
