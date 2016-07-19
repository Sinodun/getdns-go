package getdns

// #cgo LDFLAGS: -lgetdns
// #include <getdns/getdns_extra.h>
import "C"

import (
    "fmt"
    "net"
    "unicode"
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

func convertDictToGo(dict *C.getdns_dict) (Dict, error) {
    var keys *C.getdns_list
    var nKeys C.size_t

    rc := C.getdns_dict_get_names(dict, &keys)
    if rc != RETURN_GOOD {
        return nil, &returnCodeError{int(rc)}
    }
    rc = C.getdns_list_get_length(keys, &nKeys)
    if rc != RETURN_GOOD {
        return nil, &returnCodeError{int(rc)}
    }

    res := make(Dict)
    for i := 0; i < int(nKeys); i++ {
        var binName *C.getdns_bindata
        rc = C.getdns_list_get_bindata(keys, C.size_t(i), &binName)
        if rc != RETURN_GOOD {
            return nil, &returnCodeError{int(rc)}
        }
        cName := (*C.char)(unsafe.Pointer(binName.data))
        keyName := C.GoString(cName)

        var dataType C.getdns_data_type
        rc = C.getdns_dict_get_data_type(dict, cName, &dataType)
        if rc != RETURN_GOOD {
            return nil, &returnCodeError{int(rc)}
        }

        var listItem *C.getdns_list
        var dictItem *C.getdns_dict
        var intItem C.uint32_t
        var bindataItem *C.getdns_bindata

        switch dataType {
        case C.t_list:
            rc = C.getdns_dict_get_list(dict, cName, &listItem)
            if rc != RETURN_GOOD {
                return nil, &returnCodeError{int(rc)}
            }
            li, err := convertListToGo(listItem)
            if err != nil {
                return nil, err
            }
            res[keyName] = li

        case C.t_dict:
            rc = C.getdns_dict_get_dict(dict, cName, &dictItem)
            if rc != RETURN_GOOD {
                return nil, &returnCodeError{int(rc)}
            }
            d, err := convertDictToGo(dictItem)
            if err != nil {
                return nil, err
            }
            res[keyName] = d

        case C.t_int:
            rc = C.getdns_dict_get_int(dict, cName, &intItem)
            if rc != RETURN_GOOD {
                return nil, &returnCodeError{int(rc)}
            }
            res[keyName] = int(intItem)

        case C.t_bindata:
            rc = C.getdns_dict_get_bindata(dict, cName, &bindataItem)
            if rc != RETURN_GOOD {
                return nil, &returnCodeError{int(rc)}
            }
            res[keyName] = bindataToByteSlice(bindataItem)

        default:
            return nil, &returnCodeError{RETURN_WRONG_TYPE_REQUESTED}
        }
    }

    return res, nil
}

func convertListToGo(list *C.getdns_list) (List, error) {
    var length C.size_t
    rc := C.getdns_list_get_length(list, &length)
    if rc != RETURN_GOOD {
        return nil, &returnCodeError{int(rc)}
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
            return nil, &returnCodeError{int(rc)}
        }

        switch dataType {
        case C.t_list:
            rc = C.getdns_list_get_list(list, C.size_t(i), &listItem)
            if rc != RETURN_GOOD {
                return nil, &returnCodeError{int(rc)}
            }
            li, err := convertListToGo(listItem)
            if err != nil {
                return nil, err
            }
            res = append(res, li)

        case C.t_dict:
            rc = C.getdns_list_get_dict(list, C.size_t(i), &dictItem)
            if rc != RETURN_GOOD {
                return nil, &returnCodeError{int(rc)}
            }
            d, err := convertDictToGo(dictItem)
            if err != nil {
                return nil, err
            }
            res = append(res, d)

        case C.t_int:
            rc = C.getdns_list_get_int(list, C.size_t(i), &intItem)
            if rc != RETURN_GOOD {
                return nil, &returnCodeError{int(rc)}
            }
            res = append(res, int(intItem))

        case C.t_bindata:
            rc = C.getdns_list_get_bindata(list, C.size_t(i), &bindataItem)
            if rc != RETURN_GOOD {
                return nil, &returnCodeError{int(rc)}
            }
            res = append(res, bindataToByteSlice(bindataItem))

        default:
            return nil, &returnCodeError{RETURN_WRONG_TYPE_REQUESTED}
        }
    }

    return res, nil
}

func convertDictToC(d *Dict) (*C.getdns_dict, error) {
    var res *C.getdns_dict

    if d == nil {
        return nil, nil
    }

    res = C.getdns_dict_create()
    if res == nil {
        return nil, &returnCodeError{RETURN_MEMORY_ERROR}
    }

    for key, item := range *d {
        ckey := C.CString(key)
        defer C.free(unsafe.Pointer(ckey))

        var rc C.getdns_return_t
        switch val := item.(type) {
        case int:
            rc = C.getdns_dict_set_int(res, ckey, C.uint32_t(val))

        case []byte:
            var bindata C.getdns_bindata
            bindata.size = C.size_t(len(val))
            bindata.data = (*C.uint8_t)(&val[0])
            rc = C.getdns_dict_set_bindata(res, ckey, &bindata)

        case Dict:
            d, err := convertDictToC(&val)
            if err != nil {
                C.getdns_dict_destroy(res)
                return nil, err
            }
            rc = C.getdns_dict_set_dict(res, ckey, d)

        case List:
            l, err := convertListToC(&val)
            if err != nil {
                C.getdns_dict_destroy(res)
                return nil, err
            }
            rc = C.getdns_dict_set_list(res, ckey, l)

        default:
            C.getdns_dict_destroy(res)
            return nil, &returnCodeError{RETURN_WRONG_TYPE_REQUESTED}
        }
        if rc != RETURN_GOOD {
            C.getdns_dict_destroy(res)
            return nil, &returnCodeError{int(rc)}
        }
    }

    return res, nil
}

func convertListToC(l *List) (*C.getdns_list, error) {
    var res *C.getdns_list

    if l == nil {
        return nil, nil
    }

    res = C.getdns_list_create()
    if res == nil {
        return nil, &returnCodeError{RETURN_MEMORY_ERROR}
    }

    for i, item := range *l {
        var rc C.getdns_return_t
        switch val := item.(type) {
        case int:
            rc = C.getdns_list_set_int(res, C.size_t(i), C.uint32_t(val))

        case []byte:
            var bindata C.getdns_bindata
            bindata.size = C.size_t(len(val))
            bindata.data = (*C.uint8_t)(&val[0])
            rc = C.getdns_list_set_bindata(res, C.size_t(i), &bindata)

        case Dict:
            d, err := convertDictToC(&val)
            if err != nil {
                C.getdns_list_destroy(res)
                return nil, err
            }
            rc = C.getdns_list_set_dict(res, C.size_t(i), d)

        case List:
            l, err := convertListToC(&val)
            if err != nil {
                C.getdns_list_destroy(res)
                return nil, err
            }
            rc = C.getdns_list_set_list(res, C.size_t(i), l)

        default:
            C.getdns_list_destroy(res)
            return nil, &returnCodeError{RETURN_WRONG_TYPE_REQUESTED}
        }
        if rc != RETURN_GOOD {
            C.getdns_list_destroy(res)
            return nil, &returnCodeError{int(rc)}
        }
    }

    return res, nil
}

func checkExtensions(exts *Dict) error {
    if exts == nil {
        return nil
    }

    var retcall string
    var ok bool
    if C.GETDNS_NUMERIC_VERSION < 0x00090000 {
        retcall = "return_call_debugging"
    } else {
        retcall = "return_call_reporting"
    }
    for key, item := range *exts {
        switch key {
        case retcall,
            "add_warning_for_bad_dns",
            "dnssec_return_status",
            "dnssec_return_all_statuses",
            "dnssec_return_only_secure",
            "dnssec_return_validation_chain",
            "return_api_information",
            "return_both_v4_and_v6":
            ival, ok := item.(int)
            if !ok || (ival != EXTENSION_TRUE && ival != EXTENSION_FALSE) {
                return &returnCodeError{RETURN_EXTENSION_MISFORMAT}
            }

        case "specify_class":
            _, ok = item.(int)
            if !ok {
                return &returnCodeError{RETURN_EXTENSION_MISFORMAT}
            }

        case "add_opt_parameters":
            optdict, ok := item.(Dict)
            if !ok {
                return &returnCodeError{RETURN_EXTENSION_MISFORMAT}
            }
            for optkey, optval := range optdict {
                switch optkey {
                case "maximum_udp_payload_size",
                    "extended_rcode",
                    "version",
                    "do_bit":
                    _, ok = optval.(int)
                    if !ok {
                        return &returnCodeError{RETURN_EXTENSION_MISFORMAT}
                    }

                case "options":
                    l, ok := optval.(List)
                    if !ok {
                        return &returnCodeError{RETURN_EXTENSION_MISFORMAT}
                    }
                    for _, listitem := range l {
                        ld, ok := listitem.(Dict)
                        if !ok {
                            return &returnCodeError{RETURN_EXTENSION_MISFORMAT}
                        }
                        for lkey, ldata := range ld {
                            switch lkey {
                            case "option_code":
                                _, ok = ldata.(int)

                            case "option_data":
                                _, ok = ldata.([]byte)

                            default:
                                ok = false
                            }
                        }
                        if !ok {
                            return &returnCodeError{RETURN_EXTENSION_MISFORMAT}
                        }
                    }

                default:
                    return &returnCodeError{RETURN_EXTENSION_MISFORMAT}
                }
            }

        default:
            return &returnCodeError{RETURN_NO_SUCH_EXTENSION}
        }
    }

    return nil
}

func val2str(item interface{}, key *string) string {
    switch val := item.(type) {
    case int:
        return fmt.Sprintf("%d", val)
    case []byte:
        printable := true
        for _, c := range string(val) {
            if !unicode.IsPrint(c) {
                printable = false
            }
        }
        if printable {
            return "'" + string(val) + "'"
        }
        s, err := ConvertDNSNameToFQDN(val)
        if err == nil {
            return "'" + s + "'"
        }
        if key != nil && *key == "address_data" {
            var ip net.IP = val
            return ip.String()
        }
        return fmt.Sprintf("'% x'", string(val))
    case List:
        return val.String()
    case Dict:
        return val.String()
    default:
        return "<Unknown>"
    }
}

func (l *List) String() (res string) {
    res = "["
    first := true
    for _, item := range *l {
        if first {
            first = false
        } else {
            res = res + ", "
        }
        res = res + val2str(item, nil)
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
            res = res + ", "
        }
        res = res + fmt.Sprintf("'%s': ", key) + val2str(item, &key)
    }
    return res + "}"
}
