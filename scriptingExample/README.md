# 剪存脚本使用教程

欢迎使用剪存脚本功能！脚本功能允许你自定义处理剪贴板内容，实现自动化操作。

## 📖 快速开始

### 已有脚本
快速访问：[钉钉消息推送](pushMessageDingtalk.js) | [Pushover 消息推送](pushMessagePushover.js) | [JWT Token 生成](jwt.js) | [文本信息提取](textExtract.js) | [URL 短链接生成](urlShortener.js) | [Base64 编码](base64Encode.js) | [Base64 解码](base64Decode.js) | [时间戳转换](timestampConverter.js) | 等等...

### 1. 创建脚本

1. 打开剪存应用
2. 点击右上角设置按钮 ⚙️
3. 在设置中找到"脚本管理"
4. 点击"新建脚本"按钮
5. 填写脚本信息：
   - **名称**：脚本的显示名称（如：Base64 编码）
   - **描述**：脚本的功能说明
   - **触发时机**：选择脚本的执行时机
     - `手动执行`：通过菜单或按钮手动触发
     - `保存后执行`：剪贴板内容保存后自动执行
   - **内容类型**：限制脚本只在特定内容类型时执行（留空表示所有类型）
   - **关键词**：限制脚本只在内容包含特定关键词时执行（留空表示不过滤）
     - 支持普通字符串匹配（不区分大小写）
     - **支持正则表达式**：以 `/` 开头和结尾的关键词会被识别为正则表达式
       - 格式：`/pattern/` 或 `/pattern/flags`
       - 示例：`/^test$/i`（匹配完整的 "test"，不区分大小写）、`/\d+/`（匹配数字）
       - 支持的正则标志：`i`（忽略大小写）、`g`（全局匹配）、`m`（多行匹配）等
6. 在代码编辑器中编写 JavaScript 代码
7. 点击"保存"完成创建

### 2. 使用脚本

#### 手动执行脚本
1. 选中一个剪贴板项
2. 点击标题栏的"运行脚本"按钮
3. 在弹出的对话框中选择要执行的脚本
4. 查看执行结果

#### 自动执行脚本
- 设置为"保存后执行"的脚本会在剪贴板内容保存时自动运行
- 确保脚本的过滤条件（内容类型、关键词）匹配当前内容

### 3. 脚本编写基础

脚本运行在浏览器环境中，可以使用所有浏览器 API。脚本会接收到一个 `item` 对象，包含当前剪贴板项的信息：

```javascript
// item 对象结构
{
  ID: "唯一标识",
  Content: "内容文本",
  ContentType: "Text|Image|File|URL|Color|JSON",
  Timestamp: "时间戳",
  Source: "来源应用",
  CharCount: 100,
  WordCount: 20,
  IsFavorite: 0
}
```

#### 返回值

脚本可以返回以下类型的值：

1. **字符串**：直接返回处理后的文本
2. **对象**：返回包含 `error` 字段的对象表示错误
3. **Promise**：支持异步操作

示例：

```javascript
// 简单返回
return "处理后的内容";

// 返回错误
return {
  error: "处理失败：原因说明"
};

// 异步操作
async function process() {
  const result = await someAsyncOperation();
  return result;
}
return process();
```

#### 调试技巧

在开发脚本时，可以使用以下方式查看调试信息：

1. **使用 `alert()` 函数**：显示弹窗提示，适合快速查看中间结果
   ```javascript
   alert("当前内容: " + item.Content);
   alert("处理结果: " + result);
   ```

2. **使用 `return` 返回值**：返回值会显示在脚本执行结果区域
   ```javascript
   // 返回调试信息
   return {
     debug: "调试信息",
     content: item.Content,
     result: processedResult
   };
   
   // 或者直接返回字符串
   return "调试信息: " + JSON.stringify(item, null, 2);
   ```

## 📚 现有脚本示例

我们提供了一些实用的脚本示例，你可以直接复制使用或作为参考：

### 🔐 [Base64 编码](base64Encode.js)
将文本内容编码为 Base64 格式。

**功能特点：**
- 支持 Unicode 字符的编码
- 使用 TextEncoder 确保正确编码多字节字符
- 简单直接，专注于编码操作

**适用场景：** 需要将文本编码为 Base64 的场景

---

### 🔓 [Base64 解码](base64Decode.js)
将 Base64 编码的字符串解码为原始文本。

**功能特点：**
- 自动检测输入是否为有效的 Base64 字符串
- 支持 UTF-8 编码的文本解码
- 提供清晰的错误提示

**适用场景：** 需要将 Base64 字符串解码为原始文本的场景

---

### 🎫 [JWT Token 生成](jwt.js)
根据剪贴板内容生成 JWT（JSON Web Token）令牌。

**功能特点：**
- 提取剪贴板内容
- 使用 Web Crypto API 生成签名
- 返回完整的 JWT token

**适用场景：** API 开发、身份验证测试

---

### 📝 [文本信息提取](textExtract.js)
从文本中提取各种结构化信息，包括邮箱、URL、手机号、身份证号、银行卡号等。

**功能特点：**
- 提取邮箱地址
- 提取 URL 链接
- 提取手机号码（支持中国格式）
- 提取 IP 地址
- 提取身份证号（15位/18位，带校验）
- 提取银行卡号（带 Luhn 算法校验）
- 返回提取结果和统计信息

**适用场景：** 数据清洗、信息整理、隐私保护检查

---

### 📤 [Pushover 消息推送](pushMessagePushover.js)
将剪贴板内容通过 Pushover 推送到你的设备。

**功能特点：**
- 支持 Pushover API
- 可配置设备、优先级等参数
- 异步推送，不阻塞主流程

**适用场景：** 跨设备通知、重要内容提醒

---

### 💬 [钉钉消息推送](pushMessageDingtalk.js)
将剪贴板内容推送到钉钉群聊。

**功能特点：**
- 支持钉钉机器人 Webhook API
- 使用通用 `httpRequest` 函数绕过 CORS 限制
- 自动处理响应和错误信息
- 配置简单，只需填写 access_token

**配置说明：**
1. 在钉钉群中添加自定义机器人
2. 获取 Webhook URL 中的 `access_token`
3. 在脚本中修改 `access_token` 变量

**适用场景：** 团队协作、工作通知、重要信息推送

---

### 🔗 [URL 短链接生成](urlShortener.js)
将长 URL 转换为短链接。

**功能特点：**
- 支持 TinyURL API
- 自动检测 URL 类型内容
- 返回可直接使用的短链接

**适用场景：** 分享长链接、社交媒体发布

---

### ⏰ [时间戳转换](timestampConverter.js)
在 Unix 时间戳和可读日期时间之间转换。

**功能特点：**
- 自动识别时间戳格式（秒/毫秒）
- 支持多种日期格式输出
- 双向转换（时间戳 ↔ 日期）

**适用场景：** 日志分析、时间格式化

---

## 🛠️ 脚本开发技巧

### 1. 错误处理
始终使用 try-catch 处理可能的错误：

```javascript
try {
  // 你的代码
  return result;
} catch (error) {
  return {
    error: `操作失败: ${error.message}`
  };
}
```

### 2. 内容类型检查
根据内容类型执行不同逻辑：

```javascript
if (item.ContentType !== "Text") {
  return {
    error: "此脚本仅支持文本类型"
  };
}
```

### 3. 异步操作
使用 async/await 处理异步操作：

```javascript
async function fetchData() {
  const response = await fetch('https://api.example.com/data');
  const data = await response.json();
  return data;
}
return await fetchData();
```

### 4. 使用浏览器 API
脚本可以访问所有浏览器 API：

- `fetch()` - HTTP 请求（注意：可能受 CORS 限制）
- `btoa()` / `atob()` - Base64 编码/解码
- `TextEncoder` / `TextDecoder` - 文本编码转换
- `crypto.subtle` - 加密操作
- `Date` - 日期时间处理
- 等等...

### 5. 使用 httpRequest 函数绕过 CORS
当需要调用外部 API 时，如果遇到 CORS 限制，可以使用内置的 `httpRequest` 函数：

```javascript
// httpRequest 函数签名
// httpRequest(method, url, headersJson, bodyJson)
// - method: HTTP 方法（GET, POST, PUT, DELETE 等）
// - url: 请求 URL
// - headersJson: 请求头 JSON 字符串，如 '{"Content-Type": "application/json"}'
// - bodyJson: 请求体 JSON 字符串（GET 请求可为空字符串）

// 示例：发送 POST 请求
const responseJson = await httpRequest(
  'POST',
  'https://api.example.com/endpoint',
  JSON.stringify({ 'Content-Type': 'application/json' }),
  JSON.stringify({ key: 'value' })
);

const response = JSON.parse(responseJson);
console.log(response.status);    // HTTP 状态码
console.log(response.body);      // 响应体（可能是对象或字符串）
```

**注意：** `httpRequest` 函数由脚本执行器自动注入，无需手动导入。

### 6. 关键词过滤的正则表达式
在脚本配置中，关键词字段支持正则表达式，可以更精确地匹配内容：

```javascript
// 示例关键词配置：
// - "test"          → 普通字符串匹配（不区分大小写）
// - "/^test$/"      → 精确匹配 "test"（区分大小写）
// - "/^test$/i"     → 精确匹配 "test"（不区分大小写）
// - "/\d{4}-\d{2}-\d{2}/" → 匹配日期格式（如：2024-01-01）
// - "/^https?:\/\//" → 匹配以 http:// 或 https:// 开头的 URL
```

正则表达式匹配失败时会自动回退到普通字符串匹配，确保脚本不会因为正则错误而无法执行。

## 🤝 贡献脚本

我们欢迎社区贡献更多实用的脚本！

### 如何贡献

1. **Fork 本项目**
2. **创建你的脚本文件**
   - 在 `scriptingExample` 目录下创建新的 `.js` 文件
   - 使用有意义的文件名（如：`myAwesomeScript.js`）
3. **添加脚本说明**
   - 在脚本文件开头添加注释说明功能
   - 更新本 README，添加你的脚本介绍
4. **提交 Pull Request**
   - 描述脚本的功能和使用场景
   - 确保代码格式清晰，注释完整

### 脚本规范

- ✅ 添加清晰的注释说明
- ✅ 处理错误情况
- ✅ 检查内容类型和格式
- ✅ 返回有意义的错误信息
- ✅ 使用有意义的变量名
- ❌ 避免执行危险操作
- ❌ 不要访问敏感信息

## 📞 反馈与支持

如果你在使用脚本功能时遇到问题，或有好的建议：

- 🐛 [提交 Issue](https://github.com/snsogbl/clip-save/issues)
- ⭐ 给项目点个 Star 支持我们！

---

