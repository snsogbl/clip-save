/**
 * 调用阿里云百炼（DashScope）API 进行文本分析
 * @author ClipSave
 * @param {string} apiKey - 阿里云百炼 API Key
 * @param {string} appId - 应用 ID
 * @param {string} content - 要分析的文本内容，从剪贴板内容获取
 * @returns {object} - API 返回的结果数据
 * @returns {object} - 错误信息
 */

// ===== 配置区域 =====
const apiKey = "your-api-key";
const appId = "your-app-id";

// 从剪贴板项获取内容
if (item.ContentType !== "Text") {
  return {
    error: "只支持文本类型的剪贴板内容",
  };
}

if (apiKey === "your-api-key" || appId === "your-app-id") {
  return {
    error: "请编辑脚本，更新代码中的 apiKey 和 appId 为你的 API Key 和 App ID",
  };
}

const content = item.Content || "";
if (!content) {
  return {
    error: "剪贴板内容为空，无法进行分析",
  };
}

// 构建请求 URL
const url = `https://dashscope.aliyuncs.com/api/v1/apps/${appId}/completion`;

// 创建请求体
const requestBody = {
  input: {
    prompt: content,
  },
  parameters: {},
  debug: {},
};

// ===== 发送请求 =====
try {
  // 发送 HTTP POST 请求
  const response = await fetch(url, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${apiKey}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(requestBody),
  });

  // 检查 HTTP 状态码
  if (!response.ok) {
    const errorText = await response.text();
    return {
      error: `请求失败: ${response.status} ${response.statusText} - ${errorText}`,
    };
  }

  // 读取响应
  const result = await response.json();

  // 处理响应
  if (response.status === 200 && result.output) {
    // 提取 output.text
    let text = result.output.text || "";

    // 清理 JSON 代码块标记
    text = text.replace(/```json\n/g, "").replace(/\n```/g, "");

    // 解析 JSON
    let resultData;
    try {
      resultData = JSON.parse(text);
    } catch (parseError) {
      return {
        error: `解析响应 JSON 失败: ${parseError.message}`,
        rawText: text,
      };
    }

    return {
      success: true,
      data: resultData,
      usage: result.usage,
      requestId: result.request_id,
    };
  } else {
    return {
      error: `请求失败，状态码: ${response.status}`,
      response: result,
    };
  }
} catch (error) {
  return {
    error: `请求失败: ${error.message || String(error)}`,
  };
}

