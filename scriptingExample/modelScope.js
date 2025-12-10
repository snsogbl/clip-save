/**
 * 调用ModelScope API 进行文本错别字检测
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
    content: `你是严格的中文文本错别字检测助手，仅检测【错别字】，不评判用词好坏、不优化表达，严格遵循以下规则：
                1.  明确错别字定义（仅检测以下4类）：
                    - 形近字错误（如：园→圆、爱→受）
                    - 音近字错误（如：的→得、再→在）
                    - 错写/漏写/多字（如：可爱→小可耐、公园→公圆）
                    - 通用规范汉字错误（不识别生僻字、网络用语，如“可耐”不算错别字，仅算口语化表达）
                2.  绝对不做以下操作：
                    - 不把“常用规范用词”判为错别字（如“原装”“取卡针”是标准用词，不能改为“专用”“SIM卡针”）
                    - 不优化语句通顺度、不替换同义词语
                    - 不纠结标点符号（仅检测文字错误）
                3.  输出格式要求（严格遵守，无错误则明确说明）：
                    - 若有错误：分点列出「错误位置（按正文顺序，如第5字）、错误字、正确字、错误原因」
                    - 若无错误：直接输出“经检测，该文本无错别字，所有用词均符合中文通用规范。”
                    - 最后无需额外补充，不添加修正文本（无错误时不输出修正文本）
                4.  位置判断规则：按文本实际字符顺序计数（含中文、标点，不含空格），不出现“第1个字位于第5个字符”的错误表述。
                5.  最后给出完整的修正后文本`,
  },
  {
    role: "user",
    content: `请检测以下文本的错别字：${content}`,
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
      stream: false, // 关闭流式
      temperature: 0.0, // 低随机性，保证检测结果精准
      max_tokens: 2000, // 足够的生成长度
      top_p: 0.1, // 限制模型输出的随机性，聚焦核心判断
      stop: ["\n\n3."], // 防止模型输出多余内容
    }),
  }
);

// 检查响应状态
if (!response.ok) {
  throw new Error(`请求失败：${response.status} ${response.statusText}`);
}

// 解析 JSON 结果
const result = await response.json();
if (
  result["choices"] &&
  result["choices"][0] &&
  result["choices"][0]["message"] &&
  result["choices"][0]["message"]["content"]
) {
  return result["choices"][0]["message"]["content"];
} else {
  return result;
}
