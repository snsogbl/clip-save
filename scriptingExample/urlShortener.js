/**
 * 使用 TinyURL API 生成短链接
 * @author ClipSave
 * @param {string} url - 输入的 URL，从剪贴板内容获取
 * @returns {object} - 短链接
 * @returns {object} - 错误信息
 */

// 检查剪贴板内容是否为 URL
if (item.ContentType !== 'URL') {
  return {
    error: '只支持 URL 类型的剪贴板内容'
  };
}

const url = item.Content ? item.Content.trim() : '';

if (!url) {
  return {
    error: '剪贴板内容为空，无法生成短链接'
  };
}

// 验证是否为有效的 URL
try {
  new URL(url);
} catch (e) {
  return {
    error: '剪贴板内容不是有效的 URL'
  };
}

// 使用 TinyURL API 生成短链接
try {
  const apiUrl = `https://tinyurl.com/api-create.php?url=${encodeURIComponent(url)}`;
  
  const response = await fetch(apiUrl, {
    method: 'GET',
  });

  if (!response.ok) {
    return {
      error: `生成短链接失败: ${response.status} ${response.statusText}`
    };
  }

  const shortUrl = await response.text();
  
  if (!shortUrl || shortUrl.startsWith('Error')) {
    return {
      error: '生成短链接失败，请检查 URL 是否有效'
    };
  }

  // 直接返回短链接字符串，方便复制和使用
  return shortUrl.trim();
} catch (error) {
  return {
    error: `生成短链接失败: ${error.message || String(error)}`
  };
}

