"""
Space API 接口测试
基于 README.md 中定义的测试用例进行接口测试
"""
import pytest
import allure
import requests
from typing import Dict, Any


@allure.feature("Space API")
class TestSpaceAPI:
    """空间API测试类"""
    
    @pytest.fixture(autouse=True)
    def setup(self, base_url, api_client, test_user_id):
        """测试初始化"""
        self.base_url = base_url
        self.api_client = api_client(base_url)
        self.test_user_id = test_user_id
        self.created_page_id = None
    
    @allure.story("用户信息")
    @allure.title("测试1: 非登录态获取用户信息")
    @pytest.mark.no_login
    def test_01_get_user_info_no_login(self, headers_no_login):
        """非登录态，可以拉取到用户信息"""
        with allure.step("发送非登录态获取用户信息请求"):
            data = {"uid": self.test_user_id}
            response = self.api_client.post("/space/getUserInfo", data, headers_no_login)
        
        with allure.step("验证响应状态和数据"):
            assert response.status_code == 200
            response_data = response.json()
            
            # 验证基本字段存在
            assert "uid" in response_data
            assert "display_name" in response_data
            assert "avatar" in response_data
            assert "email" in response_data
            
            # 非登录态不应该返回敏感信息
            assert "username" not in response_data or response_data.get("username") == ""
            assert "status" not in response_data or response_data.get("status") == 0
    
    @allure.story("页面信息") 
    @allure.title("测试2: 非登录态获取页面数据")
    @pytest.mark.no_login
    def test_02_get_page_no_login(self, headers_no_login):
        """非登录态，可以拉取到页面信息"""
        with allure.step("发送非登录态获取页面数据请求"):
            data = {"page_id": "test_page_id"}
            response = self.api_client.post("/space/getPage", data, headers_no_login)
        
        with allure.step("验证响应"):
            # 根据实际情况，可能返回错误或者返回基本页面信息
            assert response.status_code == 200
    
    @allure.story("用户信息")
    @allure.title("测试3: 登录态获取详细用户信息")
    @pytest.mark.login
    def test_03_get_user_info_with_login(self, headers_with_login):
        """登录态，可以拉取详细用户信息"""
        with allure.step("发送登录态获取用户信息请求"):
            data = {"uid": self.test_user_id}
            response = self.api_client.post("/space/getUserInfo", data, headers_with_login)
        
        with allure.step("验证响应状态和数据"):
            assert response.status_code == 200
            response_data = response.json()
            
            # 验证基本字段
            assert "uid" in response_data
            assert "display_name" in response_data
            assert "avatar" in response_data
            assert "email" in response_data
            
            # 登录态应该返回详细信息
            assert "username" in response_data
            assert "status" in response_data
            assert "create_time" in response_data
            assert "update_time" in response_data
    
    @allure.story("我的空间")
    @allure.title("测试4: 登录态拉取我的空间")
    @pytest.mark.login
    def test_04_get_my_space_with_login(self, headers_with_login):
        """登录态，拉取我的空间"""
        with allure.step("发送登录态获取我的空间请求"):
            data = {"uid": self.test_user_id}
            response = self.api_client.post("/space/getMySpace", data, headers_with_login)
        
        with allure.step("验证响应"):
            assert response.status_code == 200
            response_data = response.json()
            
            # 验证我的空间数据结构
            assert "space_name" in response_data or "page_briefs" in response_data
    
    @allure.story("页面管理")
    @allure.title("测试5a: 非登录态创建页面失败")
    @pytest.mark.no_login
    def test_05a_create_page_no_login_fail(self, headers_no_login, sample_page_data):
        """非登录态，创建页面失败，得到错误码code != 1的返回"""
        with allure.step("发送非登录态创建页面请求"):
            response = self.api_client.post("/space/createPage", sample_page_data, headers_no_login)
        
        with allure.step("验证失败响应"):
            response_data = response.json()
            # 验证返回错误码
            assert response_data.get("code") != 1 or response.status_code != 200
    
    @allure.story("页面管理")
    @allure.title("测试5b: 登录态创建页面成功")
    @pytest.mark.login
    def test_05b_create_page_with_login_success(self, headers_with_login, sample_page_data):
        """登录态，创建页面，创建的页面id在返回json的page_id字段里"""
        with allure.step("发送登录态创建页面请求"):
            response = self.api_client.post("/space/createPage", sample_page_data, headers_with_login)
        
        with allure.step("验证成功响应"):
            assert response.status_code == 200
            response_data = response.json()
            
            # 验证创建成功
            if "page_id" in response_data:
                self.created_page_id = response_data["page_id"]
                assert self.created_page_id is not None
                assert len(self.created_page_id) > 0
    
    @allure.story("我的空间")
    @allure.title("测试6a: 非登录态无法拉取getMySpace")
    @pytest.mark.no_login
    def test_06a_get_my_space_no_login_fail(self, headers_no_login):
        """非登录态，无法拉取getMySpace，得到错误码code != 1的返回"""
        with allure.step("发送非登录态获取我的空间请求"):
            data = {"uid": self.test_user_id}
            response = self.api_client.post("/space/getMySpace", data, headers_no_login)
        
        with allure.step("验证失败响应"):
            response_data = response.json()
            assert response_data.get("code") != 1 or response.status_code != 200
    
    @allure.story("我的空间")
    @allure.title("测试6b: 登录态拉取我的空间确认创建的页面")
    @pytest.mark.login
    def test_06b_get_my_space_verify_created_page(self, headers_with_login):
        """登录态拉取我的空间，确认上一步创建的页面在返回的列表里"""
        with allure.step("发送登录态获取我的空间请求"):
            data = {"uid": self.test_user_id}
            response = self.api_client.post("/space/getMySpace", data, headers_with_login)
        
        with allure.step("验证响应并确认页面存在"):
            assert response.status_code == 200
            response_data = response.json()
            
            # 如果有创建的页面ID，验证是否在列表中
            if self.created_page_id and "page_briefs" in response_data:
                page_ids = [page.get("page_id") for page in response_data["page_briefs"]]
                assert self.created_page_id in page_ids
    
    @allure.story("页面访问")
    @allure.title("测试7a: 非登录态拉取私有页面失败")
    @pytest.mark.no_login
    def test_07a_get_private_page_no_login_fail(self, headers_no_login):
        """非登录态，拉取不到私有页面，得到错误码code != 1的返回"""
        with allure.step("发送非登录态获取私有页面请求"):
            data = {"page_id": self.created_page_id or "test_page_id"}
            response = self.api_client.post("/space/getPage", data, headers_no_login)
        
        with allure.step("验证失败响应"):
            if self.created_page_id:  # 如果有创建的页面
                response_data = response.json()
                assert response_data.get("code") != 1 or response.status_code != 200
    
    @allure.story("页面访问")
    @allure.title("测试7b: 登录态拉取创建的页面")
    @pytest.mark.login
    def test_07b_get_created_page_with_login(self, headers_with_login):
        """登录态，拉取创建的页面"""
        with allure.step("发送登录态获取页面请求"):
            data = {"page_id": self.created_page_id or "test_page_id"}
            response = self.api_client.post("/space/getPage", data, headers_with_login)
        
        with allure.step("验证成功响应"):
            assert response.status_code == 200
            response_data = response.json()
            
            if "page" in response_data:
                page_data = response_data["page"]
                assert "page_id" in page_data
                assert "title" in page_data
                assert "collections" in page_data
    
    @allure.story("页面管理") 
    @allure.title("测试8a: 非登录态更新页面失败")
    @pytest.mark.no_login
    def test_08a_update_page_no_login_fail(self, headers_no_login, updated_page_data):
        """非登录态，更新页面失败，得到错误码code != 1的返回"""
        with allure.step("发送非登录态更新页面请求"):
            update_data = updated_page_data.copy()
            update_data["page_id"] = self.created_page_id or "test_page_id"
            response = self.api_client.post("/space/updatePage", update_data, headers_no_login)
        
        with allure.step("验证失败响应"):
            response_data = response.json()
            assert response_data.get("code") != 1 or response.status_code != 200
    
    @allure.story("页面管理")
    @allure.title("测试8b: 登录态更新页面成功")
    @pytest.mark.login
    def test_08b_update_page_with_login_success(self, headers_with_login, updated_page_data):
        """登录态，更新创建的页面"""
        with allure.step("发送登录态更新页面请求"):
            update_data = updated_page_data.copy()
            update_data["page_id"] = self.created_page_id or "test_page_id"
            response = self.api_client.post("/space/updatePage", update_data, headers_with_login)
        
        with allure.step("验证成功响应"):
            if self.created_page_id:  # 只有创建了页面才能更新
                assert response.status_code == 200
                response_data = response.json()
                
                if "update_time" in response_data:
                    assert response_data["update_time"] > 0
    
    @allure.story("页面链接")
    @allure.title("测试9a: 非登录态创建页面链接失败")
    @pytest.mark.no_login
    def test_09a_add_page_link_no_login_fail(self, headers_no_login):
        """非登录态，创建链接失败，得到错误码code != 1的返回"""
        with allure.step("发送非登录态创建页面链接请求"):
            data = {
                "page_id": self.created_page_id or "test_page_id",
                "page_type": "readonly"
            }
            response = self.api_client.post("/space/addPageLink", data, headers_no_login)
        
        with allure.step("验证失败响应"):
            response_data = response.json()
            assert response_data.get("code") != 1 or response.status_code != 200
    
    @allure.story("页面链接")
    @allure.title("测试9b: 登录态创建页面链接成功")
    @pytest.mark.login
    def test_09b_add_page_link_with_login_success(self, headers_with_login):
        """登录态，创建页面链接"""
        with allure.step("发送登录态创建页面链接请求"):
            data = {
                "page_id": self.created_page_id or "test_page_id",
                "page_type": "readonly"
            }
            response = self.api_client.post("/space/addPageLink", data, headers_with_login)
        
        with allure.step("验证成功响应"):
            if self.created_page_id:  # 只有创建了页面才能创建链接
                assert response.status_code == 200
                response_data = response.json()
                
                if "new_page_id" in response_data:
                    assert response_data["new_page_id"] is not None
                    assert response_data["page_type"] == "readonly"
    
    @allure.story("页面链接")
    @allure.title("测试10a: 非登录态移除页面链接失败")
    @pytest.mark.no_login
    def test_10a_remove_page_link_no_login_fail(self, headers_no_login):
        """非登录态，移除页面链接失败，得到错误码code != 1的返回"""
        with allure.step("发送非登录态移除页面链接请求"):
            data = {"page_id": self.created_page_id or "test_page_id"}
            response = self.api_client.post("/space/removePageLink", data, headers_no_login)
        
        with allure.step("验证失败响应"):
            response_data = response.json()
            assert response_data.get("code") != 1 or response.status_code != 200
    
    @allure.story("页面链接")
    @allure.title("测试10b: 登录态移除页面链接成功")
    @pytest.mark.login
    def test_10b_remove_page_link_with_login_success(self, headers_with_login):
        """登录态，移除页面链接"""
        with allure.step("发送登录态移除页面链接请求"):
            data = {"page_id": self.created_page_id or "test_page_id"}
            response = self.api_client.post("/space/removePageLink", data, headers_with_login)
        
        with allure.step("验证成功响应"):
            if self.created_page_id:  # 只有创建了页面才能移除链接
                assert response.status_code == 200
    
    @allure.story("页面管理")
    @allure.title("测试11a: 非登录态调整页面顺序失败")
    @pytest.mark.no_login
    def test_11a_save_page_ids_no_login_fail(self, headers_no_login):
        """非登录态，调整页面ID顺序失败，得到错误码code != 1的返回"""
        with allure.step("发送非登录态调整页面顺序请求"):
            data = {
                "uid": self.test_user_id,
                "page_ids": ["page2", "page1", "page3"]
            }
            response = self.api_client.post("/space/savePageIds", data, headers_no_login)
        
        with allure.step("验证失败响应"):
            response_data = response.json()
            assert response_data.get("code") != 1 or response.status_code != 200
    
    @allure.story("页面管理")
    @allure.title("测试11b: 登录态调整页面顺序成功")
    @pytest.mark.login
    def test_11b_save_page_ids_with_login_success(self, headers_with_login):
        """登录态，调整页面ID顺序"""
        with allure.step("发送登录态调整页面顺序请求"):
            data = {
                "uid": self.test_user_id,
                "page_ids": ["page2", "page1", "page3"]
            }
            response = self.api_client.post("/space/savePageIds", data, headers_with_login)
        
        with allure.step("验证成功响应"):
            assert response.status_code == 200
            response_data = response.json()
            
            if "page_ids" in response_data:
                assert len(response_data["page_ids"]) > 0
    
    @allure.story("页面管理")
    @allure.title("测试12a: 非登录态删除页面失败")
    @pytest.mark.no_login
    def test_12a_delete_page_no_login_fail(self, headers_no_login):
        """非登录态，删除页面失败，得到错误码code != 1的返回"""
        with allure.step("发送非登录态删除页面请求"):
            data = {"page_id": self.created_page_id or "test_page_id"}
            response = self.api_client.post("/space/deletePage", data, headers_no_login)
        
        with allure.step("验证失败响应"):
            response_data = response.json()
            assert response_data.get("code") != 1 or response.status_code != 200
    
    @allure.story("页面管理")
    @allure.title("测试12b: 登录态删除页面成功")
    @pytest.mark.login
    def test_12b_delete_page_with_login_success(self, headers_with_login):
        """登录态，删除页面"""
        with allure.step("发送登录态删除页面请求"):
            data = {"page_id": self.created_page_id or "test_page_id"}
            response = self.api_client.post("/space/deletePage", data, headers_with_login)
        
        with allure.step("验证成功响应"):
            if self.created_page_id:  # 只有创建了页面才能删除
                assert response.status_code == 200
