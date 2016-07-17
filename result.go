package getdns

// #cgo LDFLAGS: -lgetdns
// #include <getdns/getdns.h>
import "C"

import (
    "net"
    "runtime"
    "unsafe"
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

func (r *Result) getInt(key string) (uint32, error) {
    var res C.uint32_t
    rc := C.getdns_dict_get_int(r.res, C.CString(key), &res)
    if rc != RETURN_GOOD {
        return 0, &Error{int(rc)}
    }
    return uint32(res), nil
}

func (r *Result) AnswerType() (uint32, error) {
    return r.getInt("answer_type")
}

func (r *Result) JustAddressAnswers() ([]map[string]string, error) {
    var res []map[string]string
    var list *C.getdns_list

    rc := C.getdns_dict_get_list(r.res, C.CString("just_address_answers"), &list)
    if rc != RETURN_GOOD {
        return nil, &Error{int(rc)}
    }

    var length C.size_t
    rc = C.getdns_list_get_length(list, &length)
    if rc != RETURN_GOOD {
        return nil, &Error{int(rc)}
    }

    for i := 0; i < int(length); i++ {
        var dataType C.getdns_data_type
        rc = C.getdns_list_get_data_type(list, C.size_t(i), &dataType)
        if rc != RETURN_GOOD {
            return nil, &Error{int(rc)}
        }

        if dataType != C.t_dict {
            return nil, &Error{RETURN_GENERIC_ERROR}
        }

        var dict *C.getdns_dict
        rc = C.getdns_list_get_dict(list, C.size_t(i), &dict)
        if rc != RETURN_GOOD {
            return nil, &Error{int(rc)}
        }

        var addrType, addrData *C.getdns_bindata
        rc = C.getdns_dict_get_bindata(dict, C.CString("address_type"), &addrType)
        if rc != RETURN_GOOD {
            return nil, &Error{int(rc)}
        }
        rc := C.getdns_dict_get_bindata(dict, C.CString("address_data"), &addrData)
        if rc != RETURN_GOOD {
            return nil, &Error{int(rc)}
        }

        item := make(map[string]string)
        item["address_type"] = C.GoString((*C.char)(unsafe.Pointer(addrType.data)))
        var addr net.IP = bindataToByteSlice(addrData)
        item["address_data"] = addr.String()

        res = append(res, item)
    }

    return res, nil
}

func (r *Result) Status() (uint32, error) {
    return r.getInt("status")
}
