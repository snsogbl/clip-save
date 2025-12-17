/**
 * 调用ModelScope API 通用AI流式请求
 * @author ClipSave
 * @param {string} apiKey - API Key
 * @param {string} model - 模型名称
 * @param {string} messages - 消息列表
 * @returns {object} - API 返回的结果数据
 */

// 从剪贴板项获取内容
if (item.ContentType !== "Text") {
  return {
    error: "只支持文本类型的剪贴板内容",
  };
}

const content = item.Content || "";
if (!content) {
  return {
    error: "剪贴板内容为空，无法进行分析",
  };
}

const apiKey = "your api key";
const model = "Qwen/Qwen3-235B-A22B-Instruct-2507";
const messages = [
  {
    role: "system",
    content: `你是专业知识型 AI 助手，核心响应规则如下：
                1、优先直接给出精准答案，不额外补充原理、背景、拓展说明（仅当我明确询问 “为什么”“原理是什么”“机制是什么” 等相关问题时，才补充对应解释）；
                2、答案简洁明了，避免冗余表述，聚焦问题核心诉求；
                3、若问题涉及多维度答案，用清晰结构（分点 / 编号）呈现，方便快速读取；
                4、保持专业性和准确性，不添加无关话术，直接交付有效信息。`,
  },
  {
    role: "user",
    content: `${content}`,
  },
];

const response = await fetch(
  "https://api-inference.modelscope.cn/v1/chat/completions",
  {
    method: "POST",
    headers: {
      Authorization: `Bearer ${apiKey}`, // Token 认证
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      model: model,
      messages: messages,
      stream: true, // 开启流式
      temperature: 0.0, // 低随机性，保证检测结果精准
      max_tokens: 500, // 足够的生成长度
      top_p: 0.1, // 限制模型输出的随机性，聚焦核心判断
    }),
  }
);

// 检查响应状态
if (!response.ok) {
  throw new Error(`请求失败：${response.status} ${response.statusText}`);
}

// 流式响应处理
const reader = response.body.getReader();
const decoder = new TextDecoder();
let fullContent = "";
let buffer = "";
let chunkBuffer = "";
let lastEmitTime = 0;
const EMIT_INTERVAL = 16; // 每 16ms 最多发送一次（约 60fps）

// 发送流开始事件
window.dispatchEvent(new CustomEvent('script-stream-start', {
  detail: { itemId: item.ID }
}));

// 处理数据行并提取内容
function processDataLine(data) {
  if (data === "[DONE]") {
    // 发送剩余缓冲区并结束
    if (chunkBuffer.length > 0) {
      window.dispatchEvent(new CustomEvent('script-stream-chunk', {
        detail: { itemId: item.ID, chunk: chunkBuffer }
      }));
    }
    window.dispatchEvent(new CustomEvent('script-stream-end', {
      detail: { itemId: item.ID }
    }));
    return true; // 表示流结束
  }

  try {
    const json = JSON.parse(data);
    const content = json.choices?.[0]?.delta?.content;
    if (content) {
      fullContent += content;
      chunkBuffer += content;

      // 时间节流发送
      const now = Date.now();
      if (now - lastEmitTime >= EMIT_INTERVAL && chunkBuffer.length > 0) {
        window.dispatchEvent(new CustomEvent('script-stream-chunk', {
          detail: { itemId: item.ID, chunk: chunkBuffer }
        }));
        chunkBuffer = "";
        lastEmitTime = now;
      }
    }
  } catch (e) {
    console.warn("解析流式数据失败:", e, "数据:", data);
  }
  return false;
}

try {
  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    // 解码并分割成行
    buffer += decoder.decode(value, { stream: true });
    const lines = buffer.split("\n");
    buffer = lines.pop() || ""; // 保留最后一个不完整的行

    // 处理每一行
    for (const line of lines) {
      if (!line.trim()) continue;
      if (!line.startsWith("data: ")) continue;

      const data = line.slice(6);
      if (processDataLine(data)) {
        return fullContent || "无响应内容";
      }
    }
  }

  // 处理最后一行
  if (buffer.trim() && buffer.startsWith("data: ")) {
    const data = buffer.slice(6);
    if (data !== "[DONE]") {
      processDataLine(data);
    }
  }

  // 发送剩余缓冲区并结束
  if (chunkBuffer.length > 0) {
    window.dispatchEvent(new CustomEvent('script-stream-chunk', {
      detail: { itemId: item.ID, chunk: chunkBuffer }
    }));
  }
  window.dispatchEvent(new CustomEvent('script-stream-end', {
    detail: { itemId: item.ID }
  }));

  return fullContent || "无响应内容";
} finally {
  reader.releaseLock();
}
