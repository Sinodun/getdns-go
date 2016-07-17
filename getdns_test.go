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
    }

    if !c.IsValid() {
        t.Error("Context not valid on creation")
    }

    c.Destroy()
    if c.IsValid() {
        t.Error("Context not destroyed")
    }
}

func TestAddress(t *testing.T) {
    c, err := getdns.CreateContext(true)
    if c == nil {
        t.Error(fmt.Sprintf("No Context created: %s", err))
    }
    defer c.Destroy()

    res, err := c.Address("www.lunch.org.uk")
    if res == nil {
        t.Error(fmt.Sprintf("No Result created: %s", err))
    }
    status, err := res.Status()
    if err != nil {
        t.Error(fmt.Sprintf("No Status: %s", err))
    }
    if status != getdns.RESPSTATUS_GOOD {
        t.Error(fmt.Sprintf("Bad Status: %d", status))
    }
    ansType, err := res.AnswerType()
    if err != nil {
        t.Error(fmt.Sprintf("No AnswerType: %s", err))
    }
    if ansType != getdns.NAMETYPE_DNS {
        t.Error(fmt.Sprintf("Bad AnswerType: %d", ansType))
    }

    addrAns, err := res.JustAddressAnswers()
    if err != nil {
        t.Error(fmt.Sprintf("No JustAddressAnswers: %s", err))
    }

    if addrAns[0]["address_type"] != "IPv6" ||
        addrAns[0]["address_data"] != "2001:41c8:51:189:feff:ff:fe00:b1c" {
        t.Error("Bad IPv6 address")
    }
    if addrAns[1]["address_type"] != "IPv4" ||
        addrAns[1]["address_data"] != "213.138.101.137" {
        t.Error("Bad IPv6 address")
    }
}
