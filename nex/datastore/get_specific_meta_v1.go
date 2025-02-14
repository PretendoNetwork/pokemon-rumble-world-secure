package nex_datastore

import (
	"github.com/PretendoNetwork/pokemon-rumble-world/database"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"

	"github.com/PretendoNetwork/nex-go/v2"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
)

func GetSpecificMetaV1(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStoreGetSpecificMetaParamV1) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.InvalidArgument, err.Error())
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	rawDataIDs := make([]uint32, len(param.DataIDs))
	for i, dataID := range param.DataIDs {
		rawDataIDs[i] = uint32(dataID)
	}

	pMetaInfos := database.GetNotificationMetasByDataIDs(rawDataIDs)

	rmcResponseStream := nex.NewByteStreamOut(globals.HPPServer.LibraryVersions(), globals.HPPServer.ByteStreamSettings())

	pMetaInfos.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodGetSpecificMetaV1
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
