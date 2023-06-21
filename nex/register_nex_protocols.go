package nex

import (
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	"github.com/PretendoNetwork/pokemon-rumble-world-secure/globals"
	nex_datastore "github.com/PretendoNetwork/pokemon-rumble-world-secure/nex/datastore"
)

func registerNEXProtocols() {
	dataStoreProtocol := datastore.NewDataStoreProtocol(globals.NEXServer)

	dataStoreProtocol.SearchObject(nex_datastore.SearchObject)
	dataStoreProtocol.PostMetaBinary(nex_datastore.PostMetaBinary)
	dataStoreProtocol.ChangeMeta(nex_datastore.ChangeMeta)
	dataStoreProtocol.GetMetas(nex_datastore.GetMetas)
}
