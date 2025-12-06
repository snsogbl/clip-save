// Base64 编码脚本
// 将文本内容编码为 Base64 格式

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

try {
  // 编码为 Base64
  // 使用 TextEncoder 处理 Unicode 字符
  const encoder = new TextEncoder();
  const bytes = encoder.encode(input);
  const binary = String.fromCharCode(...bytes);
  const encoded = btoa(binary);

  return encoded;
} catch (error) {
  return {
    error: `编码失败: ${error.message || String(error)}`,
  };
}
