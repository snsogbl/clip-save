# 快速开始

## 1. 开发模式运行

```bash
# 在项目根目录下
wails dev
```

这会启动开发服务器，应用会自动打开。

## 2. 构建生产版本

```bash
# 构建应用
wails build

# macOS 通用版本（推荐）
wails build -platform darwin/universal
```

构建完成后，可执行文件在 `build/bin/` 目录下。

## 3. 主要功能

### 自动保存剪贴板
应用启动后会自动保存剪贴板变化，无需任何配置。

### 搜索和过滤
- 顶部搜索框：输入关键词实时搜索
- 过滤下拉框：按类型筛选（文本/图片/URL等）

### 操作
- **点击列表项**：查看详情
- **复制按钮**：将内容复制回剪贴板
- **删除按钮**：删除历史记录

## 4. 开发技巧

### 前端开发
```bash
# 前端代码在 frontend/src/ 目录
# 主组件: frontend/src/components/ClipboardHistory.vue
# 修改保存后会自动热重载
```

### 后端开发
```bash
# 后端代码在根目录
# app.go - API 定义
# common/ - 核心逻辑

# 修改后重新生成前端绑定
wails generate module
```

### 调试
开发模式下自动打开 DevTools，可以查看：
- Console 日志
- Network 请求
- Vue DevTools（需要浏览器扩展）

## 5. 目录说明

```
frontend/src/
  ├── App.vue                    # 主应用
  ├── components/
  │   └── ClipboardHistory.vue   # 剪贴板历史组件
  └── main.ts                    # 入口文件

common/
  ├── clipboard.go              # 剪贴板
  ├── clipboard_darwin.go       # macOS 特定代码
  └── db.go                     # 数据库操作

app.go                          # Wails API
main.go                         # 程序入口
```

## 6. 常用命令

```bash
# 开发
wails dev

# 构建
wails build

# 清理
wails clean

# 生成前端绑定
wails generate module

# 安装依赖
go mod tidy
cd frontend && npm install
```

## 7. 数据位置

剪贴板历史保存在：
- macOS: `~/.clipsave/clipboard.db`
- Windows: `C:\Users\<用户>\.clipsave\clipboard.db`
- Linux: `~/.clipsave/clipboard.db`

## 8. 遇到问题？

### 编译错误
```bash
# 确保依赖已安装
go mod tidy
cd frontend && npm install
```

### 前端无法调用后端
```bash
# 重新生成绑定
wails generate module
```

### 剪贴板不工作
检查系统权限设置，macOS 需要授予辅助功能权限。

## 下一步

- 阅读 [README.md](README.md) 了解详细信息
- 查看 [Wails 文档](https://wails.io) 学习更多
- 探索 [Vue 3 文档](https://vuejs.org) 了解前端开发

