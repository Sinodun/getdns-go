package getdns

// #cgo LDFLAGS: -lgetdns
// #include <getdns/getdns.h>
import "C"

import (
    "net"
    "runtime"
    "unsafe"
)

var cCANONICAL_NAME = C.CString("canonical_name")
var cJUST_ADDRESS_ANSWERS = C.CString("just_address_answers")
var cREPLIES_TREE = C.CString("replies_tree")
var cVALIDATION_CHAIN = C.CString("validation_chain")

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

func (r *Result) getInt(key string) (uint32, error) {
    var res C.uint32_t
    ckey := C.CString(key)
    defer C.free(unsafe.Pointer(ckey))
    rc := C.getdns_dict_get_int(r.res, ckey, &res)
    if rc != RETURN_GOOD {
        return 0, &Error{int(rc)}
    }
    return uint32(res), nil
}

func (r *Result) AnswerType() (uint32, error) {
    return r.getInt("answer_type")
}

func (r *Result) CanonicalName() (string, error) {
    var bindata *C.getdns_bindata

    rc := C.getdns_dict_get_bindata(r.res, cCANONICAL_NAME, &bindata)
    if rc != RETURN_GOOD {
        return "", &Error{int(rc)}
    }

    b := bindataToByteSlice(bindata)
    dname, err := ConvertDNSNameToFQDN(b)
    if err != nil {
        return string(b), err
    } else {
        return dname, nil
    }
}

func (r *Result) JustAddressAnswers() ([]map[string]string, error) {
    var list *C.getdns_list

    rc := C.getdns_dict_get_list(r.res, cJUST_ADDRESS_ANSWERS, &list)
    if rc != RETURN_GOOD {
        return nil, &Error{int(rc)}
    }

    l, err := convertListToGo(list)
    if err != nil {
        return nil, err
    }

    res := make([]map[string]string, 0, len(l))
    for _, addrs := range l {
        item := make(map[string]string)
        d, ok := addrs.(Dict)
        if !ok {
            return nil, &Error{RETURN_GENERIC_ERROR}
        }
        addrType, ok := d["address_type"].([]byte)
        if !ok {
            return nil, &Error{RETURN_GENERIC_ERROR}
        }
        item["address_type"] = string(addrType)
        var ad net.IP
        ad, ok = d["address_data"].([]byte)
        if !ok {
            return nil, &Error{RETURN_GENERIC_ERROR}
        }
        item["address_data"] = ad.String()
        res = append(res, item)
    }
    return res, nil
}

func (r *Result) RepliesTree() (List, error) {
    var list *C.getdns_list

    rc := C.getdns_dict_get_list(r.res, cREPLIES_TREE, &list)
    if rc != RETURN_GOOD {
        return nil, &Error{int(rc)}
    }
    return convertListToGo(list)
}

func (r *Result) ValidationChain() (List, error) {
    var list *C.getdns_list

    rc := C.getdns_dict_get_list(r.res, cVALIDATION_CHAIN, &list)
    if rc != RETURN_GOOD {
        return nil, &Error{int(rc)}
    }
    return convertListToGo(list)
}

func (r *Result) Status() (uint32, error) {
    return r.getInt("status")
}
