package nex_datastore

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/pokemon-rumble-world/database"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"
)

func GetMetas(err error, client *nex.Client, callID uint32, dataIDs []uint64, param *datastore_types.DataStoreGetMetaParam) {
	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)

	if err != nil {
		globals.Logger.Error(err.Error())
		rmcResponse.SetError(nex.Errors.DataStore.Unknown)
	}

	if err == nil {
		metaBinaries := database.GetMetaBinariesByDataIDs(dataIDs)

		pMetaInfo := make([]*datastore_types.DataStoreMetaInfo, 0, len(metaBinaries))
		pResults := make([]*nex.Result, 0, len(metaBinaries))

		for i := 0; i < len(metaBinaries); i++ {
			metaBinary := metaBinaries[i]
			metaInfo := datastore_types.NewDataStoreMetaInfo()

			metaInfo.DataID = uint64(metaBinary.DataID)
			metaInfo.OwnerID = metaBinary.OwnerPID
			metaInfo.Size = 0
			metaInfo.Name = metaBinary.Name
			metaInfo.DataType = metaBinary.DataType
			metaInfo.MetaBinary = metaBinary.Buffer
			metaInfo.Permission = datastore_types.NewDataStorePermission()
			metaInfo.Permission.Permission = metaBinary.Permission
			metaInfo.Permission.RecipientIDs = make([]uint32, 0)
			metaInfo.DelPermission = datastore_types.NewDataStorePermission()
			metaInfo.DelPermission.Permission = metaBinary.DeletePermission
			metaInfo.DelPermission.RecipientIDs = make([]uint32, 0)
			metaInfo.CreatedTime = metaBinary.CreationTime
			metaInfo.UpdatedTime = metaBinary.UpdatedTime
			metaInfo.Period = metaBinary.Period
			metaInfo.Status = 0      // TODO - Figure this out
			metaInfo.ReferredCnt = 0 // TODO - Figure this out
			metaInfo.ReferDataID = 0 // TODO - Figure this out
			metaInfo.Flag = metaBinary.Flag
			metaInfo.ReferredTime = metaBinary.ReferredTime
			metaInfo.ExpireTime = metaBinary.ExpireTime
			metaInfo.Tags = metaBinary.Tags
			metaInfo.Ratings = make([]*datastore_types.DataStoreRatingInfoWithSlot, 0)

			pMetaInfo = append(pMetaInfo, metaInfo)

			result := nex.NewResultSuccess(nex.Errors.DataStore.Unknown)
			pResults = append(pResults, result)
		}

		rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

		rmcResponseStream.WriteListStructure(pMetaInfo)
		rmcResponseStream.WriteListResult(pResults)

		rmcResponseBody := rmcResponseStream.Bytes()

		rmcResponse.SetSuccess(datastore.MethodGetMetas, rmcResponseBody)
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
