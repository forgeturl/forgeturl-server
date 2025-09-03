"""
测试工具类
提供测试过程中需要的通用功能
"""
import json
from typing import Dict, Any, Optional


class TestDataManager:
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


# 全局测试数据管理器实例
test_data_manager = TestDataManager()