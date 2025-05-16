#!/bin/bash

# 脚本名称: check_go_env.sh
# 用途: 检查和显示 Go 开发环境信息

# 颜色定义，便于输出高亮
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # 无颜色

echo -e "${YELLOW}=== Go 环境检查脚本 ===${NC}"
echo ""

# 1. 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: Go 未安装，请先安装 Go 环境！${NC}"
    echo "安装指引: https://go.dev/doc/install"
    exit 1
else
    echo -e "${GREEN}Go 已安装${NC}"
    go_version=$(go version)
    echo "Go 版本: $go_version"
fi

# 2. 检查 GOPATH
echo ""
echo -e "${YELLOW}--- GOPATH 信息 ---${NC}"
gopath=$(go env GOPATH)
if [ -z "$gopath" ]; then
    echo -e "${RED}警告: GOPATH 未设置${NC}"
    echo "默认 GOPATH 可能为: ~/go"
else
    echo "GOPATH: $gopath"
    # 检查 GOPATH 目录是否存在
    if [ -d "$gopath" ]; then
        echo -e "${GREEN}GOPATH 目录存在${NC}"
        # 检查 bin、pkg、src 目录
        for dir in bin pkg src; do
            if [ -d "$gopath/$dir" ]; then
                echo "  - $dir: 存在"
            else
                echo "  - $dir: 不存在"
            fi
        done
    else
        echo -e "${RED}错误: GOPATH 目录 $gopath 不存在${NC}"
    fi
fi

# 3. 检查 GO111MODULE
echo ""
echo -e "${YELLOW}--- Go 模块设置 ---${NC}"
go111module=$(go env GO111MODULE)
if [ -z "$go111module" ]; then
    echo "GO111MODULE: 未显式设置（默认根据上下文）"
else
    echo "GO111MODULE: $go111module"
    case $go111module in
        "on")
            echo -e "${GREEN}模块模式启用，推荐用于新项目${NC}"
            ;;
        "off")
            echo -e "${YELLOW}警告: 模块模式关闭，可能影响依赖管理${NC}"
            ;;
        "auto")
            echo "自动模式：根据是否有 go.mod 文件决定"
            ;;
        *)
            echo -e "${RED}未知的 GO111MODULE 设置: $go111module${NC}"
            ;;
    esac
fi

# 4. 检查当前目录的 Go 模块状态
echo ""
echo -e "${YELLOW}--- 当前项目模块信息 ---${NC}"
if [ -f "go.mod" ]; then
    echo -e "${GREEN}当前目录是 Go 模块项目${NC}"
    echo "模块文件: go.mod"
    # 显示模块名称和 Go 版本
    module_name=$(grep "^module" go.mod | awk '{print $2}')
    go_mod_version=$(grep "^go" go.mod | awk '{print $2}')
    echo "模块名称: $module_name"
    echo "Go 模块版本: $go_mod_version"
    # 检查依赖数量
    dep_count=$(grep -A 1000 "^require" go.mod | grep -v "^require" | grep -v "^)" | wc -l)
    echo "直接依赖数量: $dep_count"
else
    echo -e "${YELLOW}当前目录不是 Go 模块项目（无 go.mod 文件）${NC}"
    echo "建议: 使用 'go mod init <module-name>' 初始化模块"
fi

# 5. 检查 GOROOT
echo ""
echo -e "${YELLOW}--- GOROOT 信息 ---${NC}"
goroot=$(go env GOROOT)
if [ -z "$goroot" ]; then
    echo -e "${RED}错误: GOROOT 未设置${NC}"
else
    echo "GOROOT: $goroot"
    if [ -d "$goroot" ]; then
        echo -e "${GREEN}GOROOT 目录存在${NC}"
    else
        echo -e "${RED}错误: GOROOT 目录 $goroot 不存在${NC}"
    fi
fi

# 6. 检查常用 Go 工具
echo ""
echo -e "${YELLOW}--- 常用 Go 工具 ---${NC}"
tools=("gofmt" "goimports" "golint" "golangci-lint")
for tool in "${tools[@]}"; do
    if command -v $tool &> /dev/null; then
        echo -e "${GREEN}$tool: 已安装${NC}"
    else
        echo -e "${YELLOW}$tool: 未安装${NC}"
    fi
done

# 7. 检查环境变量
echo ""
echo -e "${YELLOW}--- 其他 Go 环境变量 ---${NC}"
echo "GOMODCACHE: $(go env GOMODCACHE)"
echo "GOPROXY: $(go env GOPROXY)"
echo "GOSUMDB: $(go env GOSUMDB)"

# 8. 总结和建议
echo ""
echo -e "${YELLOW}=== 环境检查总结 ===${NC}"
if [ -n "$gopath" ] && [ -d "$gopath" ] && [ -f "go.mod" ]; then
    echo -e "${GREEN}Go 环境配置良好，适合模块化开发${NC}"
else
    echo -e "${YELLOW}建议：${NC}"
    if [ -z "$gopath" ] || [ ! -d "$gopath" ]; then
        echo "  - 设置并确保 GOPATH 有效（例如：export GOPATH=~/go）"
    fi
    if [ ! -f "go.mod" ]; then
        echo "  - 在项目目录运行 'go mod init <module-name>' 启用模块"
    fi
    if [ "$go111module" = "off" ]; then
        echo "  - 设置 GO111MODULE=on 以启用模块支持（export GO111MODULE=on）"
    fi
fi
echo -e "${GREEN}检查完成！${NC}"