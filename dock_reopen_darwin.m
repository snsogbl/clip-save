#import <Cocoa/Cocoa.h>
#import <objc/runtime.h>

extern void goOnAppReopen(void);
extern void goSetForceQuit(void);

@interface DockReopenObserver : NSObject
@end

@implementation DockReopenObserver
- (void)handleActive:(NSNotification*)note {
    goOnAppReopen();
}
@end

static id dockReopenObserverInstance = nil;

void RegisterReopenObserver(void) {
    if (dockReopenObserverInstance != nil) {
        return;
    }
    dockReopenObserverInstance = [DockReopenObserver new];
    [[NSNotificationCenter defaultCenter] addObserver:dockReopenObserverInstance
                                             selector:@selector(handleActive:)
                                                 name:NSApplicationDidBecomeActiveNotification
                                               object:nil];

    // Dynamically add applicationShouldHandleReopen to Wails' AppDelegate at runtime
    Class delegateClass = NSClassFromString(@"AppDelegate");
    if (delegateClass) {
        SEL sel = @selector(applicationShouldHandleReopen:hasVisibleWindows:);
        // Implementation that calls back into Go and returns NO (preserve default behavior)
        IMP imp = imp_implementationWithBlock(^BOOL(id _self, NSApplication *sender, BOOL flag) {
            goOnAppReopen();
            return NO;
        });

        // If not already implemented, add it (use encoding for: BOOL)id self, SEL _cmd, (NSApplication*)@(object), (BOOL)char
        // "c@:@c" is commonly used: return char, self, _cmd, object, char
        if (!class_addMethod(delegateClass, sel, imp, "c@:@c")) {
            // If exists, swizzle to inject our behavior first
            Method m = class_getInstanceMethod(delegateClass, sel);
            if (m) {
                IMP original = method_getImplementation(m);
                IMP replacement = imp_implementationWithBlock(^BOOL(id _self, NSApplication *sender, BOOL flag){
                    goOnAppReopen();
                    typedef BOOL (*Fn)(id, SEL, NSApplication*, BOOL);
                    return ((Fn)original)(_self, sel, sender, flag);
                });
                method_setImplementation(m, replacement);
            }
        }
    }

    // Also hook applicationShouldTerminate to mark force quit
    SEL quitSel = @selector(applicationShouldTerminate:);
    Method quitM = class_getInstanceMethod(delegateClass, quitSel);
    if (quitM) {
        IMP originalQuit = method_getImplementation(quitM);
        IMP replacementQuit = imp_implementationWithBlock(^NSInteger(id _self, NSApplication *sender){
            goSetForceQuit();
            typedef NSInteger (*Fn)(id, SEL, NSApplication*);
            return ((Fn)originalQuit)(_self, quitSel, sender);
        });
        method_setImplementation(quitM, replacementQuit);
    }
}


