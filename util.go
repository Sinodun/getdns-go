package getdns

// #cgo LDFLAGS: -lgetdns
// #include <getdns/getdns_extra.h>
import "C"

import (
    "unsafe"
)

func bindataToByteSlice(bindata *C.getdns_bindata) []byte {
    p := uintptr(unsafe.Pointer(bindata.data))
    var res = make([]byte, int(bindata.size))

    for i := 0; i < len(res); i++ {
        res[i] = byte(*(*C.uint8_t)(unsafe.Pointer(p)))
        p++
    }

    return res
}

type Error struct {
    rc int
}

func (err *Error) Code() int {
    return err.rc
}

func (err *Error) Error() string {
    return C.GoString(C.getdns_get_errorstr_by_id(C.uint16_t(err.rc)))
}
