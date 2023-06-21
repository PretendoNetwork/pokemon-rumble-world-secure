package nex

import (
	secureconnection "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"
	"github.com/PretendoNetwork/pokemon-rumble-world-secure/globals"
)

func registerCommonProtocols() {
	_ = secureconnection.NewCommonSecureConnectionProtocol(globals.NEXServer)
}
