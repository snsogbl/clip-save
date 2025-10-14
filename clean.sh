#!/bin/bash

# 剪存项目清理脚本
# 用于清理构建产物和临时文件

echo "🧹 开始清理剪存项目..."

# 清理前端构建产物
if [ -d "frontend/dist" ]; then
    echo "删除前端构建产物..."
    rm -rf frontend/dist
fi

# 清理前端依赖（可选，需要重新安装）
if [ "$1" = "--deep" ]; then
    if [ -d "frontend/node_modules" ]; then
        echo "删除前端依赖..."
        rm -rf frontend/node_modules
    fi
    echo "请运行 'cd frontend && npm install' 重新安装依赖"
fi

# 清理应用构建产物
if [ -d "build/bin" ]; then
    echo "删除应用构建产物..."
    rm -rf build/bin
fi

# 清理 Go 模块缓存（可选）
if [ "$1" = "--deep" ]; then
    echo "清理 Go 模块缓存..."
    go clean -modcache
fi

# 清理临时文件
find . -name "*.tmp" -delete 2>/dev/null
find . -name "*.log" -delete 2>/dev/null
find . -name ".DS_Store" -delete 2>/dev/null

echo "✅ 清理完成！"
echo ""
echo "使用方法："
echo "  ./clean.sh        - 清理构建产物"
echo "  ./clean.sh --deep - 深度清理（包括依赖）"
