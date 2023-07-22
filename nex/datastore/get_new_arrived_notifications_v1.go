package nex_datastore

import (
	"github.com/PretendoNetwork/pokemon-rumble-world/database"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
)

func GetNewArrivedNotificationsV1(err error, client *nex.Client, callID uint32, param *datastore_types.DataStoreGetNewArrivedNotificationsParam) {
	rmcResponse := nex.NewRMCResponse(0, callID)

	if err != nil {
		globals.Logger.Error(err.Error())
		rmcResponse.SetError(nex.Errors.DataStore.Unknown)
	} else {
		pResult := database.GetNotificationsByPIDAndParam(client.PID(), param)

		rmcResponseStream := nex.NewStreamOut(globals.HPPServer)

		rmcResponseStream.WriteListStructure(pResult)
		rmcResponseStream.WriteBool(false) // pHasNext. TODO - Handle this

		rmcResponseBody := rmcResponseStream.Bytes()

		rmcResponse.SetSuccess(datastore.MethodGetNewArrivedNotificationsV1, rmcResponseBody)
	}

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewHPPPacket(client, nil)

	responsePacket.SetPayload(rmcResponseBytes)

	globals.HPPServer.Send(responsePacket)
}
