/**
 * 脚本执行器 - 在浏览器环境中执行用户脚本
 */

import { EventsOn } from '../../wailsjs/runtime/runtime'
import { GetEnabledUserScriptsByTrigger, GetClipboardItemByID, HttpRequest, GetUserScriptsByIDs, CopyTextToClipboard } from '../../wailsjs/go/main/App'
import { common } from '../../wailsjs/go/models'
import { ElMessageBox } from 'element-plus'

// 使用 Wails 生成的类型
export type ClipboardItem = common.ClipboardItem
export type UserScript = common.UserScript

export interface ScriptResult {
  error?: string
  returnValue?: any  // 脚本的返回值
}

/**
 * 检查脚本是否应该触发
 * 导出供外部使用（如脚本选择器）
 */
export function shouldTriggerScript(script: UserScript, item: ClipboardItem): boolean {
  if (!script.Enabled) {
    return false
  }

  // 检查内容类型
  if (script.ContentType && script.ContentType.length > 0) {
    if (!script.ContentType.includes(item.ContentType)) {
      return false
    }
  }

  // 检查关键词（支持正则表达式）
  if (script.Keywords && script.Keywords.length > 0) {
    const content = item.Content
    const hasKeyword = script.Keywords.some(keyword => {
      // 检查是否是正则表达式格式（以 / 开头）
      if (keyword.startsWith('/') && keyword.length > 1) {
        try {
          // 去掉开头的 /
          const regexStr = keyword.slice(1)
          
          // 查找最后一个 / 的位置（用于分割 pattern 和 flags）
          const lastSlashIndex = regexStr.lastIndexOf('/')
          
          let pattern: string
          let flags = ''
          
          if (lastSlashIndex >= 0) {
            // 有 / 分隔符，可能是 /pattern/ 或 /pattern/flags
            pattern = regexStr.slice(0, lastSlashIndex)
            const afterSlash = regexStr.slice(lastSlashIndex + 1)
            
            if (afterSlash.length > 0) {
              // 有标志部分，如 /pattern/i
              flags = afterSlash
            }
            // 如果没有标志部分，pattern 就是 /pattern/ 格式，使用默认 flags
          } else {
            // 没有找到 /，说明格式不对，回退到字符串匹配
            return content.toLowerCase().includes(keyword.toLowerCase())
          }
          
          // 如果 pattern 为空，回退到字符串匹配
          if (!pattern || pattern.length === 0) {
            return content.toLowerCase().includes(keyword.toLowerCase())
          }
          
          const regex = new RegExp(pattern, flags)
          return regex.test(content)
        } catch (e) {
          // 正则表达式无效，回退到字符串匹配
          console.warn(`无效的正则表达式: ${keyword}`, e)
          return content.toLowerCase().includes(keyword.toLowerCase())
        }
      } else {
        // 普通字符串匹配（不区分大小写）
        return content.toLowerCase().includes(keyword.toLowerCase())
      }
    })
    
    if (!hasKeyword) {
      return false
    }
  }

  return true
}

/**
 * 解析脚本中的 ES6 格式导入语句
 * 支持格式：import { httpRequest, copyTextToClipboard } from '@clipsave/api'
 */
function parseImports(scriptCode: string): Set<string> {
  const imports = new Set<string>()
  
  // 解析 ES6 格式的导入
  // import { httpRequest, copyTextToClipboard } from '@clipsave/api'
  const es6ImportRegex = /import\s*\{([^}]+)\}\s*from\s*['"]@clipsave\/api['"]/g
  let match
  while ((match = es6ImportRegex.exec(scriptCode)) !== null) {
    const importsList = match[1].split(',').map(s => s.trim())
    importsList.forEach(imp => imports.add(imp))
  }
  
  return imports
}

/**
 * 移除脚本中的导入语句，返回纯净的脚本代码
 */
function removeImports(scriptCode: string): string {
  // 移除 ES6 格式的导入
  let cleanCode = scriptCode.replace(
    /import\s*\{[^}]+\}\s*from\s*['"]@clipsave\/api['"];?\s*\n?/g,
    ''
  )
  
  return cleanCode
}

/**
 * 在浏览器环境中执行脚本
 * 导出供外部使用（如脚本编辑器测试功能）
 */
export async function executeScriptInBrowser(
  script: UserScript,
  item: ClipboardItem
): Promise<ScriptResult> {
  const result: ScriptResult = {}

  try {
    // 解析脚本中的导入语句
    const imports = parseImports(script.Script)
    
    // 创建脚本执行上下文
    // 将 item 转换为普通对象，确保 JSON.stringify 能正常工作
    const itemData = {
      ID: item.ID,
      Content: item.Content,
      ContentType: item.ContentType,
      ContentHash: item.ContentHash,
      ImageData: item.ImageData,
      FilePaths: item.FilePaths,
      FileInfo: item.FileInfo,
      Timestamp: item.Timestamp ? (typeof item.Timestamp === 'string' ? item.Timestamp : new Date(item.Timestamp).toISOString()) : '',
      Source: item.Source,
      CharCount: item.CharCount,
      WordCount: item.WordCount,
      IsFavorite: item.IsFavorite,
    }
    
    const context: any = {
      item: itemData, // 使用转换后的对象
      // 注入 alert 函数，使用 Element Plus 的消息框
      alert: async (message: string) => {
        await ElMessageBox.alert(message, '提示', {
          confirmButtonText: '确定',
          type: 'info',
        })
      },
    }
    
    // 定义可用的 API 函数映射
    const availableAPIs: Array<{
      name: string
      func: any
    }> = [
      {
        name: 'csRequest',
        func: HttpRequest,
      },
      {
        name: 'csCopyText',
        func: CopyTextToClipboard,
      },
    ]
    
    // 根据导入语句注入函数
    const injectedFunctions: string[] = []
    for (const api of availableAPIs) {
      if (imports.has(api.name)) {
        context[api.name] = api.func
        injectedFunctions.push(api.name)
      }
    }

    // 移除导入语句，生成纯净的脚本代码
    const cleanScript = removeImports(script.Script)

    // 生成函数注入代码（基于 injectedFunctions）
    const functionInjections = injectedFunctions
      .map(
        (name) => `
        const ${name} = typeof __context !== 'undefined' && __context.${name} 
          ? __context.${name} 
          : null;
      `
      )
      .join('\n')

    // 在浏览器环境中执行脚本
    // 使用异步立即执行函数，注入剪贴板项对象和 API 函数
    // 注意：用户的脚本代码会被直接插入，所以如果脚本有 return 语句，会正确返回
    const scriptWithContext = `
      (async function() {
        // 注入剪贴板项对象
        const item = ${JSON.stringify(context.item)};
        
        // 注入 alert 函数（使用 Element Plus 的消息框）
        const alert = async function(message) {
          if (typeof __context !== 'undefined' && __context.alert) {
            await __context.alert(message);
          }
        };
        
        ${functionInjections}
        
        // 用户脚本代码（在真正的浏览器环境中执行）
        // 脚本可以直接使用 return 语句返回结果，支持 async/await
        ${cleanScript}
      })();
    `

    // 使用 Function 构造函数创建执行函数
    // 这样脚本就在真正的浏览器环境中运行，可以访问所有浏览器 API
    // 包装脚本以捕获返回值
    const executeFn = new Function('__context', `
      try {
        const result = ${scriptWithContext};
        return result;
      } catch (error) {
        return { __error: error.message || String(error) };
      }
    `)

    // 执行脚本（带超时控制）并捕获返回值
    const returnValue = await Promise.race([
      new Promise<any>(async (resolve) => {
        try {
          const value = executeFn(context)
          // 如果返回值是 Promise，等待它完成
          if (value && typeof value.then === 'function') {
            const resolvedValue = await value
            resolve(resolvedValue)
          } else {
            resolve(value)
          }
        } catch (error: any) {
          resolve({ __error: error.message || String(error) })
        }
      }),
      new Promise<any>((_, reject) => {
        setTimeout(() => {
          reject(new Error('脚本执行超时（超过10秒）'))
        }, 10000)
      }),
    ])

    // 如果返回值有错误标记，设置错误
    if (returnValue && typeof returnValue === 'object' && returnValue.__error) {
      result.error = returnValue.__error
    } else {
      // 保存返回值
      // 注意：如果脚本没有 return 语句，returnValue 可能是 undefined
      result.returnValue = returnValue
      console.log(`[脚本 ${script.Name}] 返回值:`, returnValue, '类型:', typeof returnValue)
    }
  } catch (error: any) {
    result.error = error.message || String(error)
    console.error(`[脚本 ${script.Name}] 执行失败:`, error)
  }

  return result
}

/**
 * 处理脚本执行事件
 */
async function handleScriptExecution(data: {
  itemId: string
  trigger: string
  scriptIds?: string[]  // 可选的脚本ID列表（后端已匹配）
  item?: ClipboardItem  // 可选的 item 数据（后端已传递，避免重复查询）
}) {
  const { itemId, trigger, scriptIds, item: itemFromEvent } = data

  try {
    // 如果后端已经传递了 item 数据，直接使用；否则查询数据库
    let item: ClipboardItem
    if (itemFromEvent) {
      console.log(`使用后端传递的 item 数据`)
      // 使用 ClipboardItem.createFrom 转换数据
      item = common.ClipboardItem.createFrom(itemFromEvent)
      
      // 如果 item 是图片类型且没有 ImageData，延迟加载
      if (item.ContentType === 'Image' && (!item.ImageData || item.ImageData.length === 0)) {
        console.log(`延迟加载图片数据...`)
        const fullItem = await GetClipboardItemByID(itemId)
        if (fullItem && fullItem.ImageData) {
          item.ImageData = fullItem.ImageData
        }
      }
    } else {
      // 兼容旧逻辑：如果没有传递 item，查询数据库
      console.log(`查询 item 数据...`)
      const queriedItem = await GetClipboardItemByID(itemId)
      if (!queriedItem) {
        console.error(`未找到剪贴板项: ${itemId}`)
        return
      }
      item = queriedItem
    }

    let scriptsToExecute: UserScript[] = []

    // 如果后端已经提供了匹配的脚本ID列表，直接使用
    if (scriptIds && scriptIds.length > 0) {
      console.log(`后端已匹配 ${scriptIds.length} 个脚本，批量获取执行...`)
      
      // 使用批量查询接口获取脚本
      try {
        scriptsToExecute = await GetUserScriptsByIDs(scriptIds)
        
        if (scriptsToExecute.length === 0) {
          console.log(`无法获取匹配的脚本`)
          return
        }
      } catch (error) {
        console.error(`批量获取脚本失败:`, error)
        return
      }
    } else {
      // 兼容旧逻辑：如果没有提供脚本ID，使用原有方式（用于其他trigger类型）
      console.log(`使用原有方式获取 ${trigger} 脚本...`)
      const scripts = await GetEnabledUserScriptsByTrigger(trigger)
      scriptsToExecute = scripts.filter(
        (script) => shouldTriggerScript(script, item)
      )
    }

    if (scriptsToExecute.length === 0) {
      console.log(`没有匹配的 ${trigger} 脚本`)
      return
    }

    console.log(`找到 ${scriptsToExecute.length} 个匹配的脚本，开始执行...`)

    // 执行每个脚本（静默执行，不显示执行中状态）
    for (const script of scriptsToExecute) {
      console.log(`执行脚本: ${script.Name}`)
      await executeScriptInBrowser(script, item)
    }
  } catch (error) {
    console.error('处理脚本执行事件失败:', error)
  }
}

/**
 * 初始化脚本执行器
 */
export function initScriptExecutor() {
  // 监听脚本执行事件
  EventsOn('clipboard.script.execute', handleScriptExecution)
  console.log('✅ 脚本执行器已初始化')
}

