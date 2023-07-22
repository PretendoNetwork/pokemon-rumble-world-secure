package nex

import (
	"fmt"
	"os"

	"github.com/PretendoNetwork/pokemon-rumble-world/globals"

	nex "github.com/PretendoNetwork/nex-go"
)

func StartHPPServer() {
	globals.HPPServer = nex.NewServer()
	globals.HPPServer.SetDefaultNEXVersion(&nex.NEXVersion{
		Major: 2,
		Minor: 4,
		Patch: 1,
	})
	globals.HPPServer.SetAccessKey("844f1d0c")
	globals.HPPServer.SetPasswordFromPIDFunction(globals.PasswordFromPID)

	globals.HPPServer.On("Data", func(packet *nex.HPPPacket) {
		request := packet.RMCRequest()

		fmt.Println("== Pokémon Rumble World - HPP ==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("======================")
	})

	registerHPPServerNEXProtocols()

	globals.HPPServer.HPPListen(fmt.Sprintf("%s:%s", os.Getenv("PN_PRW_HPP_SERVER_HOST"), os.Getenv("PN_PRW_HPP_SERVER_PORT")))
}
