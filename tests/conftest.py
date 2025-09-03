"""
pytest配置文件
"""
import pytest
import requests
from typing import Dict, Any


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
        
        def post(self, endpoint: str, data: Dict[Any, Any], headers: Dict[str, str] = None) -> requests.Response:
            """发送POST请求"""
            url = f"{self.base_url}{endpoint}"
            return self.session.post(url, json=data, headers=headers or {})
    
    return APIClient


@pytest.fixture(scope="session")
def test_user_id():
    """测试用户ID"""
    return 1


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
