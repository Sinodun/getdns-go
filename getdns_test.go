package getdns_test

import (
    "fmt"
    "testing"

    "getdns"
)

func TestContextCreate(t *testing.T) {
    c, err := getdns.CreateContext(true)
    if c == nil {
        t.Error(fmt.Sprintf("No Context created: %s", err))
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
        t.Error(fmt.Sprintf("No Context created: %s", err))
        return
    }
    defer c.Destroy()

    res, err := c.Address("www.lunch.org.uk")
    if res == nil {
        t.Error(fmt.Sprintf("No Result created: %s", err))
        return
    }
    status, err := res.Status()
    if err != nil {
        t.Error(fmt.Sprintf("No Status: %s", err))
        return
    }
    if status != getdns.RESPSTATUS_GOOD {
        t.Error(fmt.Sprintf("Bad Status: %d", status))
        return
    }
    ansType, err := res.AnswerType()
    if err != nil {
        t.Error(fmt.Sprintf("No AnswerType: %s", err))
        return
    }
    if ansType != getdns.NAMETYPE_DNS {
        t.Error(fmt.Sprintf("Bad AnswerType: %d", ansType))
        return
    }

    addrAns, err := res.JustAddressAnswers()
    if err != nil {
        t.Error(fmt.Sprintf("No JustAddressAnswers: %s", err))
        return
    }

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

    rt, err := res.RepliesTree()
    if err != nil {
        t.Error(fmt.Sprintf("No RepliesTree: %s", err))
        return
    }

    fmt.Print(rt.String())
}
