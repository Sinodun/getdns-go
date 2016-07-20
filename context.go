package getdns

// #cgo LDFLAGS: -lgetdns
// #include <getdns/getdns_extra.h>
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
        return nil, &returnCodeError{int(rc)}
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
    err := checkExtensions(exts)
    if err != nil {
        return nil, err
    }
    var res *C.getdns_dict
    var cexts *C.getdns_dict
    cexts, err = convertDictToC(exts)
    defer C.getdns_dict_destroy(cexts)
    if err != nil {
        return nil, err
    }
    cname := C.CString(name)
    defer C.free(unsafe.Pointer(cname))
    rc := C.getdns_address_sync(c.ctx, cname, cexts, &res)
    if rc != RETURN_GOOD {
        return nil, &returnCodeError{int(rc)}
    }

    return createResult(res), nil
}

func (c *Context) General(name string, requestType uint, exts *Dict) (*Result, error) {
    err := checkExtensions(exts)
    if err != nil {
        return nil, err
    }
    var res *C.getdns_dict
    var cexts *C.getdns_dict
    cexts, err = convertDictToC(exts)
    defer C.getdns_dict_destroy(cexts)
    if err != nil {
        return nil, err
    }
    cname := C.CString(name)
    defer C.free(unsafe.Pointer(cname))
    rc := C.getdns_general_sync(c.ctx, cname, C.uint16_t(requestType), cexts, &res)
    if rc != RETURN_GOOD {
        return nil, &returnCodeError{int(rc)}
    }

    return createResult(res), nil
}

func (c *Context) Hostname(address Dict, exts *Dict) (*Result, error) {
    getdnsAddr, err := convertAddressDictToCallTypes(address)
    if err != nil {
        return nil, err
    }
    err = checkExtensions(exts)
    if err != nil {
        return nil, err
    }
    var res *C.getdns_dict
    var caddr *C.getdns_dict
    caddr, err = convertDictToC(&getdnsAddr)
    defer C.getdns_dict_destroy(caddr)
    if err != nil {
        return nil, err
    }
    var cexts *C.getdns_dict
    cexts, err = convertDictToC(exts)
    defer C.getdns_dict_destroy(cexts)
    if err != nil {
        return nil, err
    }
    rc := C.getdns_hostname_sync(c.ctx, caddr, cexts, &res)
    if rc != RETURN_GOOD {
        return nil, &returnCodeError{int(rc)}
    }

    return createResult(res), nil
}

func (c *Context) Service(name string, exts *Dict) (*Result, error) {
    err := checkExtensions(exts)
    if err != nil {
        return nil, err
    }
    var res *C.getdns_dict
    var cexts *C.getdns_dict
    cexts, err = convertDictToC(exts)
    defer C.getdns_dict_destroy(cexts)
    if err != nil {
        return nil, err
    }
    cname := C.CString(name)
    defer C.free(unsafe.Pointer(cname))
    rc := C.getdns_service_sync(c.ctx, cname, cexts, &res)
    if rc != RETURN_GOOD {
        return nil, &returnCodeError{int(rc)}
    }

    return createResult(res), nil
}

func (c *Context) GetAppendName() (int, error) {
    var res C.getdns_append_name_t
    rc := C.getdns_context_get_append_name(c.ctx, &res)
    if rc != RETURN_GOOD {
        return 0, &returnCodeError{int(rc)}
    }
    return int(res), nil
}

func (c *Context) SetAppendName(appendName int) error {
    rc := C.getdns_context_set_append_name(c.ctx, C.getdns_append_name_t(appendName))
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }
    return nil
}

func (c *Context) GetDNSRootServers() ([]Dict, error) {
    var list *C.getdns_list
    rc := C.getdns_context_get_dns_root_servers(c.ctx, &list)
    if rc != RETURN_GOOD {
        return nil, &returnCodeError{int(rc)}
    }

    return makeAddressDictList(list)
}

func (c *Context) SetDNSRootServers(servers []Dict) error {
    callList := make(List, 0, len(servers))
    for _, server := range servers {
        callServer, err := convertAddressDictToCallTypes(server)
        if err != nil {
            return err
        }
        callList = append(callList, callServer)
    }

    clist, err := convertListToC(&callList)
    if err != nil {
        return err
    }
    defer C.getdns_list_destroy(clist)
    rc := C.getdns_context_set_dns_root_servers(c.ctx, clist)
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }
    return nil
}
