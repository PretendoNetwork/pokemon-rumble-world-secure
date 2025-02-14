package nex_datastore

import (
	"github.com/PretendoNetwork/nex-go/v2"
	nex_types "github.com/PretendoNetwork/nex-go/v2/types"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/PretendoNetwork/pokemon-rumble-world/database"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"
	"github.com/PretendoNetwork/pokemon-rumble-world/types"
)

func SearchObject(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStoreSearchParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.InvalidArgument, err.Error())
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	metaBinaries := make([]*types.MetaBinary, 0)
	var totalCount uint32

	if param.SearchTarget == 10 { // * Search for meta binary of this client
		metaBinary := database.GetMetaInfoByOwnerPID(uint32(connection.PID()))
		if metaBinary.DataID != 0 {
			metaBinaries = append(metaBinaries, metaBinary)
		}
		totalCount = uint32(len(metaBinaries))
	} else if len(param.DataTypes) > 0 { // * The data type is given inside the DataTypes param
		metaBinaries = database.GetMetaInfosByDataStoreSearchParam(param)
		dataTypes := make([]uint16, len(param.DataTypes))
		for i, dataType := range param.DataTypes {
			dataTypes[i] = uint16(dataType)
		}
		totalCount = database.GetTotalMetaInfosByDataTypes(dataTypes)
	} else if len(param.OwnerIDs) == 0 { // * Ignore unknown request for PID = 2 (Rendez-Vous)
		metaBinaries = database.GetAllMetaInfosByDataStoreSearchParam(param)
		totalCount = database.GetTotalMetaInfos()
	}

	pSearchResult := datastore_types.NewDataStoreSearchResult()

	pSearchResult.TotalCount = nex_types.NewUInt32(totalCount)

	if totalCount > uint32(len(metaBinaries)) {
		pSearchResult.TotalCountType = 1 // * Not all results are returned
	}

	pSearchResult.Result = make([]datastore_types.DataStoreMetaInfo, 0, len(metaBinaries))

	for i := 0; i < len(metaBinaries); i++ {
		metaBinary := metaBinaries[i]
		result := datastore_types.NewDataStoreMetaInfo()

		result.DataID = nex_types.NewUInt64(uint64(metaBinary.DataID))
		result.OwnerID = nex_types.NewPID(uint64(metaBinary.OwnerPID))
		result.Size = 0
		result.Name = nex_types.NewString(metaBinary.Name)
		result.DataType = nex_types.NewUInt16(metaBinary.DataType)
		result.Permission = datastore_types.NewDataStorePermission()
		result.Permission.Permission = nex_types.NewUInt8(metaBinary.Permission)
		result.Permission.RecipientIDs = make([]nex_types.PID, 0)
		result.DelPermission = datastore_types.NewDataStorePermission()
		result.DelPermission.Permission = nex_types.NewUInt8(metaBinary.DeletePermission)
		result.DelPermission.RecipientIDs = make([]nex_types.PID, 0)
		result.CreatedTime = metaBinary.CreationTime
		result.UpdatedTime = metaBinary.UpdatedTime
		result.Period = nex_types.NewUInt16(metaBinary.Period)
		result.Status = 0      // TODO - Figure this out
		result.ReferredCnt = 0 // TODO - Figure this out
		result.ReferDataID = 0 // TODO - Figure this out
		result.Flag = nex_types.NewUInt32(metaBinary.Flag)
		result.ReferredTime = metaBinary.ReferredTime
		result.ExpireTime = metaBinary.ExpireTime

		tags := make([]nex_types.String, len(metaBinaries[i].Tags))
		for j, tag := range metaBinaries[i].Tags {
			tags[j] = nex_types.NewString(tag)
		}
		result.Tags = tags

		result.Ratings = make([]datastore_types.DataStoreRatingInfoWithSlot, 0)

		pSearchResult.Result = append(pSearchResult.Result, result)
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	pSearchResult.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodSearchObject
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
