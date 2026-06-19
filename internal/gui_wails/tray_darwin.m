#import <Cocoa/Cocoa.h>
#include <stdio.h>

extern void onStatusItemClick(void);
extern void onQuitClick(void);
extern void onToggleLoggingClick(int enabled);
extern void onLaunchAtLoginClick(int enabled);
extern void onLaunchMinimizedClick(int enabled);
extern void onOpenLogsClick(void);

@interface StatusItemHandler : NSObject
- (void)statusItemClicked:(id)sender;
- (void)showClicked:(id)sender;
- (void)launchAtLoginClicked:(id)sender;
- (void)launchMinimizedClicked:(id)sender;
- (void)toggleLoggingClicked:(id)sender;
- (void)openLogsClicked:(id)sender;
- (void)quitClicked:(id)sender;
@end

static NSStatusItem *statusItem = nil;
static StatusItemHandler *handler = nil;
static NSMenu *contextMenu = nil;
static NSMenuItem *loggingMenuItem = nil;
static NSMenuItem *launchAtLoginMenuItem = nil;
static NSMenuItem *launchMinimizedMenuItem = nil;

@implementation StatusItemHandler

- (void)statusItemClicked:(id)sender {
    NSEvent *event = [NSApp currentEvent];
    if (event.type == NSEventTypeRightMouseDown || 
        (event.type == NSEventTypeLeftMouseDown && (event.modifierFlags & NSEventModifierFlagControl))) {
        
        printf("[objc] Right-click detected, popping up context menu\n");
        fflush(stdout);
        
        // Position the first menu item ("Show") exactly at the bottom-left of the status bar button (NSZeroPoint)
        if (contextMenu.numberOfItems > 0) {
            [contextMenu popUpMenuPositioningItem:[contextMenu itemAtIndex:0] atLocation:NSZeroPoint inView:statusItem.button];
        } else {
            [contextMenu popUpMenuPositioningItem:nil atLocation:NSZeroPoint inView:statusItem.button];
        }
    } else {
        printf("[objc] Left-click detected, toggling window\n");
        fflush(stdout);
        onStatusItemClick();
    }
}

- (void)showClicked:(id)sender {
    onStatusItemClick();
}

- (void)launchAtLoginClicked:(id)sender {
    BOOL newState = (launchAtLoginMenuItem.state == NSControlStateValueOff);
    [launchAtLoginMenuItem setState:newState ? NSControlStateValueOn : NSControlStateValueOff];
    onLaunchAtLoginClick(newState ? 1 : 0);
}

- (void)launchMinimizedClicked:(id)sender {
    BOOL newState = (launchMinimizedMenuItem.state == NSControlStateValueOff);
    [launchMinimizedMenuItem setState:newState ? NSControlStateValueOn : NSControlStateValueOff];
    onLaunchMinimizedClick(newState ? 1 : 0);
}

- (void)toggleLoggingClicked:(id)sender {
    BOOL newState = (loggingMenuItem.state == NSControlStateValueOff);
    [loggingMenuItem setState:newState ? NSControlStateValueOn : NSControlStateValueOff];
    onToggleLoggingClick(newState ? 1 : 0);
}

- (void)openLogsClicked:(id)sender {
    onOpenLogsClick();
}

- (void)quitClicked:(id)sender {
    onQuitClick();
}

@end

void initStatusItem(const char* title, const void* iconData, int iconLength, int initialLoggingState, int initialLaunchAtLoginState, int initialLaunchMinimizedState) {
    printf("[objc] initStatusItem called with title: %s, iconLength: %d\n", title, iconLength);
    fflush(stdout);
    
    // Copy parameters synchronously to Objective-C objects before Go frees the C memory
    NSString *nsTitle = [NSString stringWithUTF8String:title];
    NSData *nsIconData = nil;
    if (iconData != NULL && iconLength > 0) {
        nsIconData = [NSData dataWithBytes:iconData length:iconLength];
    }
    
    dispatch_async(dispatch_get_main_queue(), ^{
        printf("[objc] initStatusItem dispatch_async executing on main queue\n");
        fflush(stdout);
        
        // Set the activation policy to Accessory to hide the Dock icon programmatically
        [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
        
        NSStatusBar *bar = [NSStatusBar systemStatusBar];
        
        // Explicitly retain the status item to prevent it from being deallocated by the autorelease pool
        statusItem = [[bar statusItemWithLength:NSVariableStatusItemLength] retain];
        if (statusItem == nil) {
            printf("[objc] Error: Failed to create NSStatusItem\n");
            fflush(stdout);
            return;
        }
        
        NSButton *button = statusItem.button;
        if (button == nil) {
            printf("[objc] Error: NSStatusItem button is nil\n");
            fflush(stdout);
            return;
        }
        
        // Enable right-click event tracking on the button
        [button sendActionOn:(NSEventMaskLeftMouseDown | NSEventMaskRightMouseDown)];
        
        if (nsIconData != nil) {
            NSImage *image = [[NSImage alloc] initWithData:nsIconData];
            if (image != nil) {
                printf("[objc] Image created successfully from bytes\n");
                fflush(stdout);
                [image setSize:NSMakeSize(18, 18)];
                [image setTemplate:YES];
                button.image = image;
            } else {
                printf("[objc] Error: Failed to create NSImage from data, setting title instead\n");
                fflush(stdout);
                button.title = nsTitle;
            }
        } else {
            printf("[objc] No icon data, using title: %s\n", title);
            fflush(stdout);
            button.title = nsTitle;
        }
        
        button.toolTip = nsTitle;
        
        handler = [[StatusItemHandler alloc] init];
        button.target = handler;
        button.action = @selector(statusItemClicked:);
        
        // Setup context menu
        contextMenu = [[NSMenu alloc] init];
        [contextMenu setAutoenablesItems:NO];
        
        NSMenuItem *showItem = [[NSMenuItem alloc] initWithTitle:@"Show" action:@selector(showClicked:) keyEquivalent:@""];
        showItem.target = handler;
        [contextMenu addItem:showItem];
        
        [contextMenu addItem:[NSMenuItem separatorItem]];

        launchAtLoginMenuItem = [[NSMenuItem alloc] initWithTitle:@"Launch at Login" action:@selector(launchAtLoginClicked:) keyEquivalent:@""];
        launchAtLoginMenuItem.target = handler;
        [launchAtLoginMenuItem setState:initialLaunchAtLoginState ? NSControlStateValueOn : NSControlStateValueOff];
        [contextMenu addItem:launchAtLoginMenuItem];

        launchMinimizedMenuItem = [[NSMenuItem alloc] initWithTitle:@"Start Minimized" action:@selector(launchMinimizedClicked:) keyEquivalent:@""];
        launchMinimizedMenuItem.target = handler;
        [launchMinimizedMenuItem setState:initialLaunchMinimizedState ? NSControlStateValueOn : NSControlStateValueOff];
        [contextMenu addItem:launchMinimizedMenuItem];

        [contextMenu addItem:[NSMenuItem separatorItem]];
        
        loggingMenuItem = [[NSMenuItem alloc] initWithTitle:@"Enable Session Logging" action:@selector(toggleLoggingClicked:) keyEquivalent:@""];
        loggingMenuItem.target = handler;
        [loggingMenuItem setState:initialLoggingState ? NSControlStateValueOn : NSControlStateValueOff];
        [contextMenu addItem:loggingMenuItem];

        NSMenuItem *openLogsItem = [[NSMenuItem alloc] initWithTitle:@"Open Logs Directory" action:@selector(openLogsClicked:) keyEquivalent:@""];
        openLogsItem.target = handler;
        [contextMenu addItem:openLogsItem];
        
        [contextMenu addItem:[NSMenuItem separatorItem]];

        NSMenuItem *quitItem = [[NSMenuItem alloc] initWithTitle:@"Quit" action:@selector(quitClicked:) keyEquivalent:@"q"];
        quitItem.target = handler;
        [contextMenu addItem:quitItem];
        
        printf("[objc] initStatusItem execution completed successfully\n");
        fflush(stdout);
    });
}

void removeStatusItem(void) {
    printf("[objc] removeStatusItem called\n");
    fflush(stdout);
    dispatch_sync(dispatch_get_main_queue(), ^{
        if (statusItem != nil) {
            [[NSStatusBar systemStatusBar] removeStatusItem:statusItem];
            [statusItem release];
            statusItem = nil;
            printf("[objc] NSStatusItem removed and released successfully\n");
            fflush(stdout);
        }
    });
}

void getStatusItemPosition(double *x, double *y, double *width, double *height, double *screenWidth) {
    if (statusItem == nil) {
        *x = 0; *y = 0; *width = 0; *height = 0; *screenWidth = 0;
        return;
    }
    dispatch_sync(dispatch_get_main_queue(), ^{
        NSButton *button = statusItem.button;
        NSRect rect = button.bounds;
        // Convert button bounds to screen coordinates
        NSRect frame = [button convertRect:rect toView:nil];
        frame = [button.window convertRectToScreen:frame];
        
        NSScreen *screen = button.window.screen;
        if (screen == nil) {
            screen = [NSScreen mainScreen];
        }
        double screenHeight = screen.frame.size.height;
        *screenWidth = screen.frame.size.width;
        
        *x = frame.origin.x;
        // y is the bottom of the button in top-down coordinates
        *y = screenHeight - frame.origin.y;
        *width = frame.size.width;
        *height = frame.size.height;
    });
}

void activateApp(void) {
    dispatch_async(dispatch_get_main_queue(), ^{
        [NSApp activateIgnoringOtherApps:YES];
        for (NSWindow *window in [NSApp windows]) {
            if (![window.className containsString:@"StatusBar"]) {
                [window makeKeyAndOrderFront:nil];
            }
        }
    });
}
