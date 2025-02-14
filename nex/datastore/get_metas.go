package nex_datastore

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/PretendoNetwork/pokemon-rumble-world/database"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"
)

func GetMetas(err error, packet nex.PacketInterface, callID uint32, dataIDs types.List[types.UInt64], param datastore_types.DataStoreGetMetaParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.InvalidArgument, err.Error())
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	rawDataIDs := make([]uint64, len(dataIDs))
	for i, dataID := range dataIDs {
		rawDataIDs[i] = uint64(dataID)
	}

	metaBinaries := database.GetMetaBinariesByDataIDs(rawDataIDs)

	var pMetaInfo types.List[datastore_types.DataStoreMetaInfo] = make([]datastore_types.DataStoreMetaInfo, 0, len(metaBinaries))
	var pResults types.List[types.QResult] = make([]types.QResult, 0, len(metaBinaries))

	for i := 0; i < len(metaBinaries); i++ {
		metaBinary := metaBinaries[i]
		metaInfo := datastore_types.NewDataStoreMetaInfo()

		metaInfo.DataID = types.NewUInt64(uint64(metaBinary.DataID))
		metaInfo.OwnerID = types.NewPID(uint64(metaBinary.OwnerPID))
		metaInfo.Size = 0
		metaInfo.Name = types.NewString(metaBinary.Name)
		metaInfo.DataType = types.NewUInt16(metaBinary.DataType)
		metaInfo.MetaBinary = metaBinary.Buffer
		metaInfo.Permission = datastore_types.NewDataStorePermission()
		metaInfo.Permission.Permission = types.NewUInt8(metaBinary.Permission)
		metaInfo.Permission.RecipientIDs = make([]types.PID, 0)
		metaInfo.DelPermission = datastore_types.NewDataStorePermission()
		metaInfo.DelPermission.Permission = types.NewUInt8(metaBinary.DeletePermission)
		metaInfo.DelPermission.RecipientIDs = make([]types.PID, 0)
		metaInfo.CreatedTime = metaBinary.CreationTime
		metaInfo.UpdatedTime = metaBinary.UpdatedTime
		metaInfo.Period = types.NewUInt16(metaBinary.Period)
		metaInfo.Status = 0      // TODO - Figure this out
		metaInfo.ReferredCnt = 0 // TODO - Figure this out
		metaInfo.ReferDataID = 0 // TODO - Figure this out
		metaInfo.Flag = types.NewUInt32(metaBinary.Flag)
		metaInfo.ReferredTime = metaBinary.ReferredTime
		metaInfo.ExpireTime = metaBinary.ExpireTime

		tags := make([]types.String, len(metaBinaries[i].Tags))
		for j, tag := range metaBinaries[i].Tags {
			tags[j] = types.NewString(tag)
		}
		metaInfo.Tags = tags

		metaInfo.Ratings = make([]datastore_types.DataStoreRatingInfoWithSlot, 0)

		pMetaInfo = append(pMetaInfo, metaInfo)

		result := types.NewQResultSuccess(nex.ResultCodes.DataStore.Unknown)
		pResults = append(pResults, result)
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	pMetaInfo.WriteTo(rmcResponseStream)
	pResults.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodGetMetas
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
