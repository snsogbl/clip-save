// 时间戳转换脚本
// 自动检测输入是时间戳还是日期时间，并执行相应转换

if (item.ContentType !== 'Text') {
  return {
    error: '只支持文本类型的剪贴板内容'
  };
}

const input = item.Content ? item.Content.trim() : '';

if (!input) {
  return {
    error: '剪贴板内容为空'
  };
}

// 检测是否为时间戳（纯数字，可能是秒级或毫秒级）
function isTimestamp(str) {
  // 移除可能的空格和特殊字符
  const cleaned = str.replace(/[\s,]/g, '');
  // 检查是否为纯数字
  if (!/^\d+$/.test(cleaned)) {
    return false;
  }
  
  const num = parseInt(cleaned, 10);
  // 如果是10位数字，可能是秒级时间戳（2001-2099年）
  // 如果是13位数字，可能是毫秒级时间戳
  // 如果是其他位数，也可能是时间戳，但需要验证合理性
  if (cleaned.length === 10) {
    // 秒级时间戳：2001-01-01 到 2099-12-31
    return num >= 978307200 && num <= 4102444800;
  } else if (cleaned.length === 13) {
    // 毫秒级时间戳：2001-01-01 到 2099-12-31
    return num >= 978307200000 && num <= 4102444800000;
  } else if (cleaned.length >= 8 && cleaned.length <= 15) {
    // 其他长度的时间戳，尝试转换
    return true;
  }
  
  return false;
}

// 检测是否为日期时间字符串
function isDateTime(str) {
  // 常见的日期时间格式
  const datePatterns = [
    /^\d{4}-\d{2}-\d{2}$/, // YYYY-MM-DD
    /^\d{4}\/\d{2}\/\d{2}$/, // YYYY/MM/DD
    /^\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}$/, // YYYY-MM-DD HH:MM:SS
    /^\d{4}\/\d{2}\/\d{2}\s+\d{2}:\d{2}:\d{2}$/, // YYYY/MM/DD HH:MM:SS
    /^\d{2}\/\d{2}\/\d{4}$/, // MM/DD/YYYY
    /^\d{2}-\d{2}-\d{4}$/, // MM-DD-YYYY
  ];
  
  return datePatterns.some(pattern => pattern.test(str));
}

try {
  if (isTimestamp(input)) {
    // 转换为日期时间
    let timestamp = parseInt(input.replace(/[\s,]/g, ''), 10);
    
    // 如果是10位数字，是秒级时间戳；如果是13位，是毫秒级时间戳
    if (input.replace(/[\s,]/g, '').length === 10) {
      timestamp = timestamp * 1000; // 转换为毫秒
    }
    
    const date = new Date(timestamp);
    
    if (isNaN(date.getTime())) {
      return {
        error: '无效的时间戳'
      };
    }
    
    // 格式化为多种日期时间格式
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');
    
    const formats = {
      iso: date.toISOString(),
      local: `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`,
      dateOnly: `${year}-${month}-${day}`,
      chinese: `${year}年${month}月${day}日 ${hours}:${minutes}:${seconds}`,
      timestamp: {
        seconds: Math.floor(timestamp / 1000),
        milliseconds: timestamp
      }
    };
    
    return {
      success: true,
      operation: 'timestamp_to_datetime',
      original: input,
      timestamp: {
        seconds: Math.floor(timestamp / 1000),
        milliseconds: timestamp
      },
      datetime: formats.local,
      formats: formats,
      message: `时间戳转换成功：${formats.local}`
    };
  } else if (isDateTime(input)) {
    // 转换为时间戳
    const date = new Date(input);
    
    if (isNaN(date.getTime())) {
      return {
        error: '无效的日期时间格式'
      };
    }
    
    const timestampMs = date.getTime();
    const timestampS = Math.floor(timestampMs / 1000);
    
    return {
      success: true,
      operation: 'datetime_to_timestamp',
      original: input,
      timestamp: {
        seconds: timestampS,
        milliseconds: timestampMs
      },
      datetime: {
        iso: date.toISOString(),
        local: date.toLocaleString('zh-CN'),
        utc: date.toUTCString()
      },
      message: `日期时间转换成功：时间戳（秒）= ${timestampS}，时间戳（毫秒）= ${timestampMs}`
    };
  } else {
    // 尝试直接解析为日期
    const date = new Date(input);
    
    if (!isNaN(date.getTime())) {
      const timestampMs = date.getTime();
      const timestampS = Math.floor(timestampMs / 1000);
      
      return {
        success: true,
        operation: 'auto_convert',
        original: input,
        timestamp: {
          seconds: timestampS,
          milliseconds: timestampMs
        },
        datetime: {
          iso: date.toISOString(),
          local: date.toLocaleString('zh-CN'),
          utc: date.toUTCString()
        },
        message: `自动转换成功：时间戳（秒）= ${timestampS}`
      };
    }
    
    return {
      error: '无法识别输入格式，请选择时间戳（数字）或日期时间字符串的内容运行脚本'
    };
  }
} catch (error) {
  return {
    error: `转换失败: ${error.message || String(error)}`
  };
}

