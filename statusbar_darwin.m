#import <AppKit/AppKit.h>
#import <dispatch/dispatch.h>

void setActivationPolicy(int policy) {
    // AppKit 的操作必须在主线程执行
    dispatch_async(dispatch_get_main_queue(), ^{
        if (policy == 3) {
            // 设置为 "Prohibited " 模式，完全的后台进程，没有 UI
            [NSApp setActivationPolicy:NSApplicationActivationPolicyProhibited];
        } else if (policy == 2) {
            // 设置为 "Accessory" 模式，隐藏 Dock 图标
            [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
        } else if (policy == 1) {
            // 设置为 "Regular" 模式，显示 Dock 图标
            [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
        }
    });
}

