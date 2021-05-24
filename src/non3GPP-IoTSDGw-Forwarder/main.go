package main

import (
    "errors"
    "fmt"
    "github.com/urfave/cli"
    "log"
    "net"
    "os"
    "time"
    //"github.com/brocaar/lorawan"
)

type Config struct {
    ipv4 string
    port int
}

var (
    gwremoteAddr *net.UDPAddr
	config *Config
    )

var commandsCLi = []cli.Flag{
    cli.StringFlag{
        Name: "ipv4",
        Usage: "LoRaWan ipv4 binding address",

    },
    cli.IntFlag{
        Name: "port",
        Usage: "LoRaWan Network Server UDP binding port",
    },
}

/* A Simple function to verify error */
func checkError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}

func action(c * cli.Context) error{
    config = &Config{
        ipv4: c.String("ipv4"),
        port: c.Int("port"),
    }
    return nil
}

func startUplink(conngw *net.UDPConn, connbr1 *net.UDPConn) {

    upData := make([]byte,2048)
    for {
        n, remoteAddr, err := conngw.ReadFromUDP(upData)
        checkError(err)

        if n> 0 {

            fmt.Printf("Data %s was received from GW %s \n", upData, remoteAddr)

            _,err := connbr1.Write(upData[:n])
            checkError(err)
			fmt.Printf("Data %s was sent to bridge \n", upData)
        }
    }
}

func startDownlink(conngw *net.UDPConn, connbr1 *net.UDPConn ,  gwremoteAddr *net.UDPAddr) {

    downData := make([]byte,2048)
    for {

        n, remoteAddr, err := connbr1.ReadFromUDP(downData)
        checkError(err)
        if n> 0 {
            fmt.Printf("Data %s was received from BR! %s \n", downData, remoteAddr)
            _,err := conngw.WriteToUDP(downData[:n], gwremoteAddr)
            checkError(err)
			fmt.Printf("Data %s was sent to GW \n", downData)
        }
    }
}

func main() {

    app := cli.NewApp()
    app.Name = "IoTSDGW Forwarder"
    app.Usage = "Usage: -ipv4 {LoRAWAN Bridger} -port {UDP PORT}"
    app.Action = action
    app.Flags = commandsCLi
    if err := app.Run(os.Args); err != nil {
        log.Fatal("UE Run error: %v", err)
    }

    stringAddrGW := ":1680"
    localAddr, err := net.ResolveUDPAddr("udp", stringAddrGW)
    checkError(err)
 
    conngw, err := net.ListenUDP("udp", localAddr)
    checkError(err)

    if config.ipv4 == ""{
        checkError(errors.New("Ip LoRawan Bridge not set"))
    }
    if config.port == 0 {
        config.port = 1700 // Config set default port lorawan bridge
    }

    stringAddrBR := fmt.Sprintf("%s:%d", config.ipv4, config.port)
    fmt.Println(stringAddrBR)
    serverAddr,err := net.ResolveUDPAddr("udp", stringAddrBR)
    checkError(err)


    for{
        registerData := make([]byte,2048)
        n, remoteAddr, err := conngw.ReadFromUDP(registerData)
        checkError(err)
        if n > 0 {
            gwremoteAddr = remoteAddr
            fmt.Printf("Data %s was received from GW to register %s \n", registerData, gwremoteAddr)
            break
        }
    }


    connbr1, err := net.DialUDP("udp",nil, serverAddr)
    checkError(err)

    go startUplink(conngw, connbr1)

    go startDownlink(conngw, connbr1 , gwremoteAddr)

    for {
        time.Sleep(5)
    }
}