package nex_datastore

import (
	"github.com/PretendoNetwork/pokemon-rumble-world-secure/database"
	"github.com/PretendoNetwork/pokemon-rumble-world-secure/globals"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func ChangeMeta(err error, client *nex.Client, callID uint32, param *datastore.DataStoreChangeMetaParam) {
	// TODO - Check error
	_ = database.UpdateMetaBinaryByDataStoreChangeMetaParam(param)

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodChangeMeta, nil)

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
