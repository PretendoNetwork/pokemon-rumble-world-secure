package nex

import (
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"
	nex_datastore "github.com/PretendoNetwork/pokemon-rumble-world/nex/datastore"
)

func registerSecureServerNEXProtocols() {
	dataStoreProtocol := datastore.NewDataStoreProtocol(globals.SecureServer)

	dataStoreProtocol.SearchObject(nex_datastore.SearchObject)
	dataStoreProtocol.PostMetaBinary(nex_datastore.PostMetaBinary)
	dataStoreProtocol.ChangeMeta(nex_datastore.ChangeMeta)
	dataStoreProtocol.GetMetas(nex_datastore.GetMetas)
}
