package main

import (
    "fmt"
    "time"
    "net"
    "github.com/tatsushid/go-fastping"
)

func icmpPing(addrs []*net.IPAddr) error {
    p := fastping.NewPinger()

    for _,addr := range addrs {
        p.AddIPAddr(addr)
    }
    p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
	    fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
    }
    p.OnIdle = func() {
	    fmt.Println("finish")
    }
    err := p.Run()
    return err
}
