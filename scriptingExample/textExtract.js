// 文本提取脚本
// 从文本中提取邮箱、URL、手机号等信息

if (item.ContentType !== "Text") {
  return {
    error: "只支持文本类型的剪贴板内容",
  };
}

const text = item.Content || "";

if (!text) {
  return {
    error: "剪贴板内容为空",
  };
}

// 提取邮箱地址
const emailRegex = /[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}/g;
const emails = [...new Set(text.match(emailRegex) || [])];

// 提取 URL
const urlRegex = /https?:\/\/[^\s<>"{}|\\^`\[\]]+/gi;
const urls = [...new Set(text.match(urlRegex) || [])];

// 提取手机号（支持中国大陆手机号格式）
// 匹配 11 位数字，可能包含空格、横线等分隔符
const phoneRegex = /1[3-9]\d[\s-]?\d{4}[\s-]?\d{4}/g;
const phones = [...new Set(text.match(phoneRegex) || [])];

// 提取 IP 地址
const ipRegex = /\b(?:\d{1,3}\.){3}\d{1,3}\b/g;
const ips = [...new Set(text.match(ipRegex) || [])];

// 提取身份证号（18位或15位）- 严格验证
function isValidIdCard(idCard) {
  // 移除空格和横线
  const cleaned = idCard.replace(/[\s-]/g, '');
  
  // 18位身份证验证
  if (cleaned.length === 18) {
    // 前17位必须是数字，最后一位是数字或X
    if (!/^\d{17}[\dXx]$/.test(cleaned)) {
      return false;
    }
    
    // 验证出生日期（第7-14位）
    const year = parseInt(cleaned.substring(6, 10));
    const month = parseInt(cleaned.substring(10, 12));
    const day = parseInt(cleaned.substring(12, 14));
    
    if (year < 1900 || year > new Date().getFullYear()) {
      return false;
    }
    if (month < 1 || month > 12) {
      return false;
    }
    if (day < 1 || day > 31) {
      return false;
    }
    
    // 验证日期有效性
    const date = new Date(year, month - 1, day);
    if (date.getFullYear() !== year || date.getMonth() !== month - 1 || date.getDate() !== day) {
      return false;
    }
    
    // 验证校验位（第18位）
    const weights = [7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2];
    const checkCodes = ['1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'];
    let sum = 0;
    for (let i = 0; i < 17; i++) {
      sum += parseInt(cleaned[i]) * weights[i];
    }
    const checkCode = checkCodes[sum % 11];
    if (checkCode !== cleaned[17].toUpperCase()) {
      return false;
    }
    
    return true;
  }
  
  // 15位身份证验证（旧版）
  if (cleaned.length === 15) {
    // 必须是15位数字
    if (!/^\d{15}$/.test(cleaned)) {
      return false;
    }
    
    // 验证出生日期（第7-12位，YYMMDD）
    const year = parseInt('19' + cleaned.substring(6, 8));
    const month = parseInt(cleaned.substring(8, 10));
    const day = parseInt(cleaned.substring(10, 12));
    
    if (month < 1 || month > 12) {
      return false;
    }
    if (day < 1 || day > 31) {
      return false;
    }
    
    // 验证日期有效性
    const date = new Date(year, month - 1, day);
    if (date.getFullYear() !== year || date.getMonth() !== month - 1 || date.getDate() !== day) {
      return false;
    }
    
    return true;
  }
  
  return false;
}

// 提取身份证号
const idCardMatches = text.match(/\b\d{15}\b|\b\d{17}[\dXx]\b/g) || [];
const idCards = [...new Set(idCardMatches.filter(id => isValidIdCard(id)))];

// 提取银行卡号 - 使用Luhn算法验证
function isValidBankCard(cardNumber) {
  // 移除空格和横线
  const cleaned = cardNumber.replace(/[\s-]/g, '');
  
  // 必须是16-19位数字
  if (!/^\d{16,19}$/.test(cleaned)) {
    return false;
  }
  
  // Luhn算法验证
  let sum = 0;
  let isEven = false;
  
  // 从右往左处理
  for (let i = cleaned.length - 1; i >= 0; i--) {
    let digit = parseInt(cleaned[i]);
    
    if (isEven) {
      digit *= 2;
      if (digit > 9) {
        digit -= 9;
      }
    }
    
    sum += digit;
    isEven = !isEven;
  }
  
  return sum % 10 === 0;
}

// 提取银行卡号
const bankCardMatches = text.match(/\b\d{16,19}\b/g) || [];
const bankCards = [...new Set(bankCardMatches.filter(card => isValidBankCard(card)))];

const result = {
  success: true,
  extracted: {
    emails: emails,
    urls: urls,
    phones: phones,
    ips: ips,
    idCards: idCards,
    bankCards: bankCards,
  },
  message: `提取完成：邮箱 ${emails.length} 个，URL ${urls.length} 个，手机号 ${phones.length} 个，IP ${ips.length} 个，身份证 ${idCards.length} 个，银行卡 ${bankCards.length} 个`,
};

return result;
