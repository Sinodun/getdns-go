package getdns

// #cgo LDFLAGS: -lgetdns
// #include <getdns/getdns.h>
import "C"

import (
    "runtime"
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

func (c *Context) Address(name string) (*Result, error) {
    var res *C.getdns_dict
    rc := C.getdns_address_sync(c.ctx, C.CString(name), nil, &res)
    if rc != RETURN_GOOD {
        return nil, &Error{int(rc)}
    }

    return createResult(res), nil
}
