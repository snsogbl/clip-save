//go:build darwin

package common

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
#import <dispatch/dispatch.h>

void MoveWindowToCurrentScreen() {
    // 使用同步方式执行，确保在显示窗口之前完成移动
    if ([NSThread isMainThread]) {
        // 如果已经在主线程，直接执行
        NSWindow *mainWindow = [NSApplication sharedApplication].mainWindow;
        if (mainWindow == nil) {
            // 如果没有主窗口，尝试获取所有窗口的第一个
            NSArray *windows = [NSApplication sharedApplication].windows;
            if (windows.count > 0) {
                mainWindow = [windows objectAtIndex:0];
            }
        }

        if (mainWindow != nil) {
            // 获取当前聚焦的屏幕（鼠标所在的屏幕）
            NSScreen *currentScreen = [NSScreen mainScreen];

            // 获取所有屏幕
            NSArray *screens = [NSScreen screens];
            if (screens.count > 1) {
                // 如果有多个屏幕，找到包含鼠标光标的屏幕
                NSPoint mouseLocation = [NSEvent mouseLocation];
                for (NSScreen *screen in screens) {
                    NSRect screenFrame = [screen frame];
                    if (NSPointInRect(mouseLocation, screenFrame)) {
                        currentScreen = screen;
                        break;
                    }
                }
            }

            // 获取目标屏幕的框架
            NSRect screenFrame = [currentScreen frame];

            // 获取窗口的框架
            NSRect windowFrame = [mainWindow frame];

            // 计算窗口在目标屏幕中心的坐标
            CGFloat centerX = screenFrame.origin.x + (screenFrame.size.width - windowFrame.size.width) / 2;
            CGFloat centerY = screenFrame.origin.y + (screenFrame.size.height - windowFrame.size.height) / 2;

            // 设置窗口位置
            [mainWindow setFrameOrigin:NSMakePoint(centerX, centerY)];
        }
    } else {
        // 如果不在主线程，同步调度到主线程执行
        dispatch_sync(dispatch_get_main_queue(), ^{
            NSWindow *mainWindow = [NSApplication sharedApplication].mainWindow;
            if (mainWindow == nil) {
                NSArray *windows = [NSApplication sharedApplication].windows;
                if (windows.count > 0) {
                    mainWindow = [windows objectAtIndex:0];
                }
            }

            if (mainWindow != nil) {
                NSScreen *currentScreen = [NSScreen mainScreen];
                NSArray *screens = [NSScreen screens];
                if (screens.count > 1) {
                    NSPoint mouseLocation = [NSEvent mouseLocation];
                    for (NSScreen *screen in screens) {
                        NSRect screenFrame = [screen frame];
                        if (NSPointInRect(mouseLocation, screenFrame)) {
                            currentScreen = screen;
                            break;
                        }
                    }
                }

                NSRect screenFrame = [currentScreen frame];
                NSRect windowFrame = [mainWindow frame];
                CGFloat centerX = screenFrame.origin.x + (screenFrame.size.width - windowFrame.size.width) / 2;
                CGFloat centerY = screenFrame.origin.y + (screenFrame.size.height - windowFrame.size.height) / 2;
                [mainWindow setFrameOrigin:NSMakePoint(centerX, centerY)];
            }
        });
    }
}
*/
import "C"
import (
	"context"
)

// MoveWindowToCurrentScreen 将窗口移动到当前聚焦的屏幕中心（仅 macOS）
// 使用同步方式执行，确保在显示窗口之前完成移动，避免闪烁
func MoveWindowToCurrentScreen(ctx context.Context) {
	// 使用 Cocoa API 同步移动窗口
	C.MoveWindowToCurrentScreen()
}
