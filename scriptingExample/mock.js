/**
 * JSON 数据 Mock 生成器
 * 根据 JSON Schema 或普通 JSON 对象生成模拟数据
 * @author ClipSave
 * @param {string} input - JSON Schema 或普通 JSON 对象，从剪贴板内容获取
 * @returns {object} - 生成的 Mock 数据
 * @returns {object} - 错误信息
 */

const minArrayLength = 10;
const maxArrayLength = 20;

if (item.ContentType !== "JSON") {
  return {
    error: "只支持 JSON 或文本类型的剪贴板内容",
  };
}

const input = item.Content || "";

if (!input.trim()) {
  return {
    error: "剪贴板内容为空",
  };
}

try {
  // 尝试解析输入的 JSON
  let data;
  try {
    data = JSON.parse(input);
  } catch (e) {
    return {
      error: `无效的 JSON 格式: ${e.message}`,
    };
  }

  // 判断是否为 JSON Schema 格式（包含 type 字段）
  const isSchema = data && typeof data === "object" && data.type !== undefined;

  // 检测是否为 API 响应模板格式（值包含类型占位符如 "string", 0, null 等）
  function isApiTemplate(obj) {
    if (!obj || typeof obj !== "object" || Array.isArray(obj)) {
      return false;
    }
    // 检查是否包含常见的类型占位符
    const typePlaceholders = ["string", "number", "boolean", "object", "array"];
    let placeholderCount = 0;
    let totalFields = 0;
    
    function checkValue(value) {
      if (typeof value === "string" && typePlaceholders.includes(value.toLowerCase())) {
        placeholderCount++;
        return true;
      }
      if (value === 0 || value === null || value === false) {
        // 这些可能是占位符，也可能是真实值，需要结合上下文判断
        return false;
      }
      if (typeof value === "object" && value !== null) {
        if (Array.isArray(value)) {
          if (value.length === 1 && typeof value[0] === "string" && typePlaceholders.includes(value[0].toLowerCase())) {
            placeholderCount++;
            return true;
          }
          return false;
        }
        // 递归检查对象
        for (const val of Object.values(value)) {
          checkValue(val);
        }
      }
      return false;
    }
    
    for (const value of Object.values(obj)) {
      totalFields++;
      checkValue(value);
    }
    
    // 如果超过30%的字段是类型占位符，认为是模板格式
    return totalFields > 0 && placeholderCount / totalFields > 0.3;
  }
  
  const isTemplate = !isSchema && isApiTemplate(data);

  // 生成 UUID
  function generateUUID() {
    return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (c) => {
      const r = (Math.random() * 16) | 0;
      const v = c === "x" ? r : (r & 0x3) | 0x8;
      return v.toString(16);
    });
  }

  // 省市区数据
  const provinces = [
    "北京市",
    "上海市",
    "广东省",
    "江苏省",
    "浙江省",
    "江西省",
    "山东省",
    "四川省",
    "湖北省",
    "河南省",
    "湖南省",
  ];

  const cities = {
    北京市: [
      "东城区",
      "西城区",
      "朝阳区",
      "丰台区",
      "石景山区",
      "海淀区",
      "门头沟区",
      "房山区",
    ],
    上海市: [
      "黄浦区",
      "徐汇区",
      "长宁区",
      "静安区",
      "普陀区",
      "虹口区",
      "杨浦区",
      "浦东新区",
    ],
    广东省: [
      "广州市",
      "深圳市",
      "珠海市",
      "汕头市",
      "佛山市",
      "韶关市",
      "湛江市",
      "肇庆市",
    ],
    江苏省: [
      "南京市",
      "苏州市",
      "无锡市",
      "常州市",
      "镇江市",
      "扬州市",
      "泰州市",
      "南通市",
    ],
    浙江省: [
      "杭州市",
      "宁波市",
      "温州市",
      "嘉兴市",
      "湖州市",
      "绍兴市",
      "金华市",
      "衢州市",
    ],
    山东省: [
      "济南市",
      "青岛市",
      "淄博市",
      "枣庄市",
      "东营市",
      "烟台市",
      "潍坊市",
      "济宁市",
    ],
    四川省: [
      "成都市",
      "自贡市",
      "攀枝花市",
      "泸州市",
      "德阳市",
      "绵阳市",
      "广元市",
      "遂宁市",
    ],
    湖北省: [
      "武汉市",
      "黄石市",
      "十堰市",
      "宜昌市",
      "襄阳市",
      "鄂州市",
      "荆门市",
      "孝感市",
    ],
    河南省: [
      "郑州市",
      "开封市",
      "洛阳市",
      "平顶山市",
      "安阳市",
      "鹤壁市",
      "新乡市",
      "焦作市",
    ],
    湖南省: [
      "长沙市",
      "株洲市",
      "湘潭市",
      "衡阳市",
      "邵阳市",
      "岳阳市",
      "常德市",
      "张家界市",
    ],
  };

  const districts = [
    "解放路街道",
    "建设路街道",
    "人民路街道",
    "中山路街道",
    "胜利路街道",
    "和平路街道",
    "新华路街道",
    "文化路街道",
    "科技路街道",
    "工业路街道",
  ];

  // 中文姓名数据
  const surnames = ["张", "李", "王", "刘", "陈", "杨", "赵", "黄", "周", "吴"];
  const givenNames = ["伟", "芳", "娜", "秀英", "敏", "静", "丽", "强", "磊", "军", "洋", "勇", "艳", "杰", "涛", "明", "超", "秀兰"];

  // 英文姓名数据
  const firstNames = ["James", "John", "Robert", "Michael", "William", "David", "Richard", "Joseph", "Thomas", "Charles"];
  const lastNames = ["Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez"];

  // 生成中文姓名
  function generateChineseName() {
    const surname = surnames[Math.floor(Math.random() * surnames.length)];
    const givenName = givenNames[Math.floor(Math.random() * givenNames.length)];
    return surname + givenName;
  }

  // 生成英文姓名
  function generateEnglishName() {
    const firstName = firstNames[Math.floor(Math.random() * firstNames.length)];
    const lastName = lastNames[Math.floor(Math.random() * lastNames.length)];
    return `${firstName} ${lastName}`;
  }

  // 生成手机号
  function generatePhone() {
    const prefixes = ["130", "131", "132", "133", "134", "135", "136", "137", "138", "139", "150", "151", "152", "153", "155", "156", "157", "158", "159", "180", "181", "182", "183", "184", "185", "186", "187", "188", "189"];
    const prefix = prefixes[Math.floor(Math.random() * prefixes.length)];
    const suffix = Math.floor(Math.random() * 100000000).toString().padStart(8, "0");
    return prefix + suffix;
  }

  // 生成身份证号（18位，简化版，仅用于 mock）
  function generateIdCard() {
    const areas = ["110", "120", "130", "140", "150", "210", "220", "230", "310", "320", "330", "340", "350", "360", "370", "410", "420", "430", "440", "450", "460", "500", "510", "520", "530"];
    const area = areas[Math.floor(Math.random() * areas.length)];
    const year = Math.floor(Math.random() * 30) + 1970;
    const month = Math.floor(Math.random() * 12) + 1;
    const day = Math.floor(Math.random() * 28) + 1;
    const random = Math.floor(Math.random() * 1000).toString().padStart(3, "0");
    const checkCode = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "X"][Math.floor(Math.random() * 11)];
    return `${area}${year}${month.toString().padStart(2, "0")}${day.toString().padStart(2, "0")}${random}${checkCode}`;
  }

  // 生成银行卡号（16-19位，简化版）
  function generateBankCard() {
    const length = [16, 17, 18, 19][Math.floor(Math.random() * 4)];
    let card = "";
    for (let i = 0; i < length; i++) {
      card += Math.floor(Math.random() * 10).toString();
    }
    return card;
  }

  // 生成 IP 地址
  function generateIP() {
    const octets = Array.from({ length: 4 }, () =>
      Math.floor(Math.random() * 255)
    );
    return octets.join(".");
  }

  // 预编译正则表达式（提高性能）
  const regexes = {
    email: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
    phone: /^1[3-9]\d{9}$/,
    idCard: /^\d{15}$|^\d{17}[\dXx]$/,
    bankCard: /^\d{16,19}$/,
    ip: /^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$/,
    timestamp: /^\d{10,13}$/,
    date: /^\d{4}-\d{2}-\d{2}$/, // YYYY-MM-DD
    dateTime: /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/, // YYYY-MM-DD HH:mm:ss
    dateSlash: /^\d{4}\/\d{2}\/\d{2}/,
    iso8601: /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/, // ISO 8601 格式
    chinese: /[\u4e00-\u9fa5]/,
    imageExt: /\.(jpg|jpeg|png|gif|webp|svg|bmp|ico)(\?|$)/i,
    imageKeyword: /image|img|photo|picture|pic/i,
  };

  // 辅助函数：检查字段名是否包含关键词（使用单词边界匹配，避免误判）
  function fieldMatches(fieldName, keywords) {
    const lowerFieldName = (fieldName || "").toLowerCase();
    return keywords.some(keyword => {
      const lowerKeyword = keyword.toLowerCase();
      // 完全匹配
      if (lowerFieldName === lowerKeyword) {
        return true;
      }
      // 以关键词开头（后面跟非字母数字字符或结尾）
      if (lowerFieldName.startsWith(lowerKeyword)) {
        const nextChar = lowerFieldName[lowerKeyword.length];
        if (!nextChar || !/[a-z0-9]/.test(nextChar)) {
          return true;
        }
      }
      // 以关键词结尾（前面跟非字母数字字符或开头）
      if (lowerFieldName.endsWith(lowerKeyword)) {
        const prevChar = lowerFieldName[lowerFieldName.length - lowerKeyword.length - 1];
        if (prevChar === undefined || !/[a-z0-9]/.test(prevChar)) {
          return true;
        }
      }
      // 关键词前后都有分隔符（下划线、驼峰等）
      const index = lowerFieldName.indexOf(lowerKeyword);
      if (index > 0 && index < lowerFieldName.length - lowerKeyword.length) {
        const before = lowerFieldName[index - 1];
        const after = lowerFieldName[index + lowerKeyword.length];
        // 前后都必须是分隔符（非字母数字字符或下划线），不能是字母数字
        if (!/[a-z0-9]/.test(before) && (!after || !/[a-z0-9]/.test(after))) {
          return true;
        }
      }
      return false;
    });
  }

  // 辅助函数：检查字段名是否为图片字段（避免误判 imgSize、imageWidth 等）
  function isImageField(fieldName) {
    if (!fieldName) return false;
    const lowerFieldName = fieldName.toLowerCase();
    
    // 图片相关关键词
    const imageKeywords = ["image", "img", "photo", "picture", "pic", "avatar", "icon", "logo", "thumbnail"];
    
    // 排除的后缀（表示尺寸、宽度、高度等属性，这些不是图片URL）
    const excludeSuffixes = ["size", "width", "height", "w", "h", "length", "count", "num", "id"];
    
    // 检查是否以图片关键词开头或结尾，或者是图片关键词本身
    for (const keyword of imageKeywords) {
      // 完全匹配
      if (lowerFieldName === keyword) {
        return true;
      }
      
      // 以关键词开头
      if (lowerFieldName.startsWith(keyword)) {
        const rest = lowerFieldName.substring(keyword.length);
        
        // 如果后面为空，直接匹配（如 "image"）
        if (rest === "") {
          return true;
        }
        
        // 检查剩余部分是否以排除的后缀开头（如 "imgSize" -> "size"）
        const startsWithExclude = excludeSuffixes.some(suffix => 
          rest.startsWith(suffix) || rest.startsWith("_" + suffix)
        );
        
        // 如果以排除的后缀开头，则不是图片字段
        if (startsWithExclude) {
          continue;
        }
        
        // 其他情况（如 "imageUrl", "imgPath", "userImage" 等）认为是图片字段
        return true;
      }
      
      // 以关键词结尾（如 "userImage", "productImg"）
      if (lowerFieldName.endsWith(keyword)) {
        const prefix = lowerFieldName.substring(0, lowerFieldName.length - keyword.length);
        // 检查前缀是否包含排除的后缀（如 "imgSize" 中的 "size"）
        const hasExcludePrefix = excludeSuffixes.some(suffix => 
          prefix.endsWith(suffix) || prefix.endsWith("_" + suffix)
        );
        if (!hasExcludePrefix) {
          return true;
        }
      }
    }
    
    return false;
  }

  // 辅助函数：检查值是否匹配模式
  function valueMatches(value, patterns) {
    const lowerValue = (value || "").toLowerCase();
    return patterns.some(pattern => {
      if (typeof pattern === "string") {
        return lowerValue.includes(pattern);
      }
      if (pattern instanceof RegExp) {
        return pattern.test(value);
      }
      return false;
    });
  }

  // 根据字段名和值生成特殊类型的 mock 数据
  function generateSpecialType(fieldName, value) {
    const lowerFieldName = (fieldName || "").toLowerCase();
    const lowerValue = (value || "").toLowerCase();

    // 识别省份
    if (
      fieldMatches(fieldName, ["省", "province"]) ||
      valueMatches(value, ["省", /省$/])
    ) {
      return provinces[Math.floor(Math.random() * provinces.length)];
    }

    // 识别城市
    if (
      fieldMatches(fieldName, ["市", "city"]) ||
      valueMatches(value, ["市", /市$/])
    ) {
      const randomProvince =
        provinces[Math.floor(Math.random() * provinces.length)];
      const provinceCities = cities[randomProvince];
      if (provinceCities && provinceCities.length > 0) {
        return provinceCities[
          Math.floor(Math.random() * provinceCities.length)
        ];
      }
      // 如果没有对应省份的城市数据，生成通用城市名
      const cityNames = ["市", "区", "县"];
      return (
        randomProvince.replace("省", "") +
        cityNames[Math.floor(Math.random() * cityNames.length)]
      );
    }

    // 识别区/县
    if (
      fieldMatches(fieldName, ["区", "县", "district", "county"]) ||
      valueMatches(value, ["区", "县", /[区县]$/])
    ) {
      return districts[Math.floor(Math.random() * districts.length)];
    }

    // 识别姓名
    if (fieldMatches(fieldName, ["name", "姓名", "名字"])) {
      // 如果原值是中文，生成中文名；否则生成英文名
      if (regexes.chinese.test(value)) {
        return generateChineseName();
      }
      return generateEnglishName();
    }

    // 识别手机号
    if (
      fieldMatches(fieldName, ["phone", "mobile", "手机", "电话"]) ||
      regexes.phone.test(value)
    ) {
      return generatePhone();
    }

    // 识别身份证号
    if (
      fieldMatches(fieldName, ["idcard", "id_card", "身份证"]) ||
      regexes.idCard.test(value)
    ) {
      return generateIdCard();
    }

    // 识别银行卡号（使用精确匹配，避免误判 idcard 等字段）
    if (
      fieldMatches(fieldName, ["bankcard", "bank_card", "银行卡", "card"]) ||
      regexes.bankCard.test(value)
    ) {
      // 排除 idcard 相关字段
      if (lowerFieldName.includes("idcard") || lowerFieldName.includes("id_card")) {
        return null;
      }
      return generateBankCard();
    }

    // 识别 IP 地址（使用精确匹配，避免误判 emailAddress 等字段）
    if (
      fieldMatches(fieldName, ["ip", "ipaddress", "ip_address"]) ||
      regexes.ip.test(value)
    ) {
      // 排除 email 相关字段
      if (lowerFieldName.includes("email") || lowerFieldName.includes("mail")) {
        return null;
      }
      return generateIP();
    }

    // 识别价格、金额、费用等
    if (
      fieldMatches(fieldName, [
        "price",
        "amount",
        "cost",
        "fee",
        "money",
        "价格",
        "金额",
        "费用",
      ])
    ) {
      // 生成合理的价格（0.01 - 9999.99）
      return (Math.random() * 9999.98 + 0.01).toFixed(2);
    }

    // 识别数量、计数等
    if (
      fieldMatches(fieldName, ["count", "quantity", "num", "数量", "个数"])
    ) {
      // 生成合理的数量（1-1000）
      return Math.floor(Math.random() * 1000) + 1;
    }

    // 识别年龄（使用更精确的匹配，避免误判 message 等字段）
    if (fieldMatches(fieldName, ["age", "年龄"])) {
      // 生成合理的年龄（18-80）
      return Math.floor(Math.random() * 63) + 18;
    }

    // 识别分数、评分等
    if (fieldMatches(fieldName, ["score", "rating", "分数", "评分"])) {
      // 生成0-100的分数，或0-5的评分
      if (lowerFieldName.includes("rating") || lowerFieldName.includes("评分")) {
        return (Math.random() * 5).toFixed(1);
      }
      return Math.floor(Math.random() * 101);
    }

    // 识别时间字段（date、time、日期、时间等）
    // 注意：这里只根据字段名识别，不根据值，因为值可能是类型占位符
    // 实际格式会在 mockValueByType 中根据原始值格式生成
    if (
      fieldMatches(fieldName, ["time", "date", "时间", "日期", "datetime", "timestamp"]) ||
      lowerFieldName.endsWith("time") ||
      lowerFieldName.endsWith("date")
    ) {
      // 如果没有原始值，默认返回 ISO 8601 格式
      return new Date().toISOString();
    }

    // 识别图片字段（使用精确匹配，避免误判 imgSize、imageWidth 等）
    if (isImageField(fieldName)) {
      const width = Math.floor(Math.random() * 700) + 100;
      const height = Math.floor(Math.random() * 700) + 100;
      return `https://dummyimage.com/${width}x${height}`;
    }

    return null;
  }

  // 根据原始时间格式生成相同格式的随机时间
  function generateTimeByFormat(originalValue) {
    if (!originalValue || typeof originalValue !== "string") {
      return new Date().toISOString();
    }
    
    // 生成随机日期（过去1年到未来1年之间）
    const now = new Date();
    const oneYearAgo = new Date(now.getFullYear() - 1, now.getMonth(), now.getDate());
    const oneYearLater = new Date(now.getFullYear() + 1, now.getMonth(), now.getDate());
    const randomTime = oneYearAgo.getTime() + Math.random() * (oneYearLater.getTime() - oneYearAgo.getTime());
    const randomDate = new Date(randomTime);
    
    const year = randomDate.getFullYear();
    const month = String(randomDate.getMonth() + 1).padStart(2, "0");
    const day = String(randomDate.getDate()).padStart(2, "0");
    
    // 生成随机时间（0-23小时，0-59分钟，0-59秒）
    const hours = String(Math.floor(Math.random() * 24)).padStart(2, "0");
    const minutes = String(Math.floor(Math.random() * 60)).padStart(2, "0");
    const seconds = String(Math.floor(Math.random() * 60)).padStart(2, "0");
    
    // YYYY-MM-DD HH:mm:ss 格式
    if (regexes.dateTime.test(originalValue)) {
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    }
    
    // YYYY-MM-DD 格式
    if (regexes.date.test(originalValue)) {
      return `${year}-${month}-${day}`;
    }
    
    // ISO 8601 格式
    if (regexes.iso8601.test(originalValue)) {
      // 生成随机的 ISO 8601 格式时间
      const randomISO = new Date(randomTime);
      randomISO.setHours(parseInt(hours), parseInt(minutes), parseInt(seconds));
      return randomISO.toISOString();
    }
    
    // 其他日期格式（如 YYYY/MM/DD）
    if (regexes.dateSlash.test(originalValue)) {
      return `${year}/${month}/${day}`;
    }
    
    // 默认返回随机 ISO 8601
    const randomISO = new Date(randomTime);
    randomISO.setHours(parseInt(hours), parseInt(minutes), parseInt(seconds));
    return randomISO.toISOString();
  }

  // 提取公共的数字生成函数
  function generateNumberByValue(value) {
    if (Number.isInteger(value)) {
      // 如果是整数，保持相似的数值范围
      const absValue = Math.abs(value);
      if (absValue === 0) return 0;
      // 处理负数
      const isNegative = value < 0;
      const magnitude = Math.pow(10, Math.floor(Math.log10(absValue)));
      const min = Math.max(1, Math.floor(absValue / magnitude) * magnitude);
      const max = min + magnitude * 9;
      const result = Math.floor(Math.random() * (max - min + 1)) + min;
      return isNegative ? -result : result;
    } else {
      // 浮点数，保持相似的数值范围
      const absValue = Math.abs(value);
      if (absValue === 0) return Math.random();
      // 处理负数
      const isNegative = value < 0;
      const magnitude = Math.pow(10, Math.floor(Math.log10(absValue)));
      const result = (Math.random() * magnitude * 10).toFixed(2);
      return isNegative ? -parseFloat(result) : parseFloat(result);
    }
  }

  // 提取公共的字符串生成函数
  function generateStringByContent(value) {
    // 识别 email
    if (regexes.email.test(value)) {
      return `mock${Math.floor(Math.random() * 10000)}@example.com`;
    }
    // 识别图片链接
    if (
      value.startsWith("http") &&
      (regexes.imageExt.test(value) || regexes.imageKeyword.test(value))
    ) {
      const width = Math.floor(Math.random() * 700) + 100;
      const height = Math.floor(Math.random() * 700) + 100;
      return `https://dummyimage.com/${width}x${height}`;
    }
    // 识别普通 URL
    if (value.startsWith("http")) {
      return `https://example.com/${Math.random().toString(36).substring(7)}`;
    }
    // 识别日期/时间格式（保持原始格式）
    if (regexes.date.test(value) || regexes.dateTime.test(value) || regexes.dateSlash.test(value) || regexes.iso8601.test(value)) {
      return generateTimeByFormat(value);
    }
    // 识别时间戳
    if (regexes.timestamp.test(value)) {
      return Date.now().toString();
    }
    // 默认生成随机字符串
    const targetLength = Math.max(6, Math.min(value.length || 10, 20));
    return Math.random().toString(36).substring(2, 2 + targetLength);
  }

  // 根据值类型生成 Mock 数据（用于普通 JSON 对象）
  function mockValueByType(value, fieldName = null, depth = 0) {
    // 防止无限递归，限制嵌套深度
    if (depth > 10) {
      return null;
    }

    // 处理 undefined
    if (value === undefined) {
      return undefined;
    }

    if (value === null) {
      return null;
    }

    // 处理 NaN
    if (typeof value === "number" && isNaN(value)) {
      return NaN;
    }

    const type = typeof value;

    switch (type) {
      case "string":
        // 检查是否为类型占位符（API 模板格式）
        if (typeof value === "string" && value.toLowerCase() === "string") {
          // 这是类型占位符，根据字段名生成合适的字符串
          const specialType = generateSpecialType(fieldName, "");
          // 确保返回的是字符串类型，如果是数字或其他类型，忽略它
          if (specialType && typeof specialType === "string") {
            return specialType;
          }
          // 默认生成有意义的字符串
          return `mock${Math.floor(Math.random() * 10000)}`;
        }
        // 如果是时间字段，优先检查原始值格式
        if (
          fieldMatches(fieldName, ["time", "date", "时间", "日期", "datetime", "timestamp"]) ||
          (fieldName && (fieldName.toLowerCase().endsWith("time") || fieldName.toLowerCase().endsWith("date")))
        ) {
          // 检查原始值是否是时间格式
          if (regexes.date.test(value) || regexes.dateTime.test(value) || regexes.iso8601.test(value) || regexes.dateSlash.test(value)) {
            return generateTimeByFormat(value);
          }
        }
        
        // 先尝试识别特殊类型（省市区、姓名、手机号等）
        const specialType = generateSpecialType(fieldName, value);
        if (specialType) {
          return specialType;
        }
        // 根据字符串内容推断类型
        return generateStringByContent(value);

      case "number":
        // 检查是否为类型占位符（API 模板格式中，0 可能是占位符）
        if (value === 0 && isTemplate) {
          // 根据字段名生成合适的数字
          if (fieldMatches(fieldName, ["code", "status", "状态码"])) {
            return 0; // 状态码 0 通常表示成功
          }
          if (fieldMatches(fieldName, ["count", "quantity", "数量", "个数"])) {
            return Math.floor(Math.random() * 100) + 1;
          }
          // 默认生成一个合理的数字
          return Math.floor(Math.random() * 1000) + 1;
        }
        // 根据原始值的范围生成更合理的数字
        return generateNumberByValue(value);

      case "boolean":
        return Math.random() > 0.5;

      case "object":
        if (Array.isArray(value)) {
          // 数组：根据第一个元素类型生成
          if (value.length === 0) {
            // 空数组返回空数组，保持数据结构一致性
            return [];
          }
          // 检查是否为类型占位符数组（如 ["string"]）
          const firstItem = value[0];
          if (value.length === 1 && typeof firstItem === "string" && firstItem.toLowerCase() === "string") {
            // 这是类型占位符数组，根据字段名生成合适的数组
            const specialType = generateSpecialType(fieldName, "");
            // 确保返回的是字符串类型，如果是数字或其他类型，忽略它
            if (specialType && typeof specialType === "string") {
              // 如果识别为特殊类型（如时间），使用统一的数组长度范围
              const length = Math.floor(Math.random() * (maxArrayLength - minArrayLength + 1)) + minArrayLength;
              return Array.from({ length }, () => specialType);
            }
            // 默认生成字符串数组，使用统一的数组长度范围
            const length = Math.floor(Math.random() * (maxArrayLength - minArrayLength + 1)) + minArrayLength;
            return Array.from({ length }, () => `mock${Math.floor(Math.random() * 10000)}`);
          }
          // 非空数组：根据第一个元素类型生成
          // 使用统一的数组长度范围
          const length = Math.floor(Math.random() * (maxArrayLength - minArrayLength + 1)) + minArrayLength;
          // 统一处理数组元素生成，传递字段名以便识别特殊类型（如时间字段）
          return Array.from({ length }, () =>
            mockValueByType(firstItem, fieldName, depth + 1)
          );
        } else {
          // 对象：递归处理每个属性
          const result = {};
          for (const [key, val] of Object.entries(value)) {
            result[key] = mockValueByType(val, key, depth + 1);
          }
          return result;
        }

      default:
        return value;
    }
  }

  // 根据 JSON Schema 生成 Mock 数据
  function generateMockFromSchema(schema, depth = 0, fieldName = null) {
    // 防止无限递归
    if (depth > 10) {
      return null;
    }

    // 如果 schema 是字符串 "string"，这是 API 模板格式，应该生成字符串
    if (typeof schema === "string" && schema.toLowerCase() === "string") {
      const specialType = generateSpecialType(fieldName, "");
      // 确保返回的是字符串类型，如果是数字或其他类型，忽略它
      if (specialType && typeof specialType === "string") {
        return specialType;
      }
      // 默认生成字符串
      return `mock${Math.floor(Math.random() * 10000)}`;
    }

    // 如果 schema 是数字 0，这是 API 模板格式，应该根据字段名生成数字
    if (typeof schema === "number" && schema === 0) {
      if (fieldMatches(fieldName, ["code", "status", "状态码"])) {
        return 0;
      }
      return Math.floor(Math.random() * 1000) + 1;
    }

    // 如果 schema 是其他基本类型，直接返回（不处理）
    if (typeof schema !== "object" || schema === null) {
      // 但如果是布尔值或数字，可能需要生成 mock 数据
      if (typeof schema === "boolean") {
        return Math.random() > 0.5;
      }
      if (typeof schema === "number") {
        // 数字类型占位符，根据字段名生成
        if (fieldMatches(fieldName, ["code", "status", "状态码"])) {
          return 0;
        }
        return Math.floor(Math.random() * 1000) + 1;
      }
      return schema;
    }

    // 如果提供了示例值
    if (schema.example !== undefined) {
      return schema.example;
    }

    // 如果提供了默认值
    if (schema.default !== undefined) {
      return schema.default;
    }

    // 如果提供了枚举值
    if (schema.enum && Array.isArray(schema.enum) && schema.enum.length > 0) {
      return schema.enum[Math.floor(Math.random() * schema.enum.length)];
    }

    const type =
      schema.type || (Array.isArray(schema.type) ? schema.type[0] : "string");

    switch (type) {
      case "string":
        const format = (schema.format || "").toLowerCase();
        // 先检查字段名，识别特殊类型
        if (fieldName) {
          const lowerFieldName = fieldName.toLowerCase();
          // 识别时间字段
          if (
            fieldMatches(fieldName, ["time", "date", "时间", "日期", "datetime", "timestamp"]) ||
            lowerFieldName.endsWith("time") ||
            lowerFieldName.endsWith("date")
          ) {
            return new Date().toISOString();
          }
          // 识别邮箱字段
          if (fieldMatches(fieldName, ["email", "mail", "邮箱"])) {
            return `mock${Math.floor(Math.random() * 10000)}@example.com`;
          }
          // 识别图片字段
          if (isImageField(fieldName)) {
            const width = Math.floor(Math.random() * 700) + 100;
            const height = Math.floor(Math.random() * 700) + 100;
            return `https://dummyimage.com/${width}x${height}`;
          }
        }
        // 根据 format 识别
        if (format.includes("email")) {
          return `mock${Math.floor(Math.random() * 10000)}@example.com`;
        }
        if (format.includes("uri") || format.includes("url")) {
          return `https://example.com/${Math.random().toString(36).substring(7)}`;
        }
        if (format.includes("date") || format.includes("time")) {
          return new Date().toISOString();
        }
        if (format.includes("uuid")) {
          return generateUUID();
        }
        const minLength = schema.minLength || 0;
        const maxLength = schema.maxLength || 10;
        const length =
          Math.floor(Math.random() * (maxLength - minLength + 1)) + minLength;
        return Math.random().toString(36).substring(2, 2 + length);

      case "integer":
      case "number":
        const min = schema.minimum !== undefined ? schema.minimum : 0;
        const max = schema.maximum !== undefined ? schema.maximum : 100;
        if (type === "integer") {
          return Math.floor(Math.random() * (max - min + 1)) + min;
        }
        return Math.random() * (max - min) + min;

      case "boolean":
        return Math.random() > 0.5;

      case "array":
        // 如果没有指定 minItems 和 maxItems，使用统一的数组长度范围
        const minItems = schema.minItems !== undefined ? schema.minItems : minArrayLength;
        const maxItems = schema.maxItems !== undefined ? schema.maxItems : maxArrayLength;
        const arrayLength =
          Math.floor(Math.random() * (maxItems - minItems + 1)) + minItems;
        const items = schema.items || { type: "string" };
        return Array.from({ length: arrayLength }, () =>
          generateMockFromSchema(items, depth + 1, fieldName)
        );

      case "object":
        const properties = schema.properties || {};
        const result = {};
        for (const [key, value] of Object.entries(properties)) {
          result[key] = generateMockFromSchema(value, depth + 1, key);
        }
        return result;

      default:
        return null;
    }
  }

  // 生成 Mock 数据
  let mockData;
  if (isSchema) {
    // JSON Schema 格式
    mockData = generateMockFromSchema(data, 0, null);
  } else {
    // 普通 JSON 对象
    mockData = mockValueByType(data);
  }

  // 格式化输出
  return mockData;
} catch (error) {
  return {
    error: `生成 Mock 数据失败: ${error.message || String(error)}`,
  };
}
