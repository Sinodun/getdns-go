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
