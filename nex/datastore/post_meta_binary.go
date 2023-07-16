package nex_datastore

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/pokemon-rumble-world/database"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"
)

func PostMetaBinary(err error, client *nex.Client, callID uint32, param *datastore_types.DataStorePreparePostParam) {
	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)

	if err != nil {
		globals.Logger.Error(err.Error())
		rmcResponse.SetError(nex.Errors.DataStore.Unknown)
	}

	metaBinary := database.GetMetaInfoByOwnerPID(client.PID())

	if metaBinary.DataID != 0 {
		// * Meta binary already exists
		if param.PersistenceInitParam.DeleteLastObject {
			// * Delete existing object before uploading new one
			// TODO - Check error
			_ = database.DeleteMetaBinaryByDataID(metaBinary.DataID)
		}
	}

	dataID, err := database.InsertMetaBinaryByDataStorePreparePostParamWithOwnerPID(param, client.PID())
	if err != nil {
		globals.Logger.Error(err.Error())
		rmcResponse.SetError(nex.Errors.DataStore.Unknown)
	}

	if err == nil {
		rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

		rmcResponseStream.WriteUInt64LE(uint64(dataID))

		rmcResponseBody := rmcResponseStream.Bytes()

		rmcResponse.SetSuccess(datastore.MethodPreparePostObject, rmcResponseBody)
	}

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.SecureServer.Send(responsePacket)
}
