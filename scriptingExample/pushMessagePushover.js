/**
 * 将剪贴板内容推送到 Pushover
 * @author ClipSave
 * @param {string} token - Pushover token
 * @param {string} user - Pushover user key
 * @param {string} device - Pushover device name
 * @param {string} message - 消息内容，从剪贴板内容获取
 * @returns {object} - 推送结果
 * @returns {object} - 错误信息
 */

// 从剪贴板项获取消息内容
if (item.ContentType !== 'Text') {
  return {
    error: '不支持的类型'
  };
}
const message = item.Content || '';
if (!message) {
  return {
    error: '剪贴板内容为空，无法发送推送'
  };
}

// 构建请求数据
const reqdata = {
  token: 'your pushover token',
  user: 'your pushover user key',
  device: 'your device name',
  message: message
};

// 发送 POST 请求到 Pushover API
try {
  const urlStr = 'https://api.pushover.net/1/messages.json';
  
  const response = await fetch(urlStr, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(reqdata)
  });

  if (!response.ok) {
    const errorText = await response.text();
    return {
      error: `推送失败: ${response.status} ${response.statusText} - ${errorText}`
    };
  }

  const resultBody = await response.json();
  
  return {
    success: true,
    message: '推送成功',
    response: resultBody
  };
} catch (error) {
  return {
    error: `推送失败: ${error.message || String(error)}`
  };
}

