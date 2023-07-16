package nex_datastore

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/pokemon-rumble-world/database"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"
	"github.com/PretendoNetwork/pokemon-rumble-world/types"
)

func SearchObject(err error, client *nex.Client, callID uint32, param *datastore_types.DataStoreSearchParam) {
	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)

	if err != nil {
		globals.Logger.Error(err.Error())
		rmcResponse.SetError(nex.Errors.DataStore.Unknown)
	}

	if err == nil {
		metaBinaries := make([]*types.MetaBinary, 0)
		var totalCount uint32

		if param.SearchTarget == 10 { // * Search for meta binary of this client
			metaBinary := database.GetMetaInfoByOwnerPID(client.PID())
			if metaBinary.DataID != 0 {
				metaBinaries = append(metaBinaries, metaBinary)
			}
			totalCount = uint32(len(metaBinaries))
		} else if len(param.DataTypes) > 0 { // * The data type is given inside the DataTypes param
			metaBinaries = database.GetMetaInfosByDataStoreSearchParam(param)
			totalCount = database.GetTotalMetaInfosByDataTypes(param.DataTypes)
		} else if len(param.OwnerIDs) == 0 { // * Ignore unknown request for PID = 2 (Rendez-Vous)
			metaBinaries = database.GetAllMetaInfosByDataStoreSearchParam(param)
			totalCount = database.GetTotalMetaInfos()
		}

		pSearchResult := datastore_types.NewDataStoreSearchResult()

		pSearchResult.TotalCount = totalCount

		if totalCount > uint32(len(metaBinaries)) {
			pSearchResult.TotalCountType = 1 // * Not all results are returned
		}

		pSearchResult.Result = make([]*datastore_types.DataStoreMetaInfo, 0, len(metaBinaries))

		for i := 0; i < len(metaBinaries); i++ {
			metaBinary := metaBinaries[i]
			result := datastore_types.NewDataStoreMetaInfo()

			result.DataID = uint64(metaBinary.DataID)
			result.OwnerID = metaBinary.OwnerPID
			result.Size = 0
			result.Name = metaBinary.Name
			result.DataType = metaBinary.DataType
			result.Permission = datastore_types.NewDataStorePermission()
			result.Permission.Permission = metaBinary.Permission
			result.Permission.RecipientIDs = make([]uint32, 0)
			result.DelPermission = datastore_types.NewDataStorePermission()
			result.DelPermission.Permission = metaBinary.DeletePermission
			result.DelPermission.RecipientIDs = make([]uint32, 0)
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
			result.Ratings = make([]*datastore_types.DataStoreRatingInfoWithSlot, 0)

			pSearchResult.Result = append(pSearchResult.Result, result)
		}

		rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

		rmcResponseStream.WriteStructure(pSearchResult)

		rmcResponseBody := rmcResponseStream.Bytes()

		rmcResponse.SetSuccess(datastore.MethodSearchObject, rmcResponseBody)
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
