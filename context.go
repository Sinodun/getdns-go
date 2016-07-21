package getdns

// #cgo LDFLAGS: -lgetdns
// #include <getdns/getdns_extra.h>
import "C"

import (
    "runtime"
    "unsafe"
)

type Context struct {
    ctx                  *C.getdns_context
    implementationString string
    versionString        string
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

    apiInfo, err := res.ApiInformation()
    if err != nil {
        res.Destroy()
        return nil, err
    }
    var is, vs []byte
    is, ok := apiInfo["implementation_string"].([]byte)
    if ok {
        vs, ok = apiInfo["version_string"].([]byte)
    }
    if !ok {
        res.Destroy()
        return nil, &returnCodeError{RETURN_NO_SUCH_DICT_NAME}
    }
    res.implementationString = string(is)
    res.versionString = string(vs)

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

func (c *Context) Address(name string, exts Dict) (*Result, error) {
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

func (c *Context) General(name string, requestType uint, exts Dict) (*Result, error) {
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

func (c *Context) ApiInformation() (Dict, error) {
    res, err := convertDictToGo(C.getdns_context_get_api_information(c.ctx))
    if err != nil {
        return nil, err
    }
    return res, nil
}

func (c *Context) Hostname(address Dict, exts Dict) (*Result, error) {
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
    caddr, err = convertDictToC(getdnsAddr)
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

func (c *Context) Service(name string, exts Dict) (*Result, error) {
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

func (c *Context) AppendName() (int, error) {
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

func (c *Context) DNSRootServers() ([]Dict, error) {
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

    clist, err := convertListToC(callList)
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

func (c *Context) DNSTransportList() ([]int, error) {
    var list *C.getdns_transport_list_t
    var listSize C.size_t
    rc := C.getdns_context_get_dns_transport_list(c.ctx, &listSize, &list)
    if rc != RETURN_GOOD {
        return nil, &returnCodeError{int(rc)}
    }

    res := make([]int, int(listSize))
    cres := (*[1 << 30]C.int)(unsafe.Pointer(list))[:listSize:listSize]
    for i, val := range cres {
        res[i] = int(val)
    }
    return res, nil
}

func (c *Context) SetDNSTransportList(list []int) error {
    clist := make([]C.int, len(list))
    for i, val := range list {
        clist[i] = C.int(val)
    }
    rc := C.getdns_context_set_dns_transport_list(c.ctx, C.size_t(len(list)), (*C.getdns_transport_list_t)(unsafe.Pointer(&clist[0])))
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) DNSSECAllowedSkew() (uint32, error) {
    var res C.uint32_t
    rc := C.getdns_context_get_dnssec_allowed_skew(c.ctx, &res)
    if rc != RETURN_GOOD {
        return 0, &returnCodeError{int(rc)}
    }

    return uint32(res), nil
}

func (c *Context) SetDNSSECAllowedSkew(skew uint32) error {
    var cskew C.uint32_t = C.uint32_t(skew)
    rc := C.getdns_context_set_dnssec_allowed_skew(c.ctx, cskew)
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) DNSSECTrustAnchors() (List, error) {
    var list *C.getdns_list
    rc := C.getdns_context_get_dnssec_trust_anchors(c.ctx, &list)
    if rc != RETURN_GOOD {
        return nil, &returnCodeError{int(rc)}
    }

    return convertListToGo(list)
}

func (c *Context) SetDNSSECTrustAnchors(anchors List) error {
    canchors, err := convertListToC(anchors)
    if err != nil {
        return err
    }
    defer C.getdns_list_destroy(canchors)
    rc := C.getdns_context_set_dnssec_trust_anchors(c.ctx, canchors)
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) EDNSClientSubnetPrivate() (bool, error) {
    var val C.uint8_t
    rc := C.getdns_context_get_edns_client_subnet_private(c.ctx, &val)
    if rc != RETURN_GOOD {
        return false, &returnCodeError{int(rc)}
    }

    return (val == 1), nil
}

func (c *Context) SetEDNSClientSubnetPrivate(private bool) error {
    var cskew C.uint8_t = 0
    if private {
        cskew = 1
    }
    rc := C.getdns_context_set_edns_client_subnet_private(c.ctx, cskew)
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) EDNSDoBit() (bool, error) {
    var val C.uint8_t
    rc := C.getdns_context_get_edns_do_bit(c.ctx, &val)
    if rc != RETURN_GOOD {
        return false, &returnCodeError{int(rc)}
    }

    return (val == 1), nil
}

func (c *Context) SetEDNSDoBit(newval bool) error {
    var do C.uint8_t = 0
    if newval {
        do = 1
    }
    rc := C.getdns_context_set_edns_do_bit(c.ctx, do)
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) EDNSExtendedRcode() (uint8, error) {
    var val C.uint8_t
    rc := C.getdns_context_get_edns_extended_rcode(c.ctx, &val)
    if rc != RETURN_GOOD {
        return 0, &returnCodeError{int(rc)}
    }

    return uint8(val), nil
}

func (c *Context) SetEDNSExtendedRcode(newval uint8) error {
    rc := C.getdns_context_set_edns_extended_rcode(c.ctx, C.uint8_t(newval))
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) EDNSMaximumUDPPayloadSize() (uint16, error) {
    var val C.uint16_t
    rc := C.getdns_context_get_edns_maximum_udp_payload_size(c.ctx, &val)
    if rc != RETURN_GOOD {
        return 0, &returnCodeError{int(rc)}
    }

    return uint16(val), nil
}

func (c *Context) SetEDNSMaximumUDPPayloadSize(newval uint16) error {
    rc := C.getdns_context_set_edns_maximum_udp_payload_size(c.ctx, C.uint16_t(newval))
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) EDNSVersion() (uint8, error) {
    var val C.uint8_t
    rc := C.getdns_context_get_edns_version(c.ctx, &val)
    if rc != RETURN_GOOD {
        return 0, &returnCodeError{int(rc)}
    }

    return uint8(val), nil
}

func (c *Context) SetEDNSVersion(newval uint8) error {
    rc := C.getdns_context_set_edns_version(c.ctx, C.uint8_t(newval))
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) FollowRedirects() (int, error) {
    var val C.getdns_redirects_t
    rc := C.getdns_context_get_follow_redirects(c.ctx, &val)
    if rc != RETURN_GOOD {
        return 0, &returnCodeError{int(rc)}
    }

    return int(val), nil
}

func (c *Context) SetFollowRedirects(newval int) error {
    rc := C.getdns_context_set_follow_redirects(c.ctx, C.getdns_redirects_t(newval))
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) IdleTimeout() (uint64, error) {
    var val C.uint64_t
    rc := C.getdns_context_get_idle_timeout(c.ctx, &val)
    if rc != RETURN_GOOD {
        return 0, &returnCodeError{int(rc)}
    }

    return uint64(val), nil
}

func (c *Context) SetIdleTimeout(newval uint64) error {
    rc := C.getdns_context_set_idle_timeout(c.ctx, C.uint64_t(newval))
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) LimitOutstandingQueries() (uint16, error) {
    var val C.uint16_t
    rc := C.getdns_context_get_limit_outstanding_queries(c.ctx, &val)
    if rc != RETURN_GOOD {
        return 0, &returnCodeError{int(rc)}
    }

    return uint16(val), nil
}

func (c *Context) SetLimitOutstandingQueries(newval uint16) error {
    rc := C.getdns_context_set_limit_outstanding_queries(c.ctx, C.uint16_t(newval))
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) Namespaces() ([]int, error) {
    var list *C.getdns_namespace_t
    var listSize C.size_t
    rc := C.getdns_context_get_namespaces(c.ctx, &listSize, &list)
    if rc != RETURN_GOOD {
        return nil, &returnCodeError{int(rc)}
    }

    res := make([]int, int(listSize))
    cres := (*[1 << 30]C.int)(unsafe.Pointer(list))[:listSize:listSize]
    for i, val := range cres {
        res[i] = int(val)
    }
    return res, nil
}

func (c *Context) SetNamespaces(list []int) error {
    clist := make([]C.int, len(list))
    for i, val := range list {
        clist[i] = C.int(val)
    }
    rc := C.getdns_context_set_namespaces(c.ctx, C.size_t(len(list)), (*C.getdns_namespace_t)(unsafe.Pointer(&clist[0])))
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) ResolutionType() (int, error) {
    var res C.getdns_resolution_t
    rc := C.getdns_context_get_resolution_type(c.ctx, &res)
    if rc != RETURN_GOOD {
        return 0, &returnCodeError{int(rc)}
    }
    return int(res), nil
}

func (c *Context) SetResolutionType(newval int) error {
    rc := C.getdns_context_set_resolution_type(c.ctx, C.getdns_resolution_t(newval))
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }
    return nil
}

func (c *Context) Suffix() ([]string, error) {
    var list *C.getdns_list
    rc := C.getdns_context_get_suffix(c.ctx, &list)
    if rc != RETURN_GOOD {
        return nil, &returnCodeError{int(rc)}
    }

    glist, err := convertListToGo(list)
    if err != nil {
        return nil, err
    }

    res := make([]string, len(glist))
    for i, val := range glist {
        b, ok := val.([]byte)
        if !ok {
            return nil, &returnCodeError{RETURN_GENERIC_ERROR}
        }
        // In 0.9.0 this returns a zero byte as the last byte.
        if b[len(b)-1] == 0 {
            b = b[0 : len(b)-1]
        }
        res[i] = string(b)
    }

    return res, nil
}

func (c *Context) SetSuffix(list []string) error {
    glist := make(List, len(list))
    for i, val := range list {
        glist[i] = val
    }
    clist, err := convertListToC(glist)
    if err != nil {
        return err
    }
    defer C.getdns_list_destroy(clist)
    rc := C.getdns_context_set_suffix(c.ctx, clist)
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) Timeout() (uint64, error) {
    var val C.uint64_t
    rc := C.getdns_context_get_timeout(c.ctx, &val)
    if rc != RETURN_GOOD {
        return 0, &returnCodeError{int(rc)}
    }

    return uint64(val), nil
}

func (c *Context) SetTimeout(newval uint64) error {
    rc := C.getdns_context_set_timeout(c.ctx, C.uint64_t(newval))
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) TLSAuthentication() (int, error) {
    var val C.getdns_tls_authentication_t
    rc := C.getdns_context_get_tls_authentication(c.ctx, &val)
    if rc != RETURN_GOOD {
        return 0, &returnCodeError{int(rc)}
    }

    return int(val), nil
}

func (c *Context) SetTLSAuthentication(newval int) error {
    rc := C.getdns_context_set_tls_authentication(c.ctx, C.getdns_tls_authentication_t(newval))
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) TLSQueryPaddingBlocksize() (uint16, error) {
    var val C.uint16_t
    rc := C.getdns_context_get_tls_query_padding_blocksize(c.ctx, &val)
    if rc != RETURN_GOOD {
        return 0, &returnCodeError{int(rc)}
    }

    return uint16(val), nil
}

func (c *Context) SetTLSQueryPaddingBlocksize(newval uint16) error {
    rc := C.getdns_context_set_tls_query_padding_blocksize(c.ctx, C.uint16_t(newval))
    if rc != RETURN_GOOD {
        return &returnCodeError{int(rc)}
    }

    return nil
}

func (c *Context) ImplementationString() string {
    return c.implementationString
}

func (c *Context) VersionString() string {
    return c.versionString
}
