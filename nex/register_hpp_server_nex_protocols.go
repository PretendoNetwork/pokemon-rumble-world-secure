package nex

import (
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"
	nex_datastore "github.com/PretendoNetwork/pokemon-rumble-world/nex/datastore"
)

func registerHPPServerNEXProtocols() {
	dataStoreProtocol := datastore.NewDataStoreProtocol(globals.HPPServer)

	dataStoreProtocol.PrepareGetObjectV1(nex_datastore.PrepareGetObjectV1)
	dataStoreProtocol.PreparePostObjectV1(nex_datastore.PreparePostObjectV1)
	dataStoreProtocol.CompletePostObjectV1(nex_datastore.CompletePostObjectV1)
	dataStoreProtocol.GetNotificationURL(nex_datastore.GetNotificationURL)
	dataStoreProtocol.GetNewArrivedNotificationsV1(nex_datastore.GetNewArrivedNotificationsV1)
	dataStoreProtocol.GetSpecificMetaV1(nex_datastore.GetSpecificMetaV1)
}
