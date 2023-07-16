package nex_datastore

import (
	"github.com/PretendoNetwork/pokemon-rumble-world/database"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
)

func ChangeMeta(err error, client *nex.Client, callID uint32, param *datastore_types.DataStoreChangeMetaParam) {
	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)

	if err != nil {
		globals.Logger.Error(err.Error())
		rmcResponse.SetError(nex.Errors.DataStore.Unknown)
	} else {
		err = database.UpdateMetaBinaryByDataStoreChangeMetaParam(param)
		if err != nil {
			globals.Logger.Error(err.Error())
			rmcResponse.SetError(nex.Errors.DataStore.Unknown)
		}
	}

	if err == nil {
		rmcResponse.SetSuccess(datastore.MethodChangeMeta, nil)
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
