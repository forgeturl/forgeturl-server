"""
pytest配置文件
"""
import pytest
import requests
import logging
from typing import Dict, Any
from .test_utils import TraceUtils


@pytest.fixture(scope="session")
def base_url():
    """基础URL"""
    return "http://127.0.0.1:80"


@pytest.fixture
def headers_no_login():
    """非登录态请求头"""
    return {
        "Content-Type": "application/json"
    }


@pytest.fixture
def headers_with_login():
    """登录态请求头"""
    return {
        "Content-Type": "application/json",
        "X-Token": "test"
    }


@pytest.fixture
def api_client():
    """API客户端"""
    class APIClient:
        def __init__(self, base_url: str):
            self.base_url = base_url
            self.session = requests.Session()
            # 设置日志器
            self.logger = logging.getLogger(__name__)
        
        def _print_trace_id(self, response: requests.Response, endpoint: str):
            """打印响应头中的x-b3-traceid"""
            trace_id = TraceUtils.extract_trace_id(response)
            if trace_id:
                # 使用TraceUtils格式化输出
                trace_info = TraceUtils.format_trace_info(endpoint, response, trace_id)
                print(trace_info)
                # 记录到日志
                self.logger.info(f"API请求 {endpoint} - TraceID: {trace_id} - 状态码: {response.status_code}")
            else:
                # 即使没有trace id，也打印基本信息
                print(f"\n[INFO] 接口: {endpoint} - 状态码: {response.status_code} - TraceID: 未找到")
        
        def post(self, endpoint: str, data: Dict[Any, Any], headers: Dict[str, str] = None) -> requests.Response:
            """发送POST请求"""
            url = f"{self.base_url}{endpoint}"
            response = self.session.post(url, json=data, headers=headers or {})
            
            # 打印trace id
            self._print_trace_id(response, endpoint)
            
            return response
        
        def post_with_detailed_log(self, endpoint: str, data: Dict[Any, Any], headers: Dict[str, str] = None) -> requests.Response:
            """发送POST请求，带详细日志"""
            url = f"{self.base_url}{endpoint}"
            response = self.session.post(url, json=data, headers=headers or {})
            
            # 打印详细的请求信息
            TraceUtils.log_request_details(endpoint, data, headers or {}, response)
            
            return response
    
    return APIClient


@pytest.fixture(scope="session") 
def test_user_id():
    """测试用户ID"""
    return 1


@pytest.fixture
def enable_detailed_logging():
    """是否启用详细日志 - 可在测试中通过 pytest -k "详细" 等方式控制"""
    return False  # 默认关闭，可通过环境变量或pytest参数控制


@pytest.fixture
def sample_page_data():
    """示例页面数据"""
    return {
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
    }


@pytest.fixture
def updated_page_data():
    """更新后的页面数据"""
    return {
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
    }
