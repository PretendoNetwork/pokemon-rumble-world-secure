package nex

import (
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"
	nex_datastore "github.com/PretendoNetwork/pokemon-rumble-world/nex/datastore"
)

func registerHPPServerNEXProtocols() {
	dataStoreProtocol := datastore.NewProtocol()
	globals.HPPServer.RegisterServiceProtocol(dataStoreProtocol)

	dataStoreProtocol.SetHandlerPrepareGetObjectV1(nex_datastore.PrepareGetObjectV1)
	dataStoreProtocol.SetHandlerPreparePostObjectV1(nex_datastore.PreparePostObjectV1)
	dataStoreProtocol.SetHandlerCompletePostObjectV1(nex_datastore.CompletePostObjectV1)
	dataStoreProtocol.SetHandlerGetNotificationURL(nex_datastore.GetNotificationURL)
	dataStoreProtocol.SetHandlerGetNewArrivedNotificationsV1(nex_datastore.GetNewArrivedNotificationsV1)
	dataStoreProtocol.SetHandlerGetSpecificMetaV1(nex_datastore.GetSpecificMetaV1)
}
