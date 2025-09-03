#!/usr/bin/env python3
"""
测试运行脚本
用于快速执行不同类型的测试
"""
import subprocess
import sys
import argparse
from pathlib import Path


def run_command(command: list) -> int:
    """运行命令并返回退出码"""
    print(f"执行命令: {' '.join(command)}")
    try:
        result = subprocess.run(command, check=False)
        return result.returncode
    except FileNotFoundError as e:
        print(f"命令未找到: {e}")
        return 1


def main():
    parser = argparse.ArgumentParser(description="ForgetURL API 测试运行器")
    parser.add_argument("--type", choices=["all", "login", "no-login", "space", "page"], 
                       default="all", help="测试类型")
    parser.add_argument("--report", choices=["html", "allure"], 
                       help="报告类型")
    parser.add_argument("--verbose", "-v", action="store_true", 
                       help="详细输出")
    
    args = parser.parse_args()
    
    # 基础命令
    base_cmd = ["python", "-m", "pytest"]
    
    if args.verbose:
        base_cmd.append("-v")
    
    # 根据类型添加标记
    if args.type == "login":
        base_cmd.extend(["-m", "login"])
    elif args.type == "no-login":
        base_cmd.extend(["-m", "no_login"])
    elif args.type == "space":
        base_cmd.extend(["-m", "space"])
    elif args.type == "page":
        base_cmd.extend(["-m", "page"])
    
    # 添加报告选项
    if args.report == "html":
        base_cmd.extend(["--html=report.html", "--self-contained-html"])
    elif args.report == "allure":
        base_cmd.extend(["--alluredir=allure-results"])
    
    # 运行测试
    exit_code = run_command(base_cmd)
    
    # 如果是allure报告，尝试打开
    if args.report == "allure" and exit_code == 0:
        print("\n生成 Allure 报告...")
        serve_cmd = ["allure", "serve", "allure-results"]
        run_command(serve_cmd)
    
    return exit_code


if __name__ == "__main__":
    sys.exit(main())
