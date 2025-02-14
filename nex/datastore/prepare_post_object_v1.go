package nex_datastore

import (
	"fmt"
	"os"
	"time"

	"github.com/PretendoNetwork/pokemon-rumble-world/database"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
)

func PreparePostObjectV1(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStorePreparePostParamV1) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.InvalidArgument, err.Error())
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	dataID, err := database.InsertNotificationMetaByDataStorePreparePostParamV1WithOwnerPID(param, uint32(connection.PID()))
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	bucket := os.Getenv("PN_PRW_S3_BUCKET")
	key := fmt.Sprintf("data/%011d", dataID)

	input := &globals.PostObjectInput{
		Bucket:    bucket,
		Key:       key,
		ExpiresIn: time.Minute * 15,
	}

	res, err := globals.S3PresignPostClient.PresignPostObject(input)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	fieldKey := datastore_types.NewDataStoreKeyValue()
	fieldKey.Key = "key"
	fieldKey.Value = types.NewString(key)

	fieldCredential := datastore_types.NewDataStoreKeyValue()
	fieldCredential.Key = "X-Amz-Credential"
	fieldCredential.Value = types.NewString(res.Credential)

	fieldSecurityToken := datastore_types.NewDataStoreKeyValue()
	fieldSecurityToken.Key = "X-Amz-Security-Token"
	fieldSecurityToken.Value = ""

	fieldAlgorithm := datastore_types.NewDataStoreKeyValue()
	fieldAlgorithm.Key = "X-Amz-Algorithm"
	fieldAlgorithm.Value = "AWS4-HMAC-SHA256"

	fieldDate := datastore_types.NewDataStoreKeyValue()
	fieldDate.Key = "X-Amz-Date"
	fieldDate.Value = types.NewString(res.Date)

	fieldPolicy := datastore_types.NewDataStoreKeyValue()
	fieldPolicy.Key = "policy"
	fieldPolicy.Value = types.NewString(res.Policy)

	fieldSignature := datastore_types.NewDataStoreKeyValue()
	fieldSignature.Key = "X-Amz-Signature"
	fieldSignature.Value = types.NewString(res.Signature)

	pReqPostInfo := datastore_types.NewDataStoreReqPostInfoV1()

	pReqPostInfo.DataID = types.NewUInt32(dataID)
	pReqPostInfo.URL = types.NewString(res.URL)
	pReqPostInfo.RequestHeaders = []datastore_types.DataStoreKeyValue{}
	pReqPostInfo.FormFields = []datastore_types.DataStoreKeyValue{
		fieldKey,
		fieldCredential,
		fieldSecurityToken,
		fieldAlgorithm,
		fieldDate,
		fieldPolicy,
		fieldSignature,
	}
	pReqPostInfo.RootCACert = []byte{}

	rmcResponseStream := nex.NewByteStreamOut(globals.HPPServer.LibraryVersions(), globals.HPPServer.ByteStreamSettings())

	pReqPostInfo.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodPreparePostObjectV1
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
