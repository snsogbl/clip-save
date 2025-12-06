// Base64 解码脚本
// 将 Base64 编码的字符串解码为原始文本

if (item.ContentType !== "Text") {
  return {
    error: "只支持文本类型的剪贴板内容",
  };
}

const input = item.Content || "";

if (!input) {
  return {
    error: "剪贴板内容为空",
  };
}

// 检测是否为 Base64 编码的字符串
function isBase64(str) {
  try {
    // Base64 字符串通常只包含 A-Z, a-z, 0-9, +, /, = 字符
    // 并且长度是 4 的倍数（可能包含填充的 =）
    const base64Regex = /^[A-Za-z0-9+/]*={0,2}$/;
    if (!base64Regex.test(str.trim())) {
      return false;
    }

    // 尝试解码，如果成功则可能是 Base64
    const decoded = atob(str.trim());
    // 检查解码后的内容是否是可打印的 ASCII 字符或有效的 UTF-8
    return /^[\x20-\x7E]*$/.test(decoded) || decoded.length > 0;
  } catch (e) {
    return false;
  }
}

const trimmedInput = input.trim();

if (!isBase64(trimmedInput)) {
  return {
    error: "输入内容不是有效的 Base64 编码字符串",
  };
}

try {
  // 解码 Base64
  const decoded = atob(trimmedInput);
  
  // 尝试使用 TextDecoder 处理 UTF-8 编码
  try {
    const decoder = new TextDecoder('utf-8');
    const bytes = new Uint8Array(decoded.length);
    for (let i = 0; i < decoded.length; i++) {
      bytes[i] = decoded.charCodeAt(i);
    }
    return decoder.decode(bytes);
  } catch (e) {
    // 如果 TextDecoder 失败，直接返回解码结果
    return decoded;
  }
} catch (error) {
  return {
    error: `解码失败: ${error.message || String(error)}`,
  };
}

