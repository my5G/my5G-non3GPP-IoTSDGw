package main

import (
    "github.com/pkg/errors"
    "github.com/urfave/cli"
    "log"
    "net"
    "os"
    "sync"
    "time"
)

const (
    defaultForwarderPort  int = 1680
    defaultPortSendGW = 1690
)

type Config struct {
    ipv4Local string
    ipv4NetServer string
    port int
}

var (
    gwremoteAddr *net.UDPAddr
    config *Config
    wg            sync.WaitGroup
    metrics       Metrics
    uplinkCounter int
    downlinkCounter int
    starTimeDownload time.Time
    starTimeUpload time.Time
)

var commandsCLi = []cli.Flag{
    cli.StringFlag{
        Name: "ipv4Local",
        Usage: "LoRaWan ipv4 binding address",

    },

    cli.StringFlag{
        Name: "ipv4NetServer",
        Usage: "LoRaWan ipv4 netserver address",

    },

    cli.IntFlag{
        Name: "port",
        Usage: "LoRaWan Network Server UDP binding port",
    },
}

func init(){
    metrics = Metrics{}
    metrics.Init()
}

func action(c * cli.Context) error{
    config = &Config{
        ipv4Local: c.String("ipv4Local"),
        ipv4NetServer: c.String("ipv4NetServer"),
        port: c.Int("port"),
    }
    return nil
}

func configBindAddr(listenAddrPort *net.UDPAddr) func(string, int ) ( *net.UDPConn, error) {
    return func( ipAddress string, port int ) ( *net.UDPConn, error) {

        listenAddrPort.Port = port
        validIpAddress := net.ParseIP(ipAddress)

        if validIpAddress != nil {
            listenAddrPort.IP = validIpAddress
            listener, err := net.ListenUDP("udp", listenAddrPort )
            if err != nil {
                return nil , err
            }
            return listener, nil
        }
        return nil, errors.New("Not IPAddress valid")
    }
}

func configBindDialupAddr(addrUDP *net.UDPAddr) func(string, int ) ( *net.UDPConn, error) {
    return func( ipAddress string, port int ) ( *net.UDPConn, error) {

        addrUDP.Port = port
        validIpAddress := net.ParseIP(ipAddress)

        if validIpAddress != nil {
            addrUDP.IP = validIpAddress
            dialUp, err := net.DialUDP("udp", nil, addrUDP)
            if err != nil {
                return nil , err
            }
            return dialUp, nil
        }
        return nil, errors.New("Not IPAddress valid")
    }
}

func dispatchTraffic(buff []byte, socket *net.UDPConn) {
    _, err := socket.Write(buff)

    if err != nil {
        log.Fatalf("wrtite data to UDP Socket failed %v", err)
    }
}

func StartUplink(fromGateway  *net.UDPConn, toNetServer *net.UDPConn) {
    defer wg.Done()

    buff  := make([]byte, 65535)
    starTimeUpload = now()
    for {
        n, _, err := fromGateway.ReadFromUDP(buff)
        if err != nil {
            log.Fatalf("[Gateway] Read from UDP failed %v", err)
        } else{

            uplinkCounter++
            actualElipsedTime := now().Sub(starTimeUpload)
            metrics.Store( StoreInfo("uplink", uplinkCounter, len(buff[:n]), actualElipsedTime.Seconds() ) )

            if n > 0 {

                go dispatchTraffic(buff[:n], toNetServer)
            }
        }
    }
}

func StartDownlink(fromBridge *net.UDPConn, toGateway *net.UDPConn ) {
    defer wg.Done()

    buff := make([]byte,65535)
    starTimeDownload = now()

    for {
        n, _ , err := fromBridge.ReadFromUDP(buff)
        if  err != nil {
            log.Fatalf("[Bridge] Read from UDP failed %v", err)
        } else {

            downlinkCounter++
            actualElipsedTime := now().Sub(starTimeDownload)

            //value := bytes.NewReader( buff[:n][1:3] )
            //u := binary.BigEndian.Uint16(value)
            //var u uint16
            //binary.Read(value, binary.BigEndian, &u)
            //fmt.Printf("Message : %d \n ", u)

            metrics.Store( StoreInfo("downlink", downlinkCounter, len(buff[:n]), actualElipsedTime.Seconds() ) )

            if n > 0 {
                go dispatchTraffic(buff[:n], toGateway)
            }
        }
    }
}

func Run (){

    UpForwarderConnection := new(net.UDPAddr)
    listenerPort1680, err := configBindAddr(UpForwarderConnection)(config.ipv4Local, defaultForwarderPort )
    if err != nil {
        log.Fatalf("Error Listener port %v ", err)
        return
    }

    BridgeAddCon := new(net.UDPAddr)
    dialUpBridge , err := configBindDialupAddr(BridgeAddCon)(config.ipv4NetServer, config.port )
    if err != nil {
        log.Fatalf("Error dialUP port %v ", err)
        return
    }

    GatewayConn := new(net.UDPAddr)
    dialUpGateway, err := configBindDialupAddr(GatewayConn)(config.ipv4Local, defaultPortSendGW )
    if err != nil {
        log.Fatalf("Error dialUP port %v ", err)
        return
    }
    wg.Add(2)
    go StartUplink(listenerPort1680, dialUpBridge)
    go StartDownlink(dialUpBridge, dialUpGateway)
    wg.Wait()
}

func main() {

    defer metrics.Close()

    app := cli.NewApp()
    app.Name = "IoTSDGW Forwarder"
    app.Usage = "Usage: -ipv4Local {LoRAWAN Gateway} -ipv4NetServer {NetServer}  -port {UDP PORT}"
    app.Action = action
    app.Flags = commandsCLi
    if err := app.Run(os.Args); err != nil {
        log.Fatal("UE Run error: %v", err)
    }
    if config.ipv4Local == ""{
        config.ipv4Local = "127.0.0.1"
    }

    if config.ipv4NetServer == "" {
        config.ipv4NetServer = "127.0.0.1"
    }

    if config.port == 0 {
        config.port = 1700 // Config set default port lorawan bridge
    }
    Run()
}

var now = time.Now