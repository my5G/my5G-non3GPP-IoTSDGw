package context

import (
	"encoding/binary"
	"encoding/hex"
)

var idDev uint16

const(
	FSM_IDLE = iota
	FSM_SEND
	FSM_WAIT
	FSM_RECV
)

const (
	PacketType string = "02"
	Payload string = "0000000000000000017b227278706b223a5b7b2274696d65223a22323032312d30322d31375430383a30383a33302d3033" +
		"3a3030222c22746d6d73223a393232333337323839302c22746d7374223a393232333337322c226368616e223a302c2272666368223a312" +
		"c2266726571223a3931362e382c2273746174223a312c226d6f6475223a224c4f5241222c2264617472223a225346374257313235222c22" +
		"636f6472223a22342f35222c2272737369223a2d35372c226c736e72223a372c2273697a65223a31342c2264617461223a2251414541414" +
		"1414141674143347841393846383d227d5d7d"
)

//
// Pega o json avbri o data
/// Converte do base 64 para hex
// Alterar os bytes os bites



func init(){
	idDev = 0
}

func IdEncoding(i uint16) (string){
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return hex.EncodeToString(b)
}

func Incremment() (uint16){
	idDev = idDev + 1
	return idDev
}