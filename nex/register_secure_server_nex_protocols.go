package nex

import (
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"
	nex_datastore "github.com/PretendoNetwork/pokemon-rumble-world/nex/datastore"
)

func registerSecureServerNEXProtocols() {
	dataStoreProtocol := datastore.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(dataStoreProtocol)

	dataStoreProtocol.SetHandlerSearchObject(nex_datastore.SearchObject)
	dataStoreProtocol.SetHandlerPostMetaBinary(nex_datastore.PostMetaBinary)
	dataStoreProtocol.SetHandlerChangeMeta(nex_datastore.ChangeMeta)
	dataStoreProtocol.SetHandlerGetMetas(nex_datastore.GetMetas)
}
