package nex_datastore

import (
	"github.com/PretendoNetwork/pokemon-rumble-world/database"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
)

func GetNewArrivedNotificationsV1(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStoreGetNewArrivedNotificationsParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.InvalidArgument, err.Error())
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	pResult := database.GetNotificationsByPIDAndParam(uint32(connection.PID()), param)
	pHasNext := types.NewBool(false) // TODO - Handle this

	rmcResponseStream := nex.NewByteStreamOut(globals.HPPServer.LibraryVersions(), globals.HPPServer.ByteStreamSettings())

	pResult.WriteTo(rmcResponseStream)
	pHasNext.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodGetNewArrivedNotificationsV1
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
