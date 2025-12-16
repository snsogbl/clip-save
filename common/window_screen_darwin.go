//go:build darwin

package common

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
#import <dispatch/dispatch.h>

void EnsureWindowOnCurrentScreen() {
    // 使用同步方式执行，确保在显示窗口之前完成移动
    if ([NSThread isMainThread]) {
        // 检查 NSApplication 是否已初始化
        NSApplication *app = [NSApplication sharedApplication];
        if (app == nil) {
            return;
        }

        // 如果已经在主线程，直接执行
        NSWindow *mainWindow = app.mainWindow;
        if (mainWindow == nil) {
            // 如果没有主窗口，尝试获取所有窗口的第一个
            NSArray *windows = app.windows;
            if (windows != nil && windows.count > 0) {
                mainWindow = [windows objectAtIndex:0];
            }
        }

        if (mainWindow != nil) {
            // 设置窗口的集合行为，使其能够移动到当前活动的桌面空间
            NSWindowCollectionBehavior behavior = [mainWindow collectionBehavior];
            // 添加 NSWindowCollectionBehaviorMoveToActiveSpace，让窗口跟随当前活动的空间
            behavior |= NSWindowCollectionBehaviorMoveToActiveSpace;
            [mainWindow setCollectionBehavior:behavior];

            // 如果窗口被隐藏，强制移动到主屏幕中心
            // 这样窗口会出现在当前活动的桌面空间
            BOOL windowVisible = [mainWindow isVisible];

            if (!windowVisible) {
                // 窗口被隐藏，移动到主屏幕中心
                NSScreen *mainScreen = [NSScreen mainScreen];
                NSRect screenFrame = [mainScreen frame];
                NSRect windowFrame = [mainWindow frame];

                CGFloat centerX = screenFrame.origin.x + (screenFrame.size.width - windowFrame.size.width) / 2;
                CGFloat centerY = screenFrame.origin.y + (screenFrame.size.height - windowFrame.size.height) / 2;
                [mainWindow setFrameOrigin:NSMakePoint(centerX, centerY)];
            } else {
                // 窗口可见，检查是否在当前屏幕上
                NSScreen *currentScreen = nil;

                // 优先使用包含键盘焦点的窗口所在的屏幕
                NSWindow *keyWindow = app.keyWindow;
                if (keyWindow != nil && keyWindow != mainWindow) {
                    currentScreen = [keyWindow screen];
                }

                // 如果找不到，使用包含鼠标光标的屏幕
                if (currentScreen == nil) {
                    NSPoint mouseLocation = [NSEvent mouseLocation];
                    NSArray *screens = [NSScreen screens];
                    for (NSScreen *screen in screens) {
                        NSRect screenFrame = [screen frame];
                        if (NSPointInRect(mouseLocation, screenFrame)) {
                            currentScreen = screen;
                            break;
                        }
                    }
                }

                // 如果还是找不到，使用主屏幕
                if (currentScreen == nil) {
                    currentScreen = [NSScreen mainScreen];
                }

                NSRect windowFrame = [mainWindow frame];
                NSPoint windowCenter = NSMakePoint(
                    windowFrame.origin.x + windowFrame.size.width / 2,
                    windowFrame.origin.y + windowFrame.size.height / 2
                );
                NSRect screenFrame = [currentScreen frame];
                BOOL windowOnScreen = NSPointInRect(windowCenter, screenFrame);

                if (!windowOnScreen) {
                    // 窗口不在当前屏幕上，移动到当前屏幕中心
                    CGFloat centerX = screenFrame.origin.x + (screenFrame.size.width - windowFrame.size.width) / 2;
                    CGFloat centerY = screenFrame.origin.y + (screenFrame.size.height - windowFrame.size.height) / 2;
                    [mainWindow setFrameOrigin:NSMakePoint(centerX, centerY)];
                }
            }
        }
    } else {
        // 如果不在主线程，同步调度到主线程执行
        dispatch_sync(dispatch_get_main_queue(), ^{
            // 检查 NSApplication 是否已初始化
            NSApplication *app = [NSApplication sharedApplication];
            if (app == nil) {
                return;
            }

            NSWindow *mainWindow = app.mainWindow;
            if (mainWindow == nil) {
                NSArray *windows = app.windows;
                if (windows != nil && windows.count > 0) {
                    mainWindow = [windows objectAtIndex:0];
                }
            }

            if (mainWindow != nil) {
                // 设置窗口的集合行为，使其能够移动到当前活动的桌面空间
                NSWindowCollectionBehavior behavior = [mainWindow collectionBehavior];
                behavior |= NSWindowCollectionBehaviorMoveToActiveSpace;
                [mainWindow setCollectionBehavior:behavior];

                BOOL windowVisible = [mainWindow isVisible];

                if (!windowVisible) {
                    // 窗口被隐藏，移动到主屏幕中心
                    NSScreen *mainScreen = [NSScreen mainScreen];
                    NSRect screenFrame = [mainScreen frame];
                    NSRect windowFrame = [mainWindow frame];

                    CGFloat centerX = screenFrame.origin.x + (screenFrame.size.width - windowFrame.size.width) / 2;
                    CGFloat centerY = screenFrame.origin.y + (screenFrame.size.height - windowFrame.size.height) / 2;
                    [mainWindow setFrameOrigin:NSMakePoint(centerX, centerY)];
                } else {
                    NSScreen *currentScreen = nil;

                    NSWindow *keyWindow = app.keyWindow;
                    if (keyWindow != nil && keyWindow != mainWindow) {
                        currentScreen = [keyWindow screen];
                    }

                    if (currentScreen == nil) {
                        NSPoint mouseLocation = [NSEvent mouseLocation];
                        NSArray *screens = [NSScreen screens];
                        for (NSScreen *screen in screens) {
                            NSRect screenFrame = [screen frame];
                            if (NSPointInRect(mouseLocation, screenFrame)) {
                                currentScreen = screen;
                                break;
                            }
                        }
                    }

                    if (currentScreen == nil) {
                        currentScreen = [NSScreen mainScreen];
                    }

                    NSRect windowFrame = [mainWindow frame];
                    NSPoint windowCenter = NSMakePoint(
                        windowFrame.origin.x + windowFrame.size.width / 2,
                        windowFrame.origin.y + windowFrame.size.height / 2
                    );
                    NSRect screenFrame = [currentScreen frame];
                    BOOL windowOnScreen = NSPointInRect(windowCenter, screenFrame);

                    if (!windowOnScreen) {
                        CGFloat centerX = screenFrame.origin.x + (screenFrame.size.width - windowFrame.size.width) / 2;
                        CGFloat centerY = screenFrame.origin.y + (screenFrame.size.height - windowFrame.size.height) / 2;
                        [mainWindow setFrameOrigin:NSMakePoint(centerX, centerY)];
                    }
                }
            }
        });
    }
}
*/
import "C"
import "context"

// EnsureWindowOnCurrentScreen 确保窗口在当前聚焦的屏幕上（仅 macOS）
// 如果窗口已经在当前屏幕上，保持原位置不变
// 如果窗口不在当前屏幕上，移动到当前屏幕中心
func EnsureWindowOnCurrentScreen(ctx context.Context) {
	// 使用 Cocoa API 检查并移动窗口
	C.EnsureWindowOnCurrentScreen()
}
