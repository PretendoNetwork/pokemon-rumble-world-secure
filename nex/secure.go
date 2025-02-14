package nex

import (
	"fmt"
	"os"
	"strconv"

	"github.com/PretendoNetwork/pokemon-rumble-world/globals"

	nex "github.com/PretendoNetwork/nex-go/v2"
)

func StartSecureServer() {
	globals.SecureServer = nex.NewPRUDPServer()

	globals.SecureEndpoint = nex.NewPRUDPEndPoint(1)
	globals.SecureEndpoint.IsSecureEndPoint = true
	globals.SecureEndpoint.ServerAccount = globals.SecureServerAccount
	globals.SecureEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.SecureEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.SecureServer.BindPRUDPEndPoint(globals.SecureEndpoint)

	globals.SecureServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(3, 8, 2))
	globals.SecureServer.ByteStreamSettings.UseStructureHeader = true
	globals.SecureServer.AccessKey = "844f1d0c"

	globals.SecureEndpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("== Pokémon Rumble World - Secure ==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID)
		fmt.Printf("Method ID: %#v\n", request.MethodID)
		fmt.Println("======================")
	})

	globals.SecureEndpoint.OnError(func(err *nex.Error) {
		globals.Logger.Error(err.Error())
	})

	registerCommonSecureServerProtocols()
	registerSecureServerNEXProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_PRW_SECURE_SERVER_PORT"))

	globals.SecureServer.Listen(port)
}
