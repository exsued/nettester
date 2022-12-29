package main

import (
    "log"
    "flag"
    "io/ioutil"
    "strings"
    "github.com/exsued/httpping"
)

func parseConf(filePath string) ([]string, error) {
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
    result := make([]string, len(lines))

    for i, line := range lines {
        result[i] = line      
    }
    return result, nil
}

func gimba() {
    log.Println("Gocha")
}

func banga() {
    log.Println("!Alarm!")
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
    pinger := httpping.NewHttpPinger(addrs, interval, alarminterval)
    pinger.OnAlarm = banga
    pinger.OnRecv = gimba
    pinger.Start()
}
