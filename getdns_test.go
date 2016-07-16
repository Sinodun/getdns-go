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
        t.Error(fmt.Sprintf("Bad status: %d", status))
    }
}
