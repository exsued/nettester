package main

import (
    "net"
    "log"
    "flag"

    "io/ioutil"
    "strings"
)

func parseConf(filePath string) ([]*net.IPAddr, error) {
    bytesRead, err := ioutil.ReadFile("sites.txt")
    if err != nil {
        return nil, err    
    }
    fileContent := string(bytesRead)
    linesWithGaps := strings.Split(fileContent, "\n")
    lines := make([]string, 0)

    for _, rawline := range linesWithGaps {
        if(len(rawline) > 0) {
            lines = append(lines, rawline)        
        }
    }
    result := make([]*net.IPAddr, len(lines))

    for i, line := range lines {
        result[i], err = net.ResolveIPAddr("ip4:icmp", line)
        if err != nil {
            return nil, err
        }        
    }
    return result, nil
}

func main () {
    var cfgFilePath string
    var interval float64
    var alarminterval float64
    flag.StringVar(&cfgFilePath, "sites", "./sites.txt", "Path to file with pinged addresses")
    flag.Float64Var(&interval, "interval", 1.0, "Interval between sending requests")
    flag.Float64Var(&alarminterval, "alarmInterval", 60.0, "Internet problem alert interval (sec)")
    flag.Parse()

    //Парсим список пингуемых адресов
    addrs, err := parseConf(cfgFilePath)
    if err != nil {
        log.Println(err)    
    }

    //http пинговка
    HttpPingLoop(addrs, interval, alarminterval)
    //icmp пинговка
    /*
    for err == nil {
        err = icmpPing(addrs)        
        if err != nil {
            log.Println(err)  
        }  
    }
    */
}
