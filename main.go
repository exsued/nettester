package main

import (
    "log"
    "flag"
    "io/ioutil"
    "strings"
    "strconv"
    "github.com/exsued/httpping"
    "os/exec"
    "fmt"
)

var (
    alarmScriptPath string
    debug bool
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

func OnReceive(httpStatus int) {
    if debug {
        log.Println("Success. Returned: " + strconv.Itoa(httpStatus))
    }
}

func OnFailedReceive(err error) {
    log.Println("Failed." + err.Error())
}
func OnAlarm() {
    cmd, err := exec.Command("/bin/sh", alarmScriptPath).Output()
    log.Println("Running alarm script")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(string(cmd))
}

func main () {
    var cfgFilePath string
    var interval float64
    var alarminterval float64
    var GetTimeout float64
    flag.StringVar(&cfgFilePath, "sites", "./sites.txt", "Path to file with pinged addresses")
    flag.StringVar(&cfgFilePath, "onAlarm", "./alarm.sh", "Path to alarm script")
    flag.Float64Var(&interval, "interval", 1.0, "Interval between sending requests (sec)")
    flag.Float64Var(&alarminterval, "alarmInterval", 60.0, "Internet problem alert interval (sec)")
    flag.Float64Var(&GetTimeout, "GetTimeout", 10.0, "HTTP GET Timeout (sec)")
    flag.Parse()

    //Парсим список пингуемых адресов
    addrs, err := parseConf(cfgFilePath)
    if err != nil {
        log.Println(err)    
    }

    //http пинговка
    pinger := httpping.NewHttpPinger(addrs, interval, alarminterval, GetTimeout)
    pinger.OnRecv = OnReceive
    pinger.OnAlarm = OnAlarm
    pinger.OnFailedRecv = OnFailedReceive
    pinger.Start()
}
