package getdns

// #cgo LDFLAGS: -lgetdns
// #include <getdns/getdns_extra.h>
import "C"

import (
    "fmt"
)

// List is a Go type representing a getdns_list.
type List []interface{}

// Dict is a Go type representing a getdns_dict.
type Dict map[string]interface{}

// Error reports a getdns return code.
type Error interface {
    error
    ReturnCode() int
}

type returnCodeError struct {
    rc int
}

// Code returns the getdns numeric return code.
func (err *returnCodeError) ReturnCode() int {
    return err.rc
}

// Error implements the error interface and returns a printable
// description of the error.
func (err *returnCodeError) Error() string {
    return fmt.Sprintf("getdns error %d: %s", err.rc, C.GoString(C.getdns_get_errorstr_by_id(C.uint16_t(err.rc))))
}

// ConvertDNSNametoFQDN converts a name in DNS label format to a FQDN.
// It reimplements the getdns library routine in pure Go rather than
// calling into the library.
func ConvertDNSNameToFQDN(b []byte) (string, error) {
    res := ""
    p := 0
    if len(b) < 1 {
        return "", &returnCodeError{RETURN_BAD_DOMAIN_NAME}
    }
    for b[p] != 0 {
        labelLen := int(b[p])
        p = p + 1
        if labelLen > 63 || p+labelLen >= len(b) {
            return "", &returnCodeError{RETURN_BAD_DOMAIN_NAME}
        }
        labelContent := b[p : p+labelLen]
        res = res + string(labelContent) + "."
        p = p + int(labelLen)
    }
    if len(res) == 0 {
        res = "."
    }
    return res, nil
}
