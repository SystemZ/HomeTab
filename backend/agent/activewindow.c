
// -*- coding:utf-8-unix; mode:c; -*-
/*
  get the active window on X window system
  https://gist.github.com/kui/2622504
  http://k-ui.jp/blog/2012/05/07/get-active-window-on-x-window-system/
 */

#include <stdlib.h>
#include <stdio.h>
#include <locale.h>

#include <X11/Xlib.h>           // `apt-get install libx11-dev`
#include <X11/Xmu/WinUtil.h>    // `apt-get install libxmu-dev`

Bool xerror = False;

Display* open_display(){
    Display* d = XOpenDisplay(NULL);
    if(d == NULL){
        printf("fail\n");
        exit(1);
    }
    return d;
}

int handle_error(Display* display, XErrorEvent* error){
    printf("ERROR: X11 error\n");
    xerror = True;
    return 1;
}

Window get_focus_window(Display* d){
    Window w;
    int revert_to;
    XGetInputFocus(d, &w, &revert_to); // see man
    if(xerror){
        printf("fail\n");
        exit(1);
    }else if(w == None){
        printf("no focus window\n");
        exit(1);
    }

    return w;
}

// get the top window.
// a top window have the following specifications.
//  * the start window is contained the descendent windows.
//  * the parent window is the root window.
Window get_top_window(Display* d, Window start){
    Window w = start;
    Window parent = start;
    Window root = None;
    Window *children;
    unsigned int nchildren;
    Status s;

    while (parent != root) {
        w = parent;
        s = XQueryTree(d, w, &root, &parent, &children, &nchildren); // see man
        if (s)
            XFree(children);
        if(xerror){
            printf("fail\n");
            exit(1);
        }
    }
    return w;
}

// search a named window (that has a WM_STATE prop)
// on the descendent windows of the argment Window.
Window get_named_window(Display* d, Window start){
    Window w;
    w = XmuClientWindow(d, start); // see man
    if(w == start)
        printf("fail\n");
    return w;
}

// (XFetchName cannot get a name with multi-byte chars)
const char * print_window_name(Display* d, Window w){
    XTextProperty prop;
    Status s;

    s = XGetWMName(d, w, &prop); // see man
    if(!xerror && s){
        int count = 0, result;
        char **list = NULL;
        result = XmbTextPropertyToTextList(d, &prop, &list, &count); // see man
        if(result == Success){
            return list[0];
            //printf("\t%s\n", list[0]);
        }else{
            printf("ERROR: XmbTextPropertyToTextList\n");
        }
    }else{
        printf("ERROR: XGetWMName\n");
    }
}

const char * activeWindowName() {
    Display* d;
    Window w;

    setlocale(LC_ALL, ""); // see man locale

    d = open_display();
    XSetErrorHandler(handle_error);

    // get active window
    w = get_focus_window(d);
    w = get_top_window(d, w);
    w = get_named_window(d, w);

    return print_window_name(d, w);
}