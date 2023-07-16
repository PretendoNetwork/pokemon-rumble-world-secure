package nex

import (
	secureconnection "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"
)

func registerCommonSecureServerProtocols() {
	_ = secureconnection.NewCommonSecureConnectionProtocol(globals.SecureServer)
}
