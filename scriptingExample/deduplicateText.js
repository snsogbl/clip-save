/**
 * 统计文本中每行内容的出现次数
 * @author github.com/qyuja
 * @param {string} input - 输入的字符串
 * @returns {string} - 返回重复出现次数 eg: hello world (2) foo bar (3)
 * @returns {object} - 错误信息
 */


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
try{
    const counter = new Map();
    // 过滤空行
    const lines = text.split(/\r?\n/).map(i => i.trim()).filter(Boolean)
    // 去重并计数 
    for (const line of lines) {
      counter.set(line, (counter.get(line) ?? 0) + 1);
    }
    // 输出示例：xxx (8)
    const result = [...counter.entries()]
      .map(([value, count]) => `${value} (${count})`)
      .join("\n");
    return result;
}catch(error){
    return {
        error:"解析异常"
    }
}
