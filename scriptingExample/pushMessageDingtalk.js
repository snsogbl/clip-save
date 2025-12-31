/**
 * 将剪贴板内容推送到钉钉群
 * @author ClipSave
 * @param {string} access_token - 钉钉机器人 access_token
 * @param {string} message - 消息内容，从剪贴板内容获取
 * @returns {object} - 推送结果
 * @returns {object} - 错误信息
 */

import { csRequest } from "@clipsave/api";

const access_token = "your-access-token";

// 从剪贴板项获取消息内容
if (item.ContentType !== "Text") {
  return {
    error: "只支持文本类型的剪贴板内容",
  };
}

if (access_token === "your-access-token") {
  return {
    error: "请编辑脚本，更新代码中的 access_token 为你的钉钉机器人 access_token",
  };
}

const message = item.Content || "";
if (!message) {
  return {
    error: "剪贴板内容为空，无法发送推送",
  };
}

// ===== 配置区域 =====
// 钉钉机器人 Webhook URL（从钉钉群机器人设置中获取）
const webhookUrl =
  "https://oapi.dingtalk.com/robot/send?access_token=" + access_token;

// ===== 发送消息 =====
try {
  // 如果设置了加签密钥，生成签名并添加到 URL
  let finalUrl = webhookUrl;

  // 构建消息体
  const messageBody = {
    msgtype: "text",
    text: {
      content: message,
    },
  };

  // 使用通用 csRequest 函数发送请求
  const responseJson = await csRequest(
    "POST",
    finalUrl,
    JSON.stringify({ "Content-Type": "application/json" }),
    JSON.stringify(messageBody)
  );

  // 解析响应
  const response = JSON.parse(responseJson);

  // 检查 HTTP 状态码
  if (response.status !== 200) {
    return {
      error: `请求失败: ${response.status} ${response.statusText}`,
    };
  }

  // 解析响应体（可能是字符串或对象）
  let resultBody = response.body;
  if (typeof resultBody === "string") {
    try {
      resultBody = JSON.parse(resultBody);
    } catch (e) {
      // 如果不是 JSON，直接使用字符串
    }
  }

  // 检查钉钉返回的错误码
  if (resultBody.errcode !== undefined && resultBody.errcode !== 0) {
    return {
      error: `推送失败: ${resultBody.errmsg || "未知错误"}`,
    };
  }

  return {
    success: true,
    message: "推送成功",
    response: resultBody,
  };
} catch (error) {
  return {
    error: `推送失败: ${error.message || String(error)}`,
  };
}
