# 子模块管理相关命令
.PHONY: submodule-update submodule-status submodule-all help

help:
	@echo "子模块管理命令："
	@echo "  make submodule-update   - 更新子模块到最新版本"
	@echo "  make submodule-status   - 检查子模块状态"
	@echo "  make submodule-all      - 执行初始化和更新操作"


# 更新子模块到最新版本
submodule-update:
	@echo "更新子模块到最新版本..."
	git submodule update --remote

# 检查子模块状态
submodule-status:
	@echo "检查子模块状态..."
	git submodule status

# 完整初始化和更新（包括嵌套子模块）
submodule-all:
	@echo "完整初始化和更新子模块..."
	git submodule update --init --recursive

# 设置默认目标
.DEFAULT_GOAL := help 