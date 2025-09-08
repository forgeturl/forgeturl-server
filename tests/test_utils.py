"""
测试工具类
提供测试过程中需要的通用功能
"""
import json
import requests
from typing import Dict, Any, Optional


class DataManager:
    """测试数据管理器"""
    
    def __init__(self):
        self._shared_data = {}
    
    def set_data(self, key: str, value: Any) -> None:
        """设置共享数据"""
        self._shared_data[key] = value
    
    def get_data(self, key: str, default: Any = None) -> Any:
        """获取共享数据"""
        return self._shared_data.get(key, default)
    
    def clear_data(self) -> None:
        """清空所有数据"""
        self._shared_data.clear()


class ResponseValidator:
    """响应验证器"""
    
    @staticmethod
    def validate_success_response(response_data: Dict[str, Any], 
                                required_fields: list = None) -> bool:
        """验证成功响应"""
        if required_fields is None:
            required_fields = []
        
        # 检查是否包含必需字段
        for field in required_fields:
            if field not in response_data:
                return False
        
        # 检查错误码（如果存在）
        if "code" in response_data:
            return response_data["code"] == 1
        
        return True
    
    @staticmethod
    def validate_error_response(response_data: Dict[str, Any]) -> bool:
        """验证错误响应"""
        if "code" in response_data:
            return response_data["code"] != 1
        return False
    
    @staticmethod
    def validate_user_info_response(response_data: Dict[str, Any], 
                                   is_login: bool = False) -> bool:
        """验证用户信息响应"""
        basic_fields = ["uid", "display_name", "avatar", "email"]
        login_fields = ["username", "status", "create_time", "update_time"]
        
        # 检查基本字段
        for field in basic_fields:
            if field not in response_data:
                return False
        
        # 如果是登录态，检查额外字段
        if is_login:
            for field in login_fields:
                if field not in response_data:
                    return False
        
        return True
    
    @staticmethod
    def validate_page_response(response_data: Dict[str, Any]) -> bool:
        """验证页面响应"""
        if "page" not in response_data:
            return False
        
        page_data = response_data["page"]
        required_fields = ["page_id", "title", "collections"]
        
        for field in required_fields:
            if field not in page_data:
                return False
        
        return True


class TraceUtils:
    """链路追踪工具类"""
    
    @staticmethod
    def extract_trace_id(response: requests.Response) -> Optional[str]:
        """从响应头中提取trace id"""
        # 尝试多种可能的header名称
        possible_headers = ['x-b3-traceid', 'X-B3-TraceId', 'X-B3-TRACEID', 'x-trace-id', 'X-Trace-Id']
        for header in possible_headers:
            trace_id = response.headers.get(header)
            if trace_id:
                return trace_id
        return None
    
    @staticmethod
    def format_trace_info(endpoint: str, response: requests.Response, trace_id: str = None) -> str:
        """格式化trace信息输出"""
        if trace_id is None:
            trace_id = TraceUtils.extract_trace_id(response)
        
        if trace_id:
            return (f"\n{'='*60}\n"
                   f"[TRACE INFO] 接口: {endpoint}\n"
                   f"[TRACE INFO] TraceID: {trace_id}\n"
                   f"[TRACE INFO] 状态码: {response.status_code}\n"
                   f"[TRACE INFO] 响应时间: {response.elapsed.total_seconds():.3f}s\n"
                   f"{'='*60}")
        else:
            return (f"\n{'='*60}\n"
                   f"[TRACE INFO] 接口: {endpoint}\n"
                   f"[TRACE INFO] TraceID: 未找到\n"
                   f"[TRACE INFO] 状态码: {response.status_code}\n"
                   f"[TRACE INFO] 响应时间: {response.elapsed.total_seconds():.3f}s\n"
                   f"{'='*60}")
    
    @staticmethod
    def log_request_details(endpoint: str, data: Dict[str, Any], headers: Dict[str, str], response: requests.Response):
        """记录详细的请求信息"""
        trace_id = TraceUtils.extract_trace_id(response)
        
        print(f"\n{'='*60}")
        print(f"[REQUEST] 接口: {endpoint}")
        print(f"[REQUEST] 请求数据: {json.dumps(data, ensure_ascii=False, indent=2)}")
        print(f"[REQUEST] 请求头: {json.dumps(headers, ensure_ascii=False, indent=2)}")
        print(f"[RESPONSE] 状态码: {response.status_code}")
        print(f"[RESPONSE] TraceID: {trace_id or '未找到'}")
        print(f"[RESPONSE] 响应时间: {response.elapsed.total_seconds():.3f}s")
        
        # 如果响应不是很长，也打印响应内容
        try:
            response_text = response.text
            if len(response_text) < 1000:  # 只打印小于1000字符的响应
                response_data = response.json()
                print(f"[RESPONSE] 响应内容: {json.dumps(response_data, ensure_ascii=False, indent=2)}")
        except:
            pass
        print(f"{'='*60}")


# 全局测试数据管理器实例
test_data_manager = DataManager()