//go:build darwin
// +build darwin

package common

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
#import <stdlib.h>

// 获取剪贴板中所有可用的类型
char** getPasteboardTypes(int* count) {
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    NSArray *types = [pasteboard types];

    *count = (int)[types count];
    if (*count == 0) {
        return NULL;
    }

    char** result = (char**)malloc(sizeof(char*) * (*count));
    for (int i = 0; i < *count; i++) {
        NSString *type = [types objectAtIndex:i];
        result[i] = strdup([type UTF8String]);
    }

    return result;
}

// 释放类型数组
void freeTypes(char** types, int count) {
    if (types == NULL) return;
    for (int i = 0; i < count; i++) {
        free(types[i]);
    }
    free(types);
}

// 读取指定类型的数据
unsigned char* readPasteboardData(const char* type, int* length) {
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    NSString *typeString = [NSString stringWithUTF8String:type];
    NSData *data = [pasteboard dataForType:typeString];

    if (data == nil) {
        *length = 0;
        return NULL;
    }

    *length = (int)[data length];
    unsigned char* result = (unsigned char*)malloc(*length);
    [data getBytes:result length:*length];

    return result;
}

// 释放数据
void freeData(unsigned char* data) {
    if (data != NULL) {
        free(data);
    }
}

// 读取文件 URL 列表
char* readFileURLs(int* count) {
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];

    // 尝试读取文件 URL（新版）
    NSArray *urls = [pasteboard readObjectsForClasses:@[[NSURL class]] options:@{NSPasteboardURLReadingFileURLsOnlyKey: @YES}];

    if (urls == nil || [urls count] == 0) {
        // 尝试旧版 NSFilenamesPboardType
        NSArray *filenames = [pasteboard propertyListForType:NSFilenamesPboardType];
        if (filenames != nil && [filenames count] > 0) {
            // 转换为 JSON 数组格式
            NSError *error = nil;
            NSData *jsonData = [NSJSONSerialization dataWithJSONObject:filenames options:0 error:&error];
            if (jsonData) {
                NSString *jsonString = [[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];
                *count = (int)[filenames count];
                return strdup([jsonString UTF8String]);
            }
        }
        *count = 0;
        return NULL;
    }

    // 转换 URL 数组为路径数组
    NSMutableArray *paths = [NSMutableArray arrayWithCapacity:[urls count]];
    for (NSURL *url in urls) {
        if ([url isFileURL]) {
            [paths addObject:[url path]];
        }
    }

    if ([paths count] == 0) {
        *count = 0;
        return NULL;
    }

    // 转换为 JSON 字符串
    NSError *error = nil;
    NSData *jsonData = [NSJSONSerialization dataWithJSONObject:paths options:0 error:&error];
    if (jsonData == nil) {
        *count = 0;
        return NULL;
    }

    NSString *jsonString = [[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];
    *count = (int)[paths count];

    return strdup([jsonString UTF8String]);
}

// 释放字符串
void freeString(char* str) {
    if (str != NULL) {
        free(str);
    }
}

// 获取剪贴板的 changeCount
int getPasteboardChangeCount() {
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    return (int)[pasteboard changeCount];
}

// 写入文件 URL 到剪贴板
int writeFileURLs(const char* jsonPaths) {
    if (jsonPaths == NULL) {
        return 0;
    }

    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];

    // 解析 JSON 数组
    NSString *jsonString = [NSString stringWithUTF8String:jsonPaths];
    NSData *jsonData = [jsonString dataUsingEncoding:NSUTF8StringEncoding];

    NSError *error = nil;
    NSArray *paths = [NSJSONSerialization JSONObjectWithData:jsonData options:0 error:&error];

    if (error != nil || paths == nil || ![paths isKindOfClass:[NSArray class]]) {
        return 0;
    }

    // 转换路径为 NSURL 数组
    NSMutableArray *urls = [NSMutableArray arrayWithCapacity:[paths count]];
    for (id pathObj in paths) {
        if ([pathObj isKindOfClass:[NSString class]]) {
            NSString *path = (NSString *)pathObj;
            NSURL *url = [NSURL fileURLWithPath:path];
            if (url != nil) {
                [urls addObject:url];
            }
        }
    }

    if ([urls count] == 0) {
        return 0;
    }

    // 清空剪贴板并写入文件 URL
    [pasteboard clearContents];
    BOOL success = [pasteboard writeObjects:urls];

    return success ? 1 : 0;
}

// 获取当前活动应用程序的名称
char* getFrontmostAppName() {
    NSWorkspace *workspace = [NSWorkspace sharedWorkspace];
    if (workspace == nil) {
        return strdup("Unknown");
    }
    
    NSRunningApplication *frontApp = [workspace frontmostApplication];
    if (frontApp == nil) {
        return strdup("Unknown");
    }

    NSString *appName = [frontApp localizedName];
    if (appName == nil || [appName length] == 0) {
        appName = [frontApp bundleIdentifier];
        if (appName == nil || [appName length] == 0) {
            return strdup("Unknown");
        }
    }

    return strdup([appName UTF8String]);
}

*/
import "C"
import (
	"fmt"
	"unsafe"
)

// GetPasteboardTypes 获取剪贴板中所有可用的类型
func GetPasteboardTypes() []string {
	var count C.int
	cTypes := C.getPasteboardTypes(&count)
	if cTypes == nil {
		return []string{}
	}
	defer C.freeTypes(cTypes, count)

	types := make([]string, int(count))
	typesSlice := (*[1 << 30]*C.char)(unsafe.Pointer(cTypes))[:count:count]

	for i := 0; i < int(count); i++ {
		types[i] = C.GoString(typesSlice[i])
	}

	return types
}

// ReadPasteboardData 读取指定类型的剪贴板数据
func ReadPasteboardData(typeName string) []byte {
	cType := C.CString(typeName)
	defer C.free(unsafe.Pointer(cType))

	var length C.int
	cData := C.readPasteboardData(cType, &length)
	if cData == nil {
		return nil
	}
	defer C.freeData(cData)

	if length == 0 {
		return nil
	}

	// 复制数据到 Go slice
	data := make([]byte, int(length))
	dataSlice := (*[1 << 30]byte)(unsafe.Pointer(cData))[:length:length]
	copy(data, dataSlice)

	return data
}

// ReadFileURLs 读取剪贴板中的文件 URL 列表
// 返回 JSON 格式的文件路径数组和文件数量
func ReadFileURLs() (string, int) {
	var count C.int
	cJSON := C.readFileURLs(&count)
	if cJSON == nil {
		return "", 0
	}
	defer C.freeString(cJSON)

	jsonString := C.GoString(cJSON)
	return jsonString, int(count)
}

// WriteFileURLs 将文件 URL 写入剪贴板
// filePaths 是 JSON 格式的文件路径数组
func WriteFileURLs(filePaths string) error {
	cJSON := C.CString(filePaths)
	defer C.free(unsafe.Pointer(cJSON))

	result := C.writeFileURLs(cJSON)
	if result == 0 {
		return fmt.Errorf("写入文件 URL 失败")
	}
	return nil
}

// GetFrontmostAppName 获取当前活动应用程序的名称
func GetFrontmostAppName() string {
	cAppName := C.getFrontmostAppName()
	if cAppName == nil {
		return "Unknown"
	}
	defer C.freeString(cAppName)

	return C.GoString(cAppName)
}

// GetPasteboardChangeCount 获取剪贴板的变化计数
func GetPasteboardChangeCount() int {
	return int(C.getPasteboardChangeCount())
}
