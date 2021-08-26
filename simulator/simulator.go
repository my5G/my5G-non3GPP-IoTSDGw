package main

import (
    "encoding/hex"
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
    now = time.Now
)

var commandsCLi = []cli.Flag{
    cli.StringFlag{
        Name: "ipv4",
        Usage: "IOTSDW Forwarder ipv4 binding address",

    },
    cli.IntFlag{
        Name: "portUp",
        Usage: "IOTSDW Forwarder Server UDP binding port for upLink",
    },

    cli.IntFlag{
        Name: "portDown",
        Usage: "IOTSDW Forwarder Server UDP binding port for downLink",
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
    portUp int
    portDown int
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

    context.DevicesContext_Self().Gateway.MAC = "0000000000000001"

    //make new  Socker Port
    ok :=context.DevicesContext_Self().ConfigSocketUDPAddr(config.ipv4,config.portDown)
    if !ok {
        ErrorLogger.Println("Socket Bind Downlink  Error")
        return
    }

    //make new  Socker Port
    ok = context.DevicesContext_Self().ConfigUplink(config.ipv4,config.portUp)
    if !ok {
        ErrorLogger.Println("Socket Bind Up  Error")
        return
    }

    InfoLogger.Println("Initializing UDP Server Network")
    go udp_server.Run()

    InfoLogger.Println("Create Devices for simulator Network")
    context.CreateDevicesForSimulate(config.numDevices) // make 1000 devices

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

    str := hex.EncodeToString([]byte{device.DevAddr[0],device.DevAddr[1],device.DevAddr[2],device.DevAddr[3]} )
    InfoLogger.Printf("************ devId %s DevADDr %s \n\n",device.GetDevID() , str )

    for flag := 0; flag < config.packetPerDevices; flag++ {


        device.SetMessagePayload(fmt.Sprintf("hello dev %d", device.DevId))
        phyLoRaPayload, ok := device.Marshall()
        if !ok {
            panic("Error msg encode Device #{flag} and packet #{flag}")
            return
        }

       // fmt.Printf("%s\n\n\n", phyLoRaPayload)

        channelId := udp_server.CycleChannel()
        device.FsmState = context.FSM_SEND
        InfoLogger.Printf("Send message device %d packet %d ",i , flag )
        device.Packet_tx++
        device.Start = now()

        udp_server.SendChannelMessage(phyLoRaPayload, device.GetDevID(),  channelId)

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
        case <-time.After(1 * time.Minute):
            device.FsmState = context.FSM_IDLE
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
        portUp: c.Int("portUp"),
        portDown: c.Int("portDown"),
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

    if config.portUp == 0 {
        config.portUp = 1680 // Config set default port lorawan bridge
    }

    if config.portDown == 0 {
        config.portDown = 1690 // Config set default port lorawan bridge
    }

    if config.numDevices  <= 0 {
        config.numDevices = 1000 // Config set default port lorawan bridge
    }

    if config.packetPerDevices <=  0 {
        config.packetPerDevices = 100// Config set default port lorawan bridge
    }

    Initialize()
    wg.Wait()
}