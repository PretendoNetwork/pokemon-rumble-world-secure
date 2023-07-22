package nex_datastore

import (
	"fmt"
	"os"
	"time"

	"github.com/PretendoNetwork/pokemon-rumble-world/database"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
)

func PreparePostObjectV1(err error, client *nex.Client, callID uint32, param *datastore_types.DataStorePreparePostParamV1) {
	rmcResponse := nex.NewRMCResponse(0, callID)

	var dataID uint32
	var insertErr error
	if err != nil {
		globals.Logger.Error(err.Error())
		rmcResponse.SetError(nex.Errors.DataStore.Unknown)
	} else {
		dataID, insertErr = database.InsertNotificationMetaByDataStorePreparePostParamV1WithOwnerPID(param, client.PID())
		if insertErr != nil {
			globals.Logger.Error(insertErr.Error())
			rmcResponse.SetError(nex.Errors.DataStore.Unknown)
		}
	}

	var key string
	var res *globals.PresignedPostObject
	var presignGetErr error
	if err == nil && insertErr == nil {
		bucket := os.Getenv("PN_PRW_S3_BUCKET")
		key = fmt.Sprintf("data/%011d", dataID)

		input := &globals.PostObjectInput{
			Bucket:    bucket,
			Key:       key,
			ExpiresIn: time.Minute * 15,
		}

		res, presignGetErr = globals.S3PresignPostClient.PresignPostObject(input)
		if presignGetErr != nil {
			globals.Logger.Error(presignGetErr.Error())
			rmcResponse.SetError(nex.Errors.DataStore.Unknown)
		}
	}

	if err == nil && insertErr == nil && presignGetErr == nil {
		fieldKey := datastore_types.NewDataStoreKeyValue()
		fieldKey.Key = "key"
		fieldKey.Value = key

		fieldCredential := datastore_types.NewDataStoreKeyValue()
		fieldCredential.Key = "X-Amz-Credential"
		fieldCredential.Value = res.Credential

		fieldSecurityToken := datastore_types.NewDataStoreKeyValue()
		fieldSecurityToken.Key = "X-Amz-Security-Token"
		fieldSecurityToken.Value = ""

		fieldAlgorithm := datastore_types.NewDataStoreKeyValue()
		fieldAlgorithm.Key = "X-Amz-Algorithm"
		fieldAlgorithm.Value = "AWS4-HMAC-SHA256"

		fieldDate := datastore_types.NewDataStoreKeyValue()
		fieldDate.Key = "X-Amz-Date"
		fieldDate.Value = res.Date

		fieldPolicy := datastore_types.NewDataStoreKeyValue()
		fieldPolicy.Key = "policy"
		fieldPolicy.Value = res.Policy

		fieldSignature := datastore_types.NewDataStoreKeyValue()
		fieldSignature.Key = "X-Amz-Signature"
		fieldSignature.Value = res.Signature

		pReqPostInfo := datastore_types.NewDataStoreReqPostInfoV1()

		pReqPostInfo.DataID = dataID
		pReqPostInfo.URL = res.URL
		pReqPostInfo.RequestHeaders = []*datastore_types.DataStoreKeyValue{}
		pReqPostInfo.FormFields = []*datastore_types.DataStoreKeyValue{
			fieldKey,
			fieldCredential,
			fieldSecurityToken,
			fieldAlgorithm,
			fieldDate,
			fieldPolicy,
			fieldSignature,
		}
		pReqPostInfo.RootCACert = []byte{}

		rmcResponseStream := nex.NewStreamOut(globals.HPPServer)

		rmcResponseStream.WriteStructure(pReqPostInfo)

		rmcResponseBody := rmcResponseStream.Bytes()

		rmcResponse.SetSuccess(datastore.MethodPreparePostObjectV1, rmcResponseBody)
	}

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewHPPPacket(client, nil)

	responsePacket.SetPayload(rmcResponseBytes)

	globals.HPPServer.Send(responsePacket)
}
