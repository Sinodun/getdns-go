package getdns

// #cgo LDFLAGS: -lgetdns
// #include <getdns/getdns_extra.h>
import "C"

import (
    "fmt"
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

type List []interface{}
type Dict map[string]interface{}

func convertDictToGo(dict *C.getdns_dict) (Dict, error) {
    var keys *C.getdns_list
    var nKeys C.size_t

    rc := C.getdns_dict_get_names(dict, &keys)
    if rc != RETURN_GOOD {
        return nil, &Error{int(rc)}
    }
    rc = C.getdns_list_get_length(keys, &nKeys)
    if rc != RETURN_GOOD {
        return nil, &Error{int(rc)}
    }

    res := make(Dict)
    for i := 0; i < int(nKeys); i++ {
        var binName *C.getdns_bindata
        rc = C.getdns_list_get_bindata(keys, C.size_t(i), &binName)
        if rc != RETURN_GOOD {
            return nil, &Error{int(rc)}
        }
        cName := (*C.char)(unsafe.Pointer(binName.data))
        keyName := C.GoString(cName)

        var dataType C.getdns_data_type
        rc = C.getdns_dict_get_data_type(dict, cName, &dataType)
        if rc != RETURN_GOOD {
            return nil, &Error{int(rc)}
        }

        var listItem *C.getdns_list
        var dictItem *C.getdns_dict
        var intItem C.uint32_t
        var bindataItem *C.getdns_bindata

        switch dataType {
        case C.t_list:
            rc = C.getdns_dict_get_list(dict, cName, &listItem)
            if rc != RETURN_GOOD {
                return nil, &Error{int(rc)}
            }
            li, err := convertListToGo(listItem)
            if err != nil {
                return nil, err
            }
            res[keyName] = li

        case C.t_dict:
            rc = C.getdns_dict_get_dict(dict, cName, &dictItem)
            if rc != RETURN_GOOD {
                return nil, &Error{int(rc)}
            }
            d, err := convertDictToGo(dictItem)
            if err != nil {
                return nil, err
            }
            res[keyName] = d

        case C.t_int:
            rc = C.getdns_dict_get_int(dict, cName, &intItem)
            if rc != RETURN_GOOD {
                return nil, &Error{int(rc)}
            }
            res[keyName] = int(intItem)

        case C.t_bindata:
            rc = C.getdns_dict_get_bindata(dict, cName, &bindataItem)
            if rc != RETURN_GOOD {
                return nil, &Error{int(rc)}
            }
            res[keyName] = bindataToByteSlice(bindataItem)

        default:
            return nil, &Error{RETURN_GENERIC_ERROR}
        }
    }

    return res, nil
}

func convertListToGo(list *C.getdns_list) (List, error) {
    var length C.size_t
    rc := C.getdns_list_get_length(list, &length)
    if rc != RETURN_GOOD {
        return nil, &Error{int(rc)}
    }

    res := make(List, 0, int(length))
    for i := 0; i < int(length); i++ {
        var dataType C.getdns_data_type
        var listItem *C.getdns_list
        var dictItem *C.getdns_dict
        var intItem C.uint32_t
        var bindataItem *C.getdns_bindata

        rc = C.getdns_list_get_data_type(list, C.size_t(i), &dataType)
        if rc != RETURN_GOOD {
            return nil, &Error{int(rc)}
        }

        switch dataType {
        case C.t_list:
            rc = C.getdns_list_get_list(list, C.size_t(i), &listItem)
            if rc != RETURN_GOOD {
                return nil, &Error{int(rc)}
            }
            li, err := convertListToGo(listItem)
            if err != nil {
                return nil, err
            }
            res = append(res, li)

        case C.t_dict:
            rc = C.getdns_list_get_dict(list, C.size_t(i), &dictItem)
            if rc != RETURN_GOOD {
                return nil, &Error{int(rc)}
            }
            d, err := convertDictToGo(dictItem)
            if err != nil {
                return nil, err
            }
            res = append(res, d)

        case C.t_int:
            rc = C.getdns_list_get_int(list, C.size_t(i), &intItem)
            if rc != RETURN_GOOD {
                return nil, &Error{int(rc)}
            }
            res = append(res, int(intItem))

        case C.t_bindata:
            rc = C.getdns_list_get_bindata(list, C.size_t(i), &bindataItem)
            if rc != RETURN_GOOD {
                return nil, &Error{int(rc)}
            }
            res = append(res, bindataToByteSlice(bindataItem))

        default:
            return nil, &Error{RETURN_GENERIC_ERROR}
        }
    }

    return res, nil
}

func (l *List) String() (res string) {
    res = "["
    first := true
    for _, item := range *l {
        if first {
            first = false
        } else {
            res = res + ","
        }
        switch val := item.(type) {
        case int:
            res = res + fmt.Sprintf("%d", val)
        case []byte:
            s, err := ConvertDNSNameToFQDN(val)
            if err != nil {
                res = res + string(val)
            } else {
                res = res + s
            }
        case List:
            res = res + val.String()
        case Dict:
            res = res + val.String()
        default:
            res = res + "Unknown"
        }
    }
    return res + "]"
}

func (d *Dict) String() (res string) {
    res = "{"
    first := true
    for key, item := range *d {
        if first {
            first = false
        } else {
            res = res + ","
        }
        res = res + fmt.Sprintf("%s: ", key)
        switch val := item.(type) {
        case int:
            res = res + fmt.Sprintf("%d", val)
        case []byte:
            s, err := ConvertDNSNameToFQDN(val)
            if err != nil {
                res = res + string(val)
            } else {
                res = res + s
            }
        case List:
            res = res + val.String()
        case Dict:
            res = res + val.String()
        default:
            res = res + "Unknown"
        }
    }
    return res + "}"
}

type Error struct {
    rc int
}

func (err *Error) Code() int {
    return err.rc
}

func (err *Error) Error() string {
    return fmt.Sprintf("getdns error %d: %s", err.rc, C.GoString(C.getdns_get_errorstr_by_id(C.uint16_t(err.rc))))
}

func ConvertDNSNameToFQDN(b []byte) (string, error) {
    res := ""
    p := 0
    if len(b) < 1 {
        return "", &Error{RETURN_GENERIC_ERROR}
    }
    for b[p] != 0 {
        labelLen := int(b[p])
        p = p + 1
        if labelLen > 63 || p+labelLen >= len(b) {
            return "", &Error{RETURN_GENERIC_ERROR}
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
