// 示例脚本JavaScript代码
// item.Content - 剪贴板内容
// item.ContentType - 内容类型 (Text, Image, File, URL, Color, JSON)
// item.Timestamp - 时间戳
// item.Source - 来源应用


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
