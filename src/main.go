package main

import (
	"fmt"
	"net"
    "os"
    "time"
	//"github.com/brocaar/lorawan"
)


 var  gwremoteAddr *net.UDPAddr



/* A Simple function to verify error */
func checkError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
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


    stringAddrGW := ":1680"
    localAddr, err := net.ResolveUDPAddr("udp", stringAddrGW)
    checkError(err)
 
    conngw, err := net.ListenUDP("udp", localAddr)
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

    stringAddrBR := "31.220.52.116:1700"
    serverAddr,err := net.ResolveUDPAddr("udp", stringAddrBR)
    checkError(err)

    connbr1, err := net.DialUDP("udp",nil, serverAddr)
    checkError(err)

    
    

    go startUplink(conngw, connbr1)

    go startDownlink(conngw, connbr1 , gwremoteAddr)



    
    for {
        time.Sleep(5)
    }

}

