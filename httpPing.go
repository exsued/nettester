package main

import (
    "log"
    "time"
    "net"
    "net/http"
    "strconv"
)

func HttpPingLoop(addrs []*net.IPAddr, interval float64, alarmInterval float64) {
    var err error
    intD := time.Duration(interval)
    alrmD:= time.Duration(alarmInterval)

    timer := time.AfterFunc(alrmD, func() {
        log.Println("Timer выполняется более минуты.")
    })

    for err == nil {
        for _, addr := range addrs { 
            status, err := httpPing(addr.String())         
            if(err != nil) {
                timer.Reset(alrmD)        
            }
            log.Println(addr.String() + " : : " + strconv.Itoa(status))
            time.Sleep(time.Duration(interval) * time.Second)        
        }
    }
}

func httpPing(domain string) (int, error) {
    url := "http://" + domain
    resp, err := http.Get(url)
    if err != nil {
       return 0, err
    }
    return resp.StatusCode, nil
}
