//go:build darwin
// +build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
#import <dispatch/dispatch.h>
#import <math.h>

// 全局变量用于存储观察者和定时器
static id windowObserver = nil;
static id resizeObserver = nil;
static NSTimer *positionTimer = nil;

// 调整按钮位置的辅助函数
static void updateButtonPositions(NSWindow *window) {
    if (window == nil) return;

    // 获取标准窗口按钮
    NSButton *closeButton = [window standardWindowButton:NSWindowCloseButton];
    NSButton *minimizeButton = [window standardWindowButton:NSWindowMiniaturizeButton];
    NSButton *zoomButton = [window standardWindowButton:NSWindowZoomButton];

    // 计算按钮的目标 Y 位置（从窗口顶部计算）
    CGFloat targetY = 3.0; // 距离顶部 3px

    // 调整关闭按钮位置
    if (closeButton != nil) {
        NSRect frame = closeButton.frame;
        if (fabs(frame.origin.y - targetY) > 0.1) { // 只在位置需要调整时才更新
            frame.origin.y = targetY;
            closeButton.frame = frame;
        }
    }

    // 调整最小化按钮位置
    if (minimizeButton != nil) {
        NSRect frame = minimizeButton.frame;
        if (fabs(frame.origin.y - targetY) > 0.1) {
            frame.origin.y = targetY;
            minimizeButton.frame = frame;
        }
    }

    // 调整全屏按钮位置
    if (zoomButton != nil) {
        NSRect frame = zoomButton.frame;
        if (fabs(frame.origin.y - targetY) > 0.1) {
            frame.origin.y = targetY;
            zoomButton.frame = frame;
        }
    }
}

void AdjustWindowButtons() {
    dispatch_async(dispatch_get_main_queue(), ^{
        // 获取主窗口
        NSWindow *window = [NSApplication sharedApplication].mainWindow;
        if (window == nil) {
            // 如果没有主窗口，尝试获取所有窗口的第一个
            NSArray *windows = [NSApplication sharedApplication].windows;
            if (windows.count > 0) {
                window = [windows objectAtIndex:0];
            }
        }

        if (window != nil) {
            // 先清理之前的观察者和定时器
            if (windowObserver != nil) {
                [[NSNotificationCenter defaultCenter] removeObserver:windowObserver];
                windowObserver = nil;
            }
            if (resizeObserver != nil) {
                [[NSNotificationCenter defaultCenter] removeObserver:resizeObserver];
                resizeObserver = nil;
            }
            if (positionTimer != nil) {
                [positionTimer invalidate];
                positionTimer = nil;
            }

            // 立即调整按钮位置
            updateButtonPositions(window);

            // 监听窗口大小变化通知
            windowObserver = [[NSNotificationCenter defaultCenter] addObserverForName:NSWindowDidResizeNotification
                                                                                object:window
                                                                                 queue:[NSOperationQueue mainQueue]
                                                                            usingBlock:^(NSNotification *note) {
                NSWindow *resizedWindow = [note object];
                updateButtonPositions(resizedWindow);
            }];

            // 监听窗口结束大小调整通知
            resizeObserver = [[NSNotificationCenter defaultCenter] addObserverForName:NSWindowDidEndLiveResizeNotification
                                                                                object:window
                                                                                 queue:[NSOperationQueue mainQueue]
                                                                            usingBlock:^(NSNotification *note) {
                NSWindow *resizedWindow = [note object];
                // 延迟一小段时间，确保系统布局完成后再调整
                dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)(0.1 * NSEC_PER_SEC)), dispatch_get_main_queue(), ^{
                    updateButtonPositions(resizedWindow);
                });
            }];

            // 使用定时器持续检查并调整按钮位置（在窗口大小调整期间）
            positionTimer = [NSTimer scheduledTimerWithTimeInterval:0.05
                                                              repeats:YES
                                                                block:^(NSTimer *timer) {
                NSWindow *currentWindow = [NSApplication sharedApplication].mainWindow;
                if (currentWindow == nil) {
                    NSArray *windows = [NSApplication sharedApplication].windows;
                    if (windows.count > 0) {
                        currentWindow = [windows objectAtIndex:0];
                    }
                }
                if (currentWindow != nil && currentWindow == window) {
                    updateButtonPositions(currentWindow);
                }
            }];
        }
    });
}

void CleanupWindowButtonsObserver() {
    if (windowObserver != nil) {
        [[NSNotificationCenter defaultCenter] removeObserver:windowObserver];
        windowObserver = nil;
    }
    if (resizeObserver != nil) {
        [[NSNotificationCenter defaultCenter] removeObserver:resizeObserver];
        resizeObserver = nil;
    }
    if (positionTimer != nil) {
        [positionTimer invalidate];
        positionTimer = nil;
    }
}
*/
import "C"

// AdjustWindowButtons 调整 macOS 窗口控制按钮位置（仅 macOS 生效）
// 会自动监听窗口大小变化，保持按钮位置
func AdjustWindowButtons() {
	C.AdjustWindowButtons()
}

// CleanupWindowButtonsObserver 清理观察者（在应用关闭时调用）
func CleanupWindowButtonsObserver() {
	C.CleanupWindowButtonsObserver()
}
