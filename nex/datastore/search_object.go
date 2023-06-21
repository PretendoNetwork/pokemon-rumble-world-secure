package nex_datastore

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	"github.com/PretendoNetwork/pokemon-rumble-world-secure/database"
	"github.com/PretendoNetwork/pokemon-rumble-world-secure/globals"
	"github.com/PretendoNetwork/pokemon-rumble-world-secure/types"
)

func SearchObject(err error, client *nex.Client, callID uint32, param *datastore.DataStoreSearchParam) {
	metaBinaries := make([]*types.MetaBinary, 0)

	if param.SearchTarget == 10 { // * Search for meta binary of this client
		metaBinary := database.GetMetaInfoByOwnerPID(client.PID())
		if metaBinary.DataID != 0 {
			metaBinaries = append(metaBinaries, metaBinary)
		}
	} else if len(param.DataTypes) > 0 { // * The data type is given inside the DataTypes param
		metaBinaries = database.GetMetaInfosByDataStoreSearchParam(param)
	} else {
		metaBinaries = database.GetAllMetaInfosByDataStoreSearchParam(param)
	}

	pSearchResult := datastore.NewDataStoreSearchResult()

	pSearchResult.TotalCount = uint32(len(metaBinaries))
	pSearchResult.Result = make([]*datastore.DataStoreMetaInfo, 0, len(metaBinaries))

	for i := 0; i < len(metaBinaries); i++ {
		metaBinary := metaBinaries[i]
		result := datastore.NewDataStoreMetaInfo()

		result.DataID = uint64(metaBinary.DataID)
		result.OwnerID = metaBinary.OwnerPID
		result.Size = 0
		result.Name = metaBinary.Name
		result.DataType = metaBinary.DataType
		result.Permission = datastore.NewDataStorePermission()
		result.Permission.Permission = metaBinary.Permission
		result.Permission.RecipientIds = make([]uint32, 0)
		result.DelPermission = datastore.NewDataStorePermission()
		result.DelPermission.Permission = metaBinary.DeletePermission
		result.DelPermission.RecipientIds = make([]uint32, 0)
		result.CreatedTime = metaBinary.CreationTime
		result.UpdatedTime = metaBinary.UpdatedTime
		result.Period = metaBinary.Period
		result.Status = 0      // TODO - Figure this out
		result.ReferredCnt = 0 // TODO - Figure this out
		result.ReferDataID = 0 // TODO - Figure this out
		result.Flag = metaBinary.Flag
		result.ReferredTime = metaBinary.ReferredTime
		result.ExpireTime = metaBinary.ExpireTime
		result.Tags = metaBinary.Tags
		result.Ratings = make([]*datastore.DataStoreRatingInfoWithSlot, 0)

		pSearchResult.Result = append(pSearchResult.Result, result)
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteStructure(pSearchResult)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodSearchObject, rmcResponseBody)

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