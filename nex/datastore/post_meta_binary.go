package nex_datastore

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/pokemon-rumble-world-secure/database"
	"github.com/PretendoNetwork/pokemon-rumble-world-secure/globals"
)

func PostMetaBinary(err error, client *nex.Client, callID uint32, param *datastore_types.DataStorePreparePostParam) {
	metaBinary := database.GetMetaInfoByOwnerPID(client.PID())

	if metaBinary.DataID != 0 {
		// * Meta binary already exists
		if param.PersistenceInitParam.DeleteLastObject {
			// * Delete existing object before uploading new one
			// TODO - Check error
			_ = database.DeleteMetaBinaryByDataID(metaBinary.DataID)
		}
	}

	dataID := database.InsertMetaBinaryByDataStorePreparePostParamWithOwnerPID(param, client.PID())

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteUInt64LE(uint64(dataID))

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodPreparePostObject, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)
}
