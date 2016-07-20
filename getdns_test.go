package getdns_test

import (
    "testing"

    "getdns"
)

func TestContextCreate(t *testing.T) {
    c, err := getdns.CreateContext(true)
    if c == nil {
        t.Errorf("No Context created: %s", err)
        return
    }

    if !c.IsValid() {
        t.Error("Context not valid on creation")
        return
    }

    c.Destroy()
    if c.IsValid() {
        t.Error("Context not destroyed")
        return
    }
}

func TestAddress(t *testing.T) {
    c, err := getdns.CreateContext(true)
    if c == nil {
        t.Fatalf("No Context created: %s", err)
    }
    defer c.Destroy()

    res, err := c.Address("www.lunch.org.uk", nil)
    if res == nil {
        t.Fatalf("No Result created: %s", err)
    }

    status, err := res.Status()
    if err != nil {
        t.Errorf("No Status: %s", err)
    } else if status != getdns.RESPSTATUS_GOOD {
        t.Fatalf("Bad Status: %d", status)
    }

    ansType, err := res.AnswerType()
    if err != nil {
        t.Errorf("No AnswerType: %s", err)
    } else if ansType != getdns.NAMETYPE_DNS {
        t.Errorf("Bad AnswerType: %d", ansType)
    }

    addrAns, err := res.JustAddressAnswers()
    if err != nil {
        t.Errorf("No JustAddressAnswers: %s", err)
    } else {
        if addrAns[0]["address_type"] != "IPv6" {
            t.Error("Bad IPv6 address_type")
        }
        if addrAns[1]["address_type"] != "IPv4" {
            t.Error("Bad IPv4 address_type")
        }
        if addrAns[0]["address_data"] != "2001:41c8:51:189:feff:ff:fe00:b1c" {
            t.Error("Bad IPv6 address_data")
        }
        if addrAns[1]["address_data"] != "213.138.101.137" {
            t.Error("Bad IPv4 address_data")
        }
    }

    rt, err := res.RepliesTree()
    if err != nil {
        t.Errorf("No RepliesTree: %s", err)
    } else if len(rt) == 0 {
        t.Error("RepliesTree empty")
    } else {
        d, ok := rt[0].(getdns.Dict)
        if !ok {
            t.Error("RepliesTree: no dict at [0]")
        } else {
            q, ok := d["question"].(getdns.Dict)
            if !ok {
                t.Error("RepliesTree: no question")
            } else {
                qname, ok := q["qname"].([]byte)
                if !ok {
                    t.Error("RepliesTree: no qname")
                } else {
                    fqdn, err := getdns.ConvertDNSNameToFQDN(qname)
                    if err != nil || fqdn != "www.lunch.org.uk." {
                        t.Errorf("QNAME incorrect: %s", qname)
                    }
                }
            }
        }
    }

    rf, err := res.RepliesFull()
    if err != nil {
        t.Errorf("No RepliesFull: %s", err)
    } else {
        _, ok := rf["replies_tree"].(getdns.List)
        if !ok {
            t.Error("RepliesFull: no replies_tree")
        }
    }

    can, err := res.CanonicalName()
    if err != nil {
        t.Errorf("No CanonicalName: %s", err)
    } else if can != "pigwidgeon.lunch.org.uk." {
        t.Errorf("Wrong canonical name: %s", can)
    }

    _, err = res.ValidationChain()
    if err == nil {
        t.Error("ValidationChain found!")
    }
}

func TestGeneral(t *testing.T) {
    c, err := getdns.CreateContext(true)
    if c == nil {
        t.Fatalf("No Context created: %s", err)
    }
    defer c.Destroy()

    exts := make(getdns.Dict, 1)
    exts["return_both_v4_and_v6"] = getdns.EXTENSION_TRUE
    res, err := c.General("lunch.org.uk", getdns.RRTYPE_MX, &exts)
    if res == nil {
        t.Fatalf("No Result created: %s", err)
    }

    rt, err := res.RepliesTree()
    if err != nil {
        t.Errorf("No RepliesTree: %s", err)
    } else if len(rt) == 0 {
        t.Error("RepliesTree empty")
    } else {
        d, ok := rt[0].(getdns.Dict)
        if !ok {
            t.Error("RepliesTree: no dict at [0]")
        } else {
            q, ok := d["question"].(getdns.Dict)
            if !ok {
                t.Error("RepliesTree: no question")
            } else {
                qtype, ok := q["qtype"].(int)
                if !ok {
                    t.Error("RepliesTree: no qtype")
                } else if qtype != getdns.RRTYPE_MX {
                    t.Errorf("QTYPE incorrect: %d", qtype)
                }
            }
        }
    }
}

func TestService(t *testing.T) {
    c, err := getdns.CreateContext(true)
    if c == nil {
        t.Fatalf("No Context created: %s", err)
    }
    defer c.Destroy()

    exts := make(getdns.Dict, 1)
    exts["return_both_v4_and_v6"] = getdns.EXTENSION_TRUE
    res, err := c.Service("_imap._tcp.gmail.com", &exts)
    if res == nil {
        t.Fatalf("No Result created: %s", err)
    }

    rt, err := res.RepliesTree()
    if err != nil {
        t.Errorf("No RepliesTree: %s", err)
    } else if len(rt) == 0 {
        t.Error("RepliesTree empty")
    } else {
        d, ok := rt[0].(getdns.Dict)
        if !ok {
            t.Error("RepliesTree: no dict at [0]")
        } else {
            q, ok := d["question"].(getdns.Dict)
            if !ok {
                t.Error("RepliesTree: no question")
            } else {
                qtype, ok := q["qtype"].(int)
                if !ok {
                    t.Error("RepliesTree: no qtype")
                } else if qtype != getdns.RRTYPE_SRV {
                    t.Errorf("QTYPE incorrect: %d", qtype)
                }
            }
        }
    }
}

func TestHostname(t *testing.T) {
    c, err := getdns.CreateContext(true)
    if c == nil {
        t.Fatalf("No Context created: %s", err)
    }
    defer c.Destroy()

    addr := make(getdns.Dict, 2)
    addr["address_type"] = "IPv6"
    addr["address_data"] = "2001:41c8:51:189:feff:ff:fe00:b1c"

    res, err := c.Hostname(addr, nil)
    if res == nil {
        t.Fatalf("No Result created: %s", err)
    }

    rt, err := res.RepliesTree()
    if err != nil {
        t.Errorf("No RepliesTree: %s", err)
    } else if len(rt) == 0 {
        t.Error("RepliesTree empty")
    } else {
        d, ok := rt[0].(getdns.Dict)
        if !ok {
            t.Error("RepliesTree: no dict at [0]")
        } else {
            q, ok := d["question"].(getdns.Dict)
            if !ok {
                t.Error("RepliesTree: no question")
            } else {
                qtype, ok := q["qtype"].(int)
                if !ok {
                    t.Error("RepliesTree: no qtype")
                } else if qtype != getdns.RRTYPE_PTR {
                    t.Errorf("QTYPE incorrect: %d", qtype)
                }
            }
        }
    }
}

func TestAppendName(t *testing.T) {
    c, err := getdns.CreateContext(true)
    if c == nil {
        t.Fatalf("No Context created: %s", err)
    }
    defer c.Destroy()

    err = c.SetAppendName(getdns.APPEND_NAME_NEVER)
    if err != nil {
        t.Fatalf("SetAppendName() failed: %s", err)
    }

    appendName, err := c.GetAppendName()
    if err != nil {
        t.Fatalf("No AppendName: %s", err)
    }
    if appendName != getdns.APPEND_NAME_NEVER {
        t.Fatalf("Bad AppendName: %d", appendName)
    }
}

func TestDNSRootServers(t *testing.T) {
    c, err := getdns.CreateContext(false)
    if c == nil {
        t.Fatalf("No Context created: %s", err)
    }
    defer c.Destroy()

    d := make(getdns.Dict)
    d["address_type"] = "IPv4"
    d["address_data"] = "213.138.101.137"
    r := make([]getdns.Dict, 0)
    r = append(r, d)
    err = c.SetDNSRootServers(r)
    if err != nil {
        t.Fatalf("Can't set DNS root server: %s", err)
    }

    roots, err := c.GetDNSRootServers()
    if err != nil {
        t.Fatalf("No DNS root servers: %s", err)
    }
    for _, root := range roots {
        t.Log(root.String())
    }
}

func TestDNSTransportList(t *testing.T) {
    c, err := getdns.CreateContext(true)
    if c == nil {
        t.Fatalf("No Context created: %s", err)
    }
    defer c.Destroy()

    err = c.SetDNSTransportList([]int{getdns.TRANSPORT_UDP, getdns.TRANSPORT_TCP, getdns.TRANSPORT_TLS})
    if err != nil {
        t.Fatalf("Can't set transport list: %s", err)
    }

    tl, err := c.GetDNSTransportList()
    if err != nil {
        t.Fatalf("No transport list: %s", err)
    }
    if len(tl) != 3 ||
        (tl[0] != getdns.TRANSPORT_UDP && tl[1] != getdns.TRANSPORT_TCP && tl[2] != getdns.TRANSPORT_TLS) {
        t.Fatal("Incorrect transport list")
    }
}

func TestDNSSECAllowedSkew(t *testing.T) {
    c, err := getdns.CreateContext(true)
    if c == nil {
        t.Fatalf("No Context created: %s", err)
    }
    defer c.Destroy()

    err = c.SetDNSSECAllowedSkew(1234)
    if err != nil {
        t.Fatalf("Can't set allowed skew: %s", err)
    }

    skew, err := c.GetDNSSECAllowedSkew()
    if err != nil {
        t.Fatalf("No allowed skew: %s", err)
    }
    if skew != 1234 {
        t.Fatal("Incorrect skew: %d", skew)
    }
}

func TestDNSSECTrustAnchors(t *testing.T) {
    c, err := getdns.CreateContext(true)
    if c == nil {
        t.Fatalf("No Context created: %s", err)
    }
    defer c.Destroy()

    anchors, err := c.GetDNSSECTrustAnchors()
    if err != nil {
        t.Fatalf("No trust anchors: %s", err)
    }
    err = c.SetDNSSECTrustAnchors(anchors)
    if err != nil {
        t.Fatalf("Can't set trust anchors: %s", err)
    }
}

func TestEDNSClientSubnetPrivate(t *testing.T) {
    c, err := getdns.CreateContext(true)
    if c == nil {
        t.Fatalf("No Context created: %s", err)
    }
    defer c.Destroy()

    err = c.SetEDNSClientSubnetPrivate(true)
    if err != nil {
        t.Fatalf("Can't set EDNS subnet: %s", err)
    }

    edns, err := c.GetEDNSClientSubnetPrivate()
    if err != nil {
        t.Fatalf("No EDNS subnet: %s", err)
    }
    if !edns {
        t.Fatal("Incorrect EDBS: %d", edns)
    }
}

func TestEDNSDoBit(t *testing.T) {
    c, err := getdns.CreateContext(true)
    if c == nil {
        t.Fatalf("No Context created: %s", err)
    }
    defer c.Destroy()

    err = c.SetEDNSDoBit(true)
    if err != nil {
        t.Fatalf("Can't set EDNS Do: %s", err)
    }

    edns, err := c.GetEDNSDoBit()
    if err != nil {
        t.Fatalf("No EDNS Do: %s", err)
    }
    if !edns {
        t.Fatal("Incorrect EDBS: %v", edns)
    }
}

func TestEDNSExtendedRcode(t *testing.T) {
    c, err := getdns.CreateContext(true)
    if c == nil {
        t.Fatalf("No Context created: %s", err)
    }
    defer c.Destroy()

    err = c.SetEDNSExtendedRcode(123)
    if err != nil {
        t.Fatalf("Can't set EDNS Do: %s", err)
    }

    edns, err := c.GetEDNSExtendedRcode()
    if err != nil {
        t.Fatalf("No EDNS Do: %s", err)
    }
    if edns != 123 {
        t.Fatal("Incorrect EDBS: %v", edns)
    }
}
