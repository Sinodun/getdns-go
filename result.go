package getdns

// #cgo LDFLAGS: -lgetdns
// #include <getdns/getdns.h>
import "C"

import (
    "runtime"
)

type Result struct {
    res *C.getdns_dict
}

func createResult(res *C.getdns_dict) *Result {
    r := &Result{res: res}
    runtime.SetFinalizer(r, (*Result).Destroy)
    return r
}

func (r *Result) Destroy() {
    if cres := r.res; cres != nil {
        r.res = nil
        runtime.SetFinalizer(r, nil)
        C.getdns_dict_destroy(cres)
    }
}

func (r *Result) IsValid() bool {
    return r.res != nil
}

func (r *Result) Status() (uint32, error) {
    var res C.uint32_t
    rc := C.getdns_dict_get_int(r.res, C.CString("status"), &res)
    if rc != RETURN_GOOD {
        return 0, &Error{int(rc)}
    }
    return uint32(res), nil
}
