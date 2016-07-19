package getdns

// #cgo LDFLAGS: -lgetdns
// #include <getdns/getdns.h>
import "C"

import (
    "runtime"
    "unsafe"
)

type Context struct {
    ctx *C.getdns_context
}

func CreateContext(setFromOS bool) (*Context, error) {
    var csetFromOS C.int = 0
    if setFromOS {
        csetFromOS = 1
    }
    var ctx *C.getdns_context
    rc := C.getdns_context_create(&ctx, csetFromOS)
    if rc != RETURN_GOOD {
        return nil, &Error{int(rc)}
    }

    res := &Context{ctx: ctx}
    runtime.SetFinalizer(res, (*Context).Destroy)
    return res, nil
}

func (c *Context) Destroy() {
    if ctx := c.ctx; c != nil {
        c.ctx = nil
        runtime.SetFinalizer(c, nil)
        C.getdns_context_destroy(ctx)
    }
}

func (c *Context) IsValid() bool {
    return c.ctx != nil
}

func (c *Context) Address(name string, exts *Dict) (*Result, error) {
    var res *C.getdns_dict
    var cexts *C.getdns_dict
    cexts, err := convertDictToC(exts)
    defer C.getdns_dict_destroy(cexts)
    if err != nil {
        return nil, err
    }
    cname := C.CString(name)
    defer C.free(unsafe.Pointer(cname))
    rc := C.getdns_address_sync(c.ctx, cname, cexts, &res)
    if rc != RETURN_GOOD {
        return nil, &Error{int(rc)}
    }

    return createResult(res), nil
}

func (c *Context) General(name string, requestType uint, exts *Dict) (*Result, error) {
    var res *C.getdns_dict
    var cexts *C.getdns_dict
    cexts, err := convertDictToC(exts)
    defer C.getdns_dict_destroy(cexts)
    if err != nil {
        return nil, err
    }
    cname := C.CString(name)
    defer C.free(unsafe.Pointer(cname))
    rc := C.getdns_general_sync(c.ctx, cname, C.uint16_t(requestType), cexts, &res)
    if rc != RETURN_GOOD {
        return nil, &Error{int(rc)}
    }

    return createResult(res), nil
}

func (c *Context) Service(name string, exts *Dict) (*Result, error) {
    var res *C.getdns_dict
    var cexts *C.getdns_dict
    cexts, err := convertDictToC(exts)
    defer C.getdns_dict_destroy(cexts)
    if err != nil {
        return nil, err
    }
    cname := C.CString(name)
    defer C.free(unsafe.Pointer(cname))
    rc := C.getdns_service_sync(c.ctx, cname, cexts, &res)
    if rc != RETURN_GOOD {
        return nil, &Error{int(rc)}
    }

    return createResult(res), nil
}
