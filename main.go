package main

import (
    "log"
    "flag"
    "io/ioutil"
    "strconv"
    "strings"
    "github.com/exsued/httpping"
    "os/exec"
    "fmt"
    "time"
    "os"
    "net"
    "encoding/gob"
)

var (
    alarmScriptPath string
    debug bool
    logDirPath string
    cfgFilePath string
    sessionServer string
    deviceName string
    packetPrefix = "name_pref"
    observedIfaceName = "eth0"

    interval float64
    pinger *httpping.HttpPinger
)

type tcpPacket struct {
        DeviceName string
        InnerAddrs []string
}
func LogFile(out string, dirpath string) {
    nowtime := time.Now()
    finalString := nowtime.Format("15:04:05\t") + out + "\n"
    fileName := dirpath + deviceName + ":" + nowtime.Format("2006-01-02") + ".txt"

    f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
    if _, err = f.WriteString(finalString); err != nil {
        log.Fatal(err)
    }
}
func parseConf(filePath string) ([]string, error) {
    bytesRead, err := ioutil.ReadFile(cfgFilePath)
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
        LogFile("Success. Returned: " + strconv.Itoa(httpStatus), logDirPath)
    }
}
func OnFailedReceive(err error) {
    log.Println("Failed." + err.Error())
    LogFile("Failed." + err.Error(), logDirPath)
}
func OnAlarm() {
    cmd, err := exec.Command("sudo", "/bin/sh", alarmScriptPath).Output()
    log.Println("Running alarm script: ", "sudo", "/bin/sh", string(cmd), alarmScriptPath)
    if err != nil {
        fmt.Println(err)
    }
    LogFile("Running alarm script: " + " sudo " + " /bin/sh\n" + string(cmd) + "\n", logDirPath)
}
func getInnerAddrs() (result []string) {
    ifaces, err := net.Interfaces()
    if err != nil {
        log.Println("Get ifaces err")
        LogFile("Get ifaces err", logDirPath)
    }
    for _, i := range ifaces {
        if i.Name != observedIfaceName {
            continue
        }
        addrs, err := i.Addrs()
        if err != nil {
        log.Println("Get addrs of iface",i.Name,"err")
        LogFile("Get addrs of iface " + i.Name + " err", logDirPath)
        }
        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
            case *net.IPNet:
                    ip = v.IP
            case *net.IPAddr:
                    ip = v.IP
            }
            print(ip.String())
            result = append(result, ip.String())
        }
    }
    return
}
func tcpClient() {
	conn, err := net.Dial("tcp", sessionServer)
    if err != nil {
        log.Println("dial err:", err)
        LogFile("dial err: " + err.Error(), logDirPath)
        return
    }
    innerAddrs:= getInnerAddrs()
    encoder := gob.NewEncoder(conn)

    for err == nil {
        packet := tcpPacket {deviceName, innerAddrs}
        err = encoder.Encode(packet)
        if err != nil {
            log.Println("send err:", err)
            LogFile("send err: " + err.Error(), logDirPath)
            conn.Close()
        }
        time.Sleep(time.Duration(interval) * time.Second)
    }
    log.Println("conn closed ", err)
    LogFile("conn closed " + err.Error(), logDirPath)
    if conn != nil {
        conn.Close()
    }
}
func main () {
    //vds1.proxinet.com
    var alarminterval float64
    var GetTimeout float64
    flag.StringVar(&sessionServer, "sessionServer", "vds1.proxicom.ru:1289", "address to long tcp session server")
    flag.StringVar(&cfgFilePath, "sites", "./sites.txt", "Path to file with pinged addresses")
    flag.StringVar(&alarmScriptPath, "onAlarm", "./alarm.sh", "Path to alarm script")
    flag.StringVar(&logDirPath, "log", "./logs/", "Path to log directory")
    flag.StringVar(&deviceName, "name", "proxicom_test", "Device name")
    flag.StringVar(&observedIfaceName, "iface", "eth0", "Observed iface name which use for check local addr")
    flag.Float64Var(&interval, "interval", 1.0, "Interval between sending requests (sec)")
    flag.Float64Var(&alarminterval, "alarmInterval", 60.0, "Internet problem alert interval (sec)")
    flag.Float64Var(&GetTimeout, "GetTimeout", 10.0, "HTTP GET Timeout (sec)")
    flag.BoolVar(&debug, "debug", false, "Set advanced output mode")
    flag.Parse()
    //log.Println(logDirPath)
    //log.Println(cfgFilePath)
    //log.Println(alarmScriptPath)
    //Парсим список пингуемых адресов
    addrs, err := parseConf(cfgFilePath)
    if err != nil {
        log.Fatalf(err.Error())
    }

    //http пинговка
    pinger = httpping.NewHttpPinger(addrs, interval, alarminterval, GetTimeout)
    pinger.OnRecv = OnReceive
    pinger.OnAlarm = OnAlarm
    pinger.OnFailedRecv = OnFailedReceive
    go pinger.Start()
    for ;; {
        tcpClient()
        time.Sleep(time.Duration(GetTimeout) * time.Second)
    }
}
