package nex

import (
	"fmt"
	"os"
	"strconv"

	"github.com/PretendoNetwork/pokemon-rumble-world/globals"

	nex "github.com/PretendoNetwork/nex-go/v2"
)

func StartHPPServer() {
	globals.HPPServer = nex.NewHPPServer()
	globals.HPPServer.LibraryVersions().SetDefault(nex.NewLibraryVersion(2, 4, 1))
	globals.HPPServer.SetAccessKey("844f1d0c")
	globals.HPPServer.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.HPPServer.AccountDetailsByUsername = globals.AccountDetailsByUsername

	globals.HPPServer.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("== Pokémon Rumble World - HPP ==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID)
		fmt.Printf("Method ID: %#v\n", request.MethodID)
		fmt.Println("======================")
	})

	registerHPPServerNEXProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_PRW_HPP_SERVER_PORT"))

	globals.HPPServer.Listen(port)
}
