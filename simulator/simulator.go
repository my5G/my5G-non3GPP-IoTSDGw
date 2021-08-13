package main

import (
    "fmt"
    "github.com/my5G/my5G-non3GPP-IoTSDGw/simulator/context"
    "github.com/my5G/my5G-non3GPP-IoTSDGw/simulator/udp_server"
    "github.com/urfave/cli"
    "log"
    "os"
    "sync"
    "time"
)

const (
    ChannelsUpLinksNumber = 8
)

var (
    InfoLogger *log.Logger
    ErrorLogger *log.Logger
    WarningLogger *log.Logger
    channelFlag int
    numPackets int
    config *Config
    wg sync.WaitGroup
    now = time.Now()
)

var commandsCLi = []cli.Flag{
    cli.StringFlag{
        Name: "ipv4",
        Usage: "IOTSDW Forwarder ipv4 binding address",

    },
    cli.IntFlag{
        Name: "port",
        Usage: "IOTSDW Forwarder Server UDP binding port",
    },
    cli.IntFlag{
        Name: "dev",
        Usage: "IOTSDW Forwarder number of devices",
    },
    cli.IntFlag{
        Name: "packets",
        Usage: "IOTSDW Forwarder lenth packets per device",
    },
}

type Config struct {
    ipv4 string
    port int
    numDevices int
    packetPerDevices int
}

func init (){
    InfoLogger = log.New(os.Stdout ,"Info: ", log.Ldate|log.Ltime|log.Lshortfile)
    ErrorLogger = log.New(os.Stdout ,"Warning: ", log.Ldate|log.Ltime|log.Lshortfile)
    WarningLogger = log.New(os.Stdout ,"Error: ", log.Ldate|log.Ltime|log.Lshortfile)
}

/* A Simple function to verify error */
func checkError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}

func Initialize(){

    //make new  Socker Port
    ok :=context.DevicesContext_Self().ConfigSocketUDPAddr(config.ipv4,config.port)
    if !ok {
        ErrorLogger.Println("Socket Bind Error")
        return
    }

    InfoLogger.Println("Initializing UDP Server Network")
    go udp_server.Run()

    InfoLogger.Println("Create Devices for simulator Network")
    context.CreateDevicesForSimulate(config.numDevices) // make 100 0 devices

    InfoLogger.Println("Run Simulator")
    run_simulator()

}

func runDevices(i int, w *sync.WaitGroup){

    defer w.Done()

    device, ok := context.DevicesContext_Self().DeviceLoad(uint16(i))
    if !ok {
        ErrorLogger.Println("Device id not Found")
        os.Exit(0)
    }

    for flag := 0; flag < config.packetPerDevices; flag++ {
        newPayload, ok := device.MakeNewPayload()
        if !ok {
            panic("Error msg encode Device #{flag} and packet #{flag}")
        }

        fmt.Printf("%s\n\n\n", newPayload)

        channelId := udp_server.CycleChannel()

        device.FsmState = context.FSM_SEND
        InfoLogger.Printf("Send message device %d packet %d ",i , flag )
        device.Packet_tx++
        udp_server.SendChannelMessage(newPayload, channelId)
        device.FsmState = context.FSM_WAIT

        /* Timeout */
        done := make(chan bool)
        go func(){
            for{
                if device.FsmState == context.FSM_RECV{
                    done <- true
                    return
                }
                if  false == <-done {
                    return
                }
            }
        }()

        select {
        case <-done:
            device.FsmState = context.FSM_IDLE
        case <-time.After(5):
            done <- false
        }

    }// End of For
}

func run_simulator(){
    //total_cycle := numPacketPerDev * numDev


    for i := 1; i <= config.numDevices; i++ {
        wg.Add(1)
        go runDevices(i, &wg)
    }
}

func action(c * cli.Context) error{
    config = &Config{
        ipv4: c.String("ipv4"),
        port: c.Int("port"),
        numDevices: c.Int("dev"),
        packetPerDevices: c.Int("packets"),
    }
    return nil
}

func main(){
    app := cli.NewApp()
    app.Name = "IoTSDGW LoRa Simulator"
    app.Usage = "Usage: -ipv4 {IOTSdw Forwarder} -port {UDP port}"
    app.Action = action
    app.Flags = commandsCLi
    if err := app.Run(os.Args); err != nil {
        log.Fatal("UE Run error: %v", err)
    }

    if config.ipv4 == ""{
       config.ipv4="127.0.0.1"
    }

    if config.port == 0 {
        config.port = 1680 // Config set default port lorawan bridge
    }

    if config.numDevices  <= 0 {
        config.numDevices = 1 // Config set default port lorawan bridge
    }

    if config.packetPerDevices <=  0 {
        config.packetPerDevices = 1// Config set default port lorawan bridge
    }

    Initialize()
    wg.Wait()
}