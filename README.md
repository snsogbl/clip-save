# 剪存 - 剪贴板历史管理工具

一个基于 Wails + Vue 3 + TypeScript 的多平台剪贴板历史管理工具。

## 已上架App Store
https://apps.apple.com/us/app/剪存/id6754015301

## 功能特性
- [新增自定义脚本-点击查看](scriptingExample/README.md)
- 自动保存剪贴板历史
- 支持文本、图片、文件等多种类型
- 实时搜索和过滤功能（支持搜索图片中的 OCR 文字）
- 本地 SQLite 数据库存储（所有数据只存本地,隐私安全有保障）
- 个性化设置：密码保护、自动清理、快捷键配置

## 项目截图

<img src="https://raw.githubusercontent.com/snsogbl/clip-save/master/screenshots/clipboard-1.png" alt="Capture" width="960">


## 技术栈

- **后端**: Go + Wails v2
- **前端**: Vue 3 + TypeScript + Vite
- **数据库**: SQLite3
- **剪贴板**: golang.design/x/clipboard

## 安装依赖

### 1. 安装 Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 2. 安装项目依赖

```bash
# Go 依赖
go mod tidy

# 前端依赖
cd frontend
npm install
```

## 开发

### 启动开发服务器

```bash
wails dev
```

这将启动热重载开发服务器：
- 后端 Go 代码修改会自动重新编译
- 前端 Vue 代码修改会自动热重载

### 开发时的调试

开发模式下，应用会自动打开开发者工具，可以查看：
- Console 日志
- Network 请求
- 前端组件状态

## 构建

### 构建生产版本

```bash
wails build
```

构建完成后，可执行文件将位于 `build/bin/` 目录下。

### macOS 构建选项

```bash
# 构建 Intel 版本
wails build -platform darwin/amd64

# 构建 Apple Silicon 版本
wails build -platform darwin/arm64

# 构建通用二进制（推荐）
wails build -platform darwin/universal
```

### Windows 构建

```bash
wails build -platform windows/amd64
```

### Linux 构建

```bash
wails build -platform linux/amd64
```

## 项目结构

```
.
├── main.go                      # 主程序入口
├── app.go                       # Wails 应用结构和 API
├── wails.json                   # Wails 配置文件
├── go.mod                       # Go 依赖管理
├── common/                      # 共享代码
│   ├── clipboard.go             # 剪贴板逻辑
│   ├── clipboard_darwin.go      # macOS 特定代码
│   └── db.go                    # 数据库操作
├── frontend/                    # 前端代码
│   ├── src/
│   │   ├── App.vue              # 主应用组件
│   │   ├── components/
│   │   │   └── ClipboardHistory.vue  # 剪贴板历史组件
│   │   ├── main.ts              # 前端入口
│   │   └── style.css            # 全局样式
│   ├── index.html               # HTML 模板
│   ├── package.json             # 前端依赖
│   ├── vite.config.ts           # Vite 配置
│   └── tsconfig.json            # TypeScript 配置
└── build/                       # 构建资源和输出
    ├── bin/                     # 编译后的可执行文件
    ├── appicon.png              # 应用图标
    └── darwin/                  # macOS 特定配置
```

## API 说明

### 后端 API（Go）

在 `app.go` 中定义的方法会自动暴露给前端：

- `SearchClipboardItems(keyword, filterType, limit)` - 搜索剪贴板项目
- `GetClipboardItems(limit)` - 获取剪贴板列表
- `GetClipboardItemByID(id)` - 根据 ID 获取项目
- `CopyToClipboard(id)` - 复制项目到剪贴板
- `DeleteClipboardItem(id)` - 删除项目
- `GetStatistics()` - 获取统计信息

### 前端调用示例

```typescript
import { SearchClipboardItems } from '../wailsjs/go/main/App'

// 搜索剪贴板项目
const items = await SearchClipboardItems('关键词', '所有类型', 100)
```

## 使用说明

1. 启动应用后，它会在后台自动保存系统剪贴板
2. 每次复制内容时，都会自动保存到历史记录
3. 使用顶部搜索框可以快速查找历史记录（支持搜索图片中的 OCR 文字）
4. 使用过滤器可以按类型筛选内容（文本/图片/URL等）
5. 点击任意历史记录可以查看详情
6. 对于图片类型，点击"提取文字"按钮可以查看 OCR 识别结果
7. 点击"复制"按钮可以将内容复制回剪贴板
8. 点击"删除"按钮可以删除历史记录

### OCR 文字识别功能

- **自动识别**：复制图片后，应用会自动在后台识别图片中的文字（macOS 10.15+）
- **智能搜索**：识别后的文字会自动加入搜索索引，输入关键词即可找到包含该文字的图片
- **手动提取**：在图片详情页面，点击"提取文字"按钮可以查看完整的 OCR 识别结果
- **支持语言**：支持中文（简体/繁体）和英文识别
- **性能优化**：相同图片只识别一次，自动复用识别结果，提升性能

## 数据存储

剪贴板历史保存在：`~/.clipsave/clipboard.db`

数据库会自动创建，包含以下字段：
- ID - 唯一标识符
- Content - 内容文本
- ContentType - 内容类型
- ImageData - 图片数据（PNG格式）
- Timestamp - 时间戳
- Source - 来源
- CharCount - 字符数
- WordCount - 单词数
- OCRText - OCR 识别的文字内容（图片类型专用，用于搜索）

## 系统要求

- **macOS**: 10.15 Catalina 或更高版本（OCR 功能需要）
- **Windows**: Windows 10/11（1809或更高版本）+ WebView2
- **Linux**: 支持 WebKit2GTK 的发行版
- **Go**: 1.21 或更高版本
- **Node.js**: 16 或更高版本

## 开发注意事项

### 更新 Go API 后

每次修改 `app.go` 中的方法后，需要重新生成前端绑定：

```bash
wails generate module
```

或者使用开发模式，会自动生成：

```bash
wails dev
```

### 前端开发

前端使用 Vite + Vue 3 + TypeScript：
- 支持 TypeScript 类型检查
- 使用 Composition API
- 自动导入 Wails 绑定
- 热模块替换（HMR）

### CGO 依赖

项目使用了 CGO（用于 SQLite 和剪贴板操作），构建时需要：
- macOS: 需要 Xcode Command Line Tools
- Windows: 需要 MinGW-w64
- Linux: 需要 gcc

## 常见问题

### 1. 构建失败

确保安装了所有依赖：
```bash
# macOS
xcode-select --install

# Windows
# 安装 MSYS2 和 MinGW-w64

# Linux
sudo apt-get install build-essential libgtk-3-dev libwebkit2gtk-4.0-dev
```

### 2. 前端无法调用后端 API

确保已经运行了 `wails generate module` 生成前端绑定。

## 许可证

MIT License

## 赞赏支持

💗 **请作者喝杯咖啡**

如果这个软件帮你省下了时间（或者至少没让你抓狂），欢迎请作者喝杯咖啡 ☕ 或者啤酒 🍺！

每一杯咖啡都会神奇地转化为：新功能、更少的 bug（我尽量）、以及作者熬夜写代码时的精神支柱。

你的 Star ⭐ 和捐赠是我修复 bug 和添加新功能的原动力（也是我买咖啡的经费）！

### 微信扫一扫

![赞赏码](frontend/src/assets/static/zs.png)

**"请我喝杯咖啡吧 ☕"**

### 给项目点个 Star ⭐

如果这个项目对你有帮助，欢迎给项目点个 Star ⭐，这是对我最大的支持！

## 致谢

- [Wails](https://wails.io) - 构建桌面应用的框架
- [Vue 3](https://vuejs.org) - 渐进式 JavaScript 框架
- [golang.design/x/clipboard](https://github.com/golang-design/clipboard) - 跨平台剪贴板库
- [json-editor-vue](https://github.com/cloydlau/json-editor-vue)
- [highlight.js](https://github.com/highlightjs/highlight.js)