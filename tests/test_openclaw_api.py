"""
OpenClaw API 接口测试
测试 OpenClaw 服务相关的所有接口

接口列表:
- POST /openclaw/getApiKey          获取API Key
- POST /openclaw/regenerateApiKey   重新生成API Key
- POST /openclaw/tmpBookmark/add    添加临时书签
- POST /openclaw/tmpBookmark/list   列出临时书签
- POST /openclaw/tmpBookmark/exists 检查临时书签是否存在
- POST /openclaw/tmpBookmark/moveToPage   将临时书签移动到页面
- POST /openclaw/tmpBookmark/moveFromPage 将页面书签移动到临时书签

认证方式:
- X-Token: 登录态token (header: X-Token)
- Bearer API Key: OpenClaw 专用认证 (header: Authorization: Bearer <api_key>)
"""
import pytest
import allure
import time
import requests
from typing import Dict, Any


@allure.feature("OpenClaw API")
class TestOpenClawAPI:
    """OpenClaw API 测试类"""

    @pytest.fixture(autouse=True)
    def setup(self, base_url, api_client, test_user_id):
        """测试初始化"""
        self.base_url = base_url
        self.api_client = api_client(base_url)
        self.test_user_id = test_user_id
        self.api_key = None  # 用于存储生成的 API Key

    # ========== 辅助方法 ==========

    def _headers_with_bearer(self, api_key: str) -> Dict[str, str]:
        """构建 Bearer token 请求头"""
        return {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {api_key}"
        }

    # ========== GetApiKey 接口测试 ==========

    @allure.story("API Key 管理")
    @allure.title("测试1a: 非登录态获取API Key失败")
    @pytest.mark.no_login
    @pytest.mark.openclaw
    def test_01a_get_api_key_no_login(self, headers_no_login):
        """非登录态，获取API Key失败，返回错误码"""
        with allure.step("发送非登录态获取API Key请求"):
            response = self.api_client.post("/openclaw/getApiKey", {}, headers_no_login)

        with allure.step("验证失败响应"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") != 1, f"非登录态应该返回错误: {response_data}"

    @allure.story("API Key 管理")
    @allure.title("测试1b: 登录态获取API Key")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_01b_get_api_key_with_login(self, headers_with_login):
        """登录态，获取API Key，首次应该没有Key"""
        with allure.step("发送登录态获取API Key请求"):
            response = self.api_client.post_with_detailed_log(
                "/openclaw/getApiKey", {}, headers_with_login
            )

        with allure.step("验证响应"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") == 1, f"登录态获取API Key应该成功: {response_data}"

            data = response_data.get("data", {})
            # has_key 字段应该存在
            assert "has_key" in data, "响应中应该包含 has_key 字段"

    # ========== RegenerateApiKey 接口测试 ==========

    @allure.story("API Key 管理")
    @allure.title("测试2a: 非登录态生成API Key失败")
    @pytest.mark.no_login
    @pytest.mark.openclaw
    def test_02a_regenerate_api_key_no_login(self, headers_no_login):
        """非登录态，生成API Key失败"""
        with allure.step("发送非登录态生成API Key请求"):
            response = self.api_client.post("/openclaw/regenerateApiKey", {}, headers_no_login)

        with allure.step("验证失败响应"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") != 1, f"非登录态应该返回错误: {response_data}"

    @allure.story("API Key 管理")
    @allure.title("测试2b: 登录态生成API Key成功")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_02b_regenerate_api_key_with_login(self, headers_with_login):
        """登录态，生成API Key成功"""
        with allure.step("发送登录态生成API Key请求"):
            response = self.api_client.post_with_detailed_log(
                "/openclaw/regenerateApiKey", {}, headers_with_login
            )

        with allure.step("验证成功响应"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") == 1, f"登录态生成API Key应该成功: {response_data}"

            data = response_data.get("data", {})
            assert "api_key" in data, "响应中应该包含 api_key 字段"
            assert data["api_key"], "api_key 不应为空"
            self.api_key = data["api_key"]

    @allure.story("API Key 管理")
    @allure.title("测试2c: 再次生成后获取API Key确认has_key为true")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_02c_get_api_key_after_regenerate(self, headers_with_login):
        """登录态生成API Key后，获取API Key应该has_key=true"""
        with allure.step("先生成API Key"):
            regen_resp = self.api_client.post("/openclaw/regenerateApiKey", {}, headers_with_login)
            assert regen_resp.status_code == 200
            regen_data = regen_resp.json()
            assert regen_data.get("code") == 1

        with allure.step("获取API Key"):
            response = self.api_client.post("/openclaw/getApiKey", {}, headers_with_login)

        with allure.step("验证has_key为true"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") == 1
            data = response_data.get("data", {})
            assert data.get("has_key") is True, f"生成后has_key应该为true: {data}"
            assert data.get("api_key"), "api_key 不应为空"

    # ========== AddTmpBookmark 接口测试 ==========

    @allure.story("临时书签管理")
    @allure.title("测试3a: 非登录态添加临时书签失败")
    @pytest.mark.no_login
    @pytest.mark.openclaw
    def test_03a_add_tmp_bookmark_no_login(self, headers_no_login):
        """非登录态，添加临时书签失败"""
        with allure.step("发送非登录态添加临时书签请求"):
            data = {
                "title": "测试书签",
                "url": "https://example.com/test"
            }
            response = self.api_client.post("/openclaw/tmpBookmark/add", data, headers_no_login)

        with allure.step("验证失败响应"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") != 1, f"非登录态应该返回错误: {response_data}"

    @allure.story("临时书签管理")
    @allure.title("测试3b: 登录态添加临时书签成功")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_03b_add_tmp_bookmark_with_login(self, headers_with_login):
        """登录态，添加临时书签成功"""
        with allure.step("发送登录态添加临时书签请求"):
            data = {
                "title": "测试书签-登录态",
                "url": "https://example.com/test-login"
            }
            response = self.api_client.post_with_detailed_log(
                "/openclaw/tmpBookmark/add", data, headers_with_login
            )

        with allure.step("验证成功响应"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") == 1, f"登录态添加书签应该成功: {response_data}"

            data = response_data.get("data", {})
            bookmark = data.get("bookmark", {})
            assert bookmark.get("id"), "书签ID不应为空"
            assert bookmark.get("title") == "测试书签-登录态"
            assert bookmark.get("url") == "https://example.com/test-login"

    @allure.story("临时书签管理")
    @allure.title("测试3c: Bearer认证添加临时书签成功")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_03c_add_tmp_bookmark_with_bearer(self, headers_with_login):
        """通过Bearer API Key认证，添加临时书签成功"""
        with allure.step("先生成API Key"):
            regen_resp = self.api_client.post("/openclaw/regenerateApiKey", {}, headers_with_login)
            assert regen_resp.status_code == 200
            regen_data = regen_resp.json()
            assert regen_data.get("code") == 1
            api_key = regen_data["data"]["api_key"]

        with allure.step("使用Bearer认证添加临时书签"):
            bookmark_data = {
                "title": "通过API Key添加的书签",
                "url": "https://example.com/api-key-test"
            }
            bearer_headers = self._headers_with_bearer(api_key)
            response = self.api_client.post_with_detailed_log(
                "/openclaw/tmpBookmark/add", bookmark_data, bearer_headers
            )

        with allure.step("验证成功响应"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") == 1, f"Bearer认证添加书签应该成功: {response_data}"

            data = response_data.get("data", {})
            bookmark = data.get("bookmark", {})
            assert bookmark.get("id"), "书签ID不应为空"

    @allure.story("临时书签管理")
    @allure.title("测试3d: 无效Bearer Token添加临时书签失败")
    @pytest.mark.no_login
    @pytest.mark.openclaw
    def test_03d_add_tmp_bookmark_invalid_bearer(self):
        """使用无效的Bearer Token，添加临时书签失败"""
        with allure.step("使用无效的Bearer Token发送请求"):
            bookmark_data = {
                "title": "无效Token测试",
                "url": "https://example.com/invalid-token"
            }
            invalid_headers = self._headers_with_bearer("invalid_api_key_12345")
            response = self.api_client.post(
                "/openclaw/tmpBookmark/add", bookmark_data, invalid_headers
            )

        with allure.step("验证失败响应"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") != 1, f"无效Token应该返回错误: {response_data}"

    @allure.story("临时书签管理")
    @allure.title("测试3e: 添加重复书签应该去重(upsert)")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_03e_add_duplicate_tmp_bookmark(self, headers_with_login):
        """添加相同title和url的书签，应该执行upsert去重"""
        with allure.step("第一次添加书签"):
            data = {
                "title": "重复测试书签",
                "url": "https://example.com/duplicate-test"
            }
            resp1 = self.api_client.post("/openclaw/tmpBookmark/add", data, headers_with_login)
            assert resp1.status_code == 200
            resp1_data = resp1.json()
            assert resp1_data.get("code") == 1
            first_id = resp1_data["data"]["bookmark"]["id"]

        with allure.step("第二次添加相同书签"):
            resp2 = self.api_client.post("/openclaw/tmpBookmark/add", data, headers_with_login)
            assert resp2.status_code == 200
            resp2_data = resp2.json()
            assert resp2_data.get("code") == 1
            second_id = resp2_data["data"]["bookmark"]["id"]

        with allure.step("验证两次返回的ID相同(upsert)"):
            assert first_id == second_id, f"重复添加应该返回相同ID: first={first_id}, second={second_id}"

    # ========== ListTmpBookmarks 接口测试 ==========

    @allure.story("临时书签管理")
    @allure.title("测试4a: 非登录态列出临时书签失败")
    @pytest.mark.no_login
    @pytest.mark.openclaw
    def test_04a_list_tmp_bookmarks_no_login(self, headers_no_login):
        """非登录态，列出临时书签失败"""
        with allure.step("发送非登录态列出书签请求"):
            response = self.api_client.post("/openclaw/tmpBookmark/list", {}, headers_no_login)

        with allure.step("验证失败响应"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") != 1, f"非登录态应该返回错误: {response_data}"

    @allure.story("临时书签管理")
    @allure.title("测试4b: 登录态列出临时书签成功")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_04b_list_tmp_bookmarks_with_login(self, headers_with_login):
        """登录态，列出临时书签"""
        with allure.step("先添加一个书签确保列表不为空"):
            add_data = {
                "title": "列表测试书签",
                "url": "https://example.com/list-test"
            }
            add_resp = self.api_client.post("/openclaw/tmpBookmark/add", add_data, headers_with_login)
            assert add_resp.status_code == 200

        with allure.step("发送登录态列出书签请求"):
            response = self.api_client.post_with_detailed_log(
                "/openclaw/tmpBookmark/list", {}, headers_with_login
            )

        with allure.step("验证成功响应"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") == 1, f"登录态列出书签应该成功: {response_data}"

            data = response_data.get("data", {})
            assert "bookmarks" in data, "响应中应该包含 bookmarks 字段"
            bookmarks = data["bookmarks"]
            assert isinstance(bookmarks, list), "bookmarks 应该是列表"

    @allure.story("临时书签管理")
    @allure.title("测试4c: 列出临时书签带limit参数")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_04c_list_tmp_bookmarks_with_limit(self, headers_with_login):
        """登录态，使用limit参数列出临时书签"""
        with allure.step("添加多个书签"):
            for i in range(3):
                add_data = {
                    "title": f"限制测试书签_{i}_{int(time.time())}",
                    "url": f"https://example.com/limit-test-{i}-{int(time.time())}"
                }
                self.api_client.post("/openclaw/tmpBookmark/add", add_data, headers_with_login)

        with allure.step("使用limit参数列出书签"):
            response = self.api_client.post(
                "/openclaw/tmpBookmark/list", {"limit": 2}, headers_with_login
            )

        with allure.step("验证响应"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") == 1
            data = response_data.get("data", {})
            bookmarks = data.get("bookmarks", [])
            assert len(bookmarks) <= 2, f"limit=2时书签数量不应超过2: 实际={len(bookmarks)}"

    # ========== ExistsTmpBookmark 接口测试 ==========

    @allure.story("临时书签管理")
    @allure.title("测试5a: 非登录态检查书签是否存在失败")
    @pytest.mark.no_login
    @pytest.mark.openclaw
    def test_05a_exists_tmp_bookmark_no_login(self, headers_no_login):
        """非登录态，检查临时书签是否存在失败"""
        with allure.step("发送非登录态检查书签请求"):
            data = {
                "title": "测试书签",
                "url": "https://example.com/test"
            }
            response = self.api_client.post("/openclaw/tmpBookmark/exists", data, headers_no_login)

        with allure.step("验证失败响应"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") != 1, f"非登录态应该返回错误: {response_data}"

    @allure.story("临时书签管理")
    @allure.title("测试5b: 登录态检查已存在的书签")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_05b_exists_tmp_bookmark_found(self, headers_with_login):
        """登录态，检查已添加的临时书签"""
        unique_suffix = str(int(time.time()))
        bookmark_title = f"存在性测试书签_{unique_suffix}"
        bookmark_url = f"https://example.com/exists-test-{unique_suffix}"

        with allure.step("先添加一个书签"):
            add_data = {
                "title": bookmark_title,
                "url": bookmark_url
            }
            add_resp = self.api_client.post("/openclaw/tmpBookmark/add", add_data, headers_with_login)
            assert add_resp.status_code == 200
            assert add_resp.json().get("code") == 1

        with allure.step("检查书签是否存在"):
            response = self.api_client.post_with_detailed_log(
                "/openclaw/tmpBookmark/exists", add_data, headers_with_login
            )

        with allure.step("验证书签存在"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") == 1, f"检查书签应该成功: {response_data}"

            data = response_data.get("data", {})
            assert data.get("exists") is True, f"书签应该存在: {data}"
            assert data.get("bookmark"), "存在时应该返回书签信息"

    @allure.story("临时书签管理")
    @allure.title("测试5c: 登录态检查不存在的书签")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_05c_exists_tmp_bookmark_not_found(self, headers_with_login):
        """登录态，检查不存在的临时书签"""
        with allure.step("检查一个不存在的书签"):
            data = {
                "title": "这个书签肯定不存在",
                "url": f"https://example.com/not-exists-{int(time.time())}"
            }
            response = self.api_client.post(
                "/openclaw/tmpBookmark/exists", data, headers_with_login
            )

        with allure.step("验证书签不存在"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") == 1, f"检查书签应该成功: {response_data}"

            data = response_data.get("data", {})
            assert data.get("exists") is False, f"书签不应该存在: {data}"

    # ========== MoveTmpBookmarkToPage 接口测试 ==========

    @allure.story("书签移动")
    @allure.title("测试6a: 非登录态移动临时书签到页面失败")
    @pytest.mark.no_login
    @pytest.mark.openclaw
    def test_06a_move_tmp_bookmark_to_page_no_login(self, headers_no_login):
        """非登录态，移动临时书签到页面失败"""
        with allure.step("发送非登录态移动书签请求"):
            data = {
                "bookmark_id": 1,
                "page_id": "test_page_id"
            }
            response = self.api_client.post(
                "/openclaw/tmpBookmark/moveToPage", data, headers_no_login
            )

        with allure.step("验证失败响应"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") != 1, f"非登录态应该返回错误: {response_data}"

    @allure.story("书签移动")
    @allure.title("测试6b: 登录态移动临时书签到页面")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_06b_move_tmp_bookmark_to_page_with_login(self, headers_with_login):
        """登录态，移动临时书签到页面（需先创建页面和书签）"""
        with allure.step("先创建一个页面"):
            page_data = {
                "title": "OpenClaw移动测试页面",
                "brief": "用于测试书签移动功能",
                "collections": [
                    {
                        "links": [
                            {
                                "title": "已有链接",
                                "url": "https://example.com/existing"
                            }
                        ]
                    }
                ]
            }
            create_resp = self.api_client.post("/space/createPage", page_data, headers_with_login)
            assert create_resp.status_code == 200
            create_data = create_resp.json()

            # 获取页面ID
            page_id = None
            if create_data.get("code") == 1:
                page_id = create_data.get("data", {}).get("page_id")

        with allure.step("添加一个临时书签"):
            unique_suffix = str(int(time.time()))
            bookmark_data = {
                "title": f"待移动书签_{unique_suffix}",
                "url": f"https://example.com/move-test-{unique_suffix}"
            }
            add_resp = self.api_client.post(
                "/openclaw/tmpBookmark/add", bookmark_data, headers_with_login
            )
            assert add_resp.status_code == 200
            add_data = add_resp.json()
            bookmark_id = None
            if add_data.get("code") == 1:
                bookmark_id = add_data["data"]["bookmark"]["id"]

        if page_id and bookmark_id:
            with allure.step("将临时书签移动到页面"):
                move_data = {
                    "bookmark_id": bookmark_id,
                    "page_id": page_id,
                    "collection_title": ""
                }
                response = self.api_client.post_with_detailed_log(
                    "/openclaw/tmpBookmark/moveToPage", move_data, headers_with_login
                )

            with allure.step("验证移动成功"):
                assert response.status_code == 200
                response_data = response.json()
                assert response_data.get("code") == 1, f"移动书签应该成功: {response_data}"

            with allure.step("验证书签已从临时列表中移除"):
                exists_data = {
                    "title": bookmark_data["title"],
                    "url": bookmark_data["url"]
                }
                exists_resp = self.api_client.post(
                    "/openclaw/tmpBookmark/exists", exists_data, headers_with_login
                )
                if exists_resp.status_code == 200:
                    exists_result = exists_resp.json()
                    if exists_result.get("code") == 1:
                        assert exists_result["data"].get("exists") is False, \
                            "移动后书签不应存在于临时列表中"

            with allure.step("清理: 删除测试页面"):
                self.api_client.post(
                    "/space/deletePage", {"page_id": page_id}, headers_with_login
                )

    @allure.story("书签移动")
    @allure.title("测试6c: 移动不存在的书签到页面失败")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_06c_move_nonexistent_bookmark_to_page(self, headers_with_login):
        """登录态，移动不存在的书签ID到页面应该失败"""
        with allure.step("发送移动不存在书签请求"):
            data = {
                "bookmark_id": 999999999,
                "page_id": "test_page_id"
            }
            response = self.api_client.post(
                "/openclaw/tmpBookmark/moveToPage", data, headers_with_login
            )

        with allure.step("验证失败响应"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") != 1, f"移动不存在的书签应该失败: {response_data}"

    # ========== MovePageBookmarkToTmp 接口测试 ==========

    @allure.story("书签移动")
    @allure.title("测试7a: 非登录态从页面移动书签到临时书签失败")
    @pytest.mark.no_login
    @pytest.mark.openclaw
    def test_07a_move_page_bookmark_to_tmp_no_login(self, headers_no_login):
        """非登录态，从页面移动书签到临时书签失败"""
        with allure.step("发送非登录态移动请求"):
            data = {
                "page_id": "test_page_id",
                "title": "测试链接",
                "url": "https://example.com"
            }
            response = self.api_client.post(
                "/openclaw/tmpBookmark/moveFromPage", data, headers_no_login
            )

        with allure.step("验证失败响应"):
            assert response.status_code == 200
            response_data = response.json()
            assert response_data.get("code") != 1, f"非登录态应该返回错误: {response_data}"

    @allure.story("书签移动")
    @allure.title("测试7b: 登录态从页面移动书签到临时书签")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_07b_move_page_bookmark_to_tmp_with_login(self, headers_with_login):
        """登录态，将页面中的链接移动到临时书签"""
        unique_suffix = str(int(time.time()))
        link_title = f"待移出链接_{unique_suffix}"
        link_url = f"https://example.com/move-from-page-{unique_suffix}"
        collection_title = "测试收藏夹"

        with allure.step("先创建一个包含链接的页面"):
            page_data = {
                "title": "移出链接测试页面",
                "brief": "用于测试从页面移动书签到临时书签",
                "collections": [
                    {
                        "title": collection_title,
                        "links": [
                            {
                                "title": link_title,
                                "url": link_url
                            }
                        ]
                    }
                ]
            }
            create_resp = self.api_client.post("/space/createPage", page_data, headers_with_login)
            assert create_resp.status_code == 200
            create_data = create_resp.json()

            page_id = None
            if create_data.get("code") == 1:
                page_id = create_data.get("data", {}).get("page_id")

        if page_id:
            with allure.step("将页面中的链接移动到临时书签"):
                move_data = {
                    "page_id": page_id,
                    "collection_title": collection_title,
                    "title": link_title,
                    "url": link_url
                }
                response = self.api_client.post_with_detailed_log(
                    "/openclaw/tmpBookmark/moveFromPage", move_data, headers_with_login
                )

            with allure.step("验证移动成功"):
                assert response.status_code == 200
                response_data = response.json()
                assert response_data.get("code") == 1, f"从页面移动书签应该成功: {response_data}"

            with allure.step("验证书签已添加到临时列表"):
                exists_data = {
                    "title": link_title,
                    "url": link_url
                }
                exists_resp = self.api_client.post(
                    "/openclaw/tmpBookmark/exists", exists_data, headers_with_login
                )
                if exists_resp.status_code == 200:
                    exists_result = exists_resp.json()
                    if exists_result.get("code") == 1:
                        assert exists_result["data"].get("exists") is True, \
                            "移动后书签应该存在于临时列表中"

            with allure.step("清理: 删除测试页面"):
                self.api_client.post(
                    "/space/deletePage", {"page_id": page_id}, headers_with_login
                )

    @allure.story("书签移动")
    @allure.title("测试7c: 在页面中移动不存在的链接到临时书签失败")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_07c_move_nonexistent_page_bookmark_to_tmp(self, headers_with_login):
        """登录态，从页面中移动不存在的链接应该失败"""
        with allure.step("创建一个空页面"):
            page_data = {
                "title": "空页面测试",
                "brief": "空的测试页面",
                "collections": []
            }
            create_resp = self.api_client.post("/space/createPage", page_data, headers_with_login)
            assert create_resp.status_code == 200
            create_data = create_resp.json()

            page_id = None
            if create_data.get("code") == 1:
                page_id = create_data.get("data", {}).get("page_id")

        if page_id:
            with allure.step("尝试移动不存在的链接"):
                move_data = {
                    "page_id": page_id,
                    "title": "不存在的链接",
                    "url": "https://example.com/not-exist"
                }
                response = self.api_client.post(
                    "/openclaw/tmpBookmark/moveFromPage", move_data, headers_with_login
                )

            with allure.step("验证失败响应"):
                assert response.status_code == 200
                response_data = response.json()
                assert response_data.get("code") != 1, f"移动不存在的链接应该失败: {response_data}"

            with allure.step("清理: 删除测试页面"):
                self.api_client.post(
                    "/space/deletePage", {"page_id": page_id}, headers_with_login
                )

    # ========== 端到端流程测试 ==========

    @allure.story("端到端流程")
    @allure.title("测试8: 完整的书签管理流程")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_08_e2e_bookmark_workflow(self, headers_with_login):
        """端到端测试: 生成API Key -> 添加书签 -> 列出 -> 检查存在 -> 移动到页面"""
        unique_suffix = str(int(time.time()))

        with allure.step("步骤1: 生成API Key"):
            regen_resp = self.api_client.post(
                "/openclaw/regenerateApiKey", {}, headers_with_login
            )
            assert regen_resp.status_code == 200
            regen_data = regen_resp.json()
            assert regen_data.get("code") == 1
            api_key = regen_data["data"]["api_key"]
            assert api_key, "API Key 不应为空"

        with allure.step("步骤2: 使用API Key添加临时书签"):
            bearer_headers = self._headers_with_bearer(api_key)
            bookmark_data = {
                "title": f"E2E测试书签_{unique_suffix}",
                "url": f"https://example.com/e2e-test-{unique_suffix}"
            }
            add_resp = self.api_client.post(
                "/openclaw/tmpBookmark/add", bookmark_data, bearer_headers
            )
            assert add_resp.status_code == 200
            add_data = add_resp.json()
            assert add_data.get("code") == 1, f"添加书签应该成功: {add_data}"
            bookmark_id = add_data["data"]["bookmark"]["id"]

        with allure.step("步骤3: 列出书签确认已添加"):
            list_resp = self.api_client.post(
                "/openclaw/tmpBookmark/list", {}, bearer_headers
            )
            assert list_resp.status_code == 200
            list_data = list_resp.json()
            assert list_data.get("code") == 1
            bookmarks = list_data["data"].get("bookmarks", [])
            bookmark_ids = [b["id"] for b in bookmarks]
            assert bookmark_id in bookmark_ids, f"新添加的书签应该在列表中: id={bookmark_id}"

        with allure.step("步骤4: 检查书签是否存在"):
            exists_resp = self.api_client.post(
                "/openclaw/tmpBookmark/exists", bookmark_data, bearer_headers
            )
            assert exists_resp.status_code == 200
            exists_data = exists_resp.json()
            assert exists_data.get("code") == 1
            assert exists_data["data"].get("exists") is True, "书签应该存在"

        with allure.step("步骤5: 创建一个目标页面"):
            page_data = {
                "title": f"E2E目标页面_{unique_suffix}",
                "brief": "端到端测试目标页面",
                "collections": []
            }
            create_resp = self.api_client.post(
                "/space/createPage", page_data, headers_with_login
            )
            assert create_resp.status_code == 200
            create_data = create_resp.json()
            page_id = None
            if create_data.get("code") == 1:
                page_id = create_data.get("data", {}).get("page_id")

        if page_id:
            with allure.step("步骤6: 将书签移动到页面"):
                move_data = {
                    "bookmark_id": bookmark_id,
                    "page_id": page_id,
                    "collection_title": "E2E测试收藏夹"
                }
                move_resp = self.api_client.post(
                    "/openclaw/tmpBookmark/moveToPage", move_data, headers_with_login
                )
                assert move_resp.status_code == 200
                move_result = move_resp.json()
                assert move_result.get("code") == 1, f"移动书签应该成功: {move_result}"

            with allure.step("步骤7: 验证书签已从临时列表中移除"):
                exists_resp2 = self.api_client.post(
                    "/openclaw/tmpBookmark/exists", bookmark_data, headers_with_login
                )
                if exists_resp2.status_code == 200:
                    exists_data2 = exists_resp2.json()
                    if exists_data2.get("code") == 1:
                        assert exists_data2["data"].get("exists") is False, \
                            "移动后书签不应在临时列表中"

            with allure.step("清理: 删除测试页面"):
                self.api_client.post(
                    "/space/deletePage", {"page_id": page_id}, headers_with_login
                )

    @allure.story("端到端流程")
    @allure.title("测试9: API Key重新生成后旧Key失效")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_09_regenerate_invalidates_old_key(self, headers_with_login):
        """重新生成API Key后，旧的Key应该失效"""
        with allure.step("步骤1: 生成第一个API Key"):
            regen1_resp = self.api_client.post(
                "/openclaw/regenerateApiKey", {}, headers_with_login
            )
            assert regen1_resp.status_code == 200
            old_key = regen1_resp.json()["data"]["api_key"]

        with allure.step("步骤2: 验证旧Key可用"):
            old_headers = self._headers_with_bearer(old_key)
            list_resp1 = self.api_client.post(
                "/openclaw/tmpBookmark/list", {}, old_headers
            )
            assert list_resp1.status_code == 200
            assert list_resp1.json().get("code") == 1, "旧Key应该可用"

        with allure.step("步骤3: 生成新API Key"):
            regen2_resp = self.api_client.post(
                "/openclaw/regenerateApiKey", {}, headers_with_login
            )
            assert regen2_resp.status_code == 200
            new_key = regen2_resp.json()["data"]["api_key"]
            assert new_key != old_key, "新旧Key应该不同"

        with allure.step("步骤4: 验证旧Key已失效"):
            list_resp2 = self.api_client.post(
                "/openclaw/tmpBookmark/list", {}, old_headers
            )
            assert list_resp2.status_code == 200
            resp_data = list_resp2.json()
            assert resp_data.get("code") != 1, f"旧Key应该失效: {resp_data}"

        with allure.step("步骤5: 验证新Key可用"):
            new_headers = self._headers_with_bearer(new_key)
            list_resp3 = self.api_client.post(
                "/openclaw/tmpBookmark/list", {}, new_headers
            )
            assert list_resp3.status_code == 200
            assert list_resp3.json().get("code") == 1, "新Key应该可用"

    # ========== 参数校验测试 ==========

    @allure.story("参数校验")
    @allure.title("测试10a: 添加书签缺少title")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_10a_add_bookmark_missing_title(self, headers_with_login):
        """添加书签时缺少title字段"""
        with allure.step("发送缺少title的请求"):
            data = {
                "url": "https://example.com/no-title"
            }
            response = self.api_client.post(
                "/openclaw/tmpBookmark/add", data, headers_with_login
            )

        with allure.step("验证响应"):
            assert response.status_code == 200
            # 具体行为取决于validate配置，至少不应该500
            response_data = response.json()
            # SetAutoValidate(false, nil, true) 表示禁用了自动validate
            # 所以这个请求可能会成功或返回业务错误

    @allure.story("参数校验")
    @allure.title("测试10b: 添加书签缺少url")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_10b_add_bookmark_missing_url(self, headers_with_login):
        """添加书签时缺少url字段"""
        with allure.step("发送缺少url的请求"):
            data = {
                "title": "缺少URL的书签"
            }
            response = self.api_client.post(
                "/openclaw/tmpBookmark/add", data, headers_with_login
            )

        with allure.step("验证响应"):
            assert response.status_code == 200
            response_data = response.json()
            # 同上，validate已禁用，需要看业务逻辑的处理

    @allure.story("参数校验")
    @allure.title("测试10c: 检查书签存在性缺少url")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_10c_exists_bookmark_missing_url(self, headers_with_login):
        """检查书签存在性时缺少url字段"""
        with allure.step("发送缺少url的请求"):
            data = {
                "title": "只有title"
            }
            response = self.api_client.post(
                "/openclaw/tmpBookmark/exists", data, headers_with_login
            )

        with allure.step("验证响应不应500"):
            assert response.status_code == 200

    @allure.story("参数校验")
    @allure.title("测试10d: 列出书签使用超出范围的limit")
    @pytest.mark.login
    @pytest.mark.openclaw
    def test_10d_list_bookmarks_invalid_limit(self, headers_with_login):
        """使用超出范围的limit值(>500)列出书签"""
        with allure.step("发送超大limit的请求"):
            response = self.api_client.post(
                "/openclaw/tmpBookmark/list", {"limit": 1000}, headers_with_login
            )

        with allure.step("验证响应"):
            assert response.status_code == 200
            # validate已禁用，所以行为取决于业务逻辑
