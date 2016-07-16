package getdns

// #cgo LDFLAGS: -lgetdns
// #include <getdns/getdns_extra.h>
import "C"

type Error struct {
    rc int
}

func (err *Error) Code() int {
    return err.rc
}

func (err *Error) Error() string {
    return C.GoString(C.getdns_get_errorstr_by_id(C.uint16_t(err.rc)))
}
