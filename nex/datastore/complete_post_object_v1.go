package nex_datastore

import (
	"fmt"
	"os"

	"github.com/PretendoNetwork/pokemon-rumble-world/database"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
)

func CompletePostObjectV1(err error, client *nex.Client, callID uint32, param *datastore_types.DataStoreCompletePostParamV1) {
	rmcResponse := nex.NewRMCResponse(0, callID)

	var bucket string
	var headRequestErr error
	if err != nil {
		globals.Logger.Error(err.Error())
		rmcResponse.SetError(nex.Errors.DataStore.Unknown)
	}

	if err == nil {
		bucket = os.Getenv("PN_PRW_S3_BUCKET")
		key := fmt.Sprintf("data/%011d", param.DataID)

		_, headRequestErr = globals.S3HeadRequest(bucket, key)
		if headRequestErr != nil {
			globals.Logger.Error(headRequestErr.Error())
			// * Report the error to the client if it isn't aware of it
			if param.IsSuccess {
				rmcResponse.SetError(nex.Errors.DataStore.Unknown)
			} else {
				rmcResponse.SetSuccess(datastore.MethodCompletePostObjectV1, nil)
			}
		}
	}

	if err == nil && headRequestErr == nil {
		friendList := globals.GetUserFriendPIDs(client.PID())
		for _, pid := range friendList {
			notificationID, notificationErr := database.InsertNotificationByDataIDAndPID(param.DataID, pid)
			if notificationErr != nil {
				globals.Logger.Critical(notificationErr.Error())
				continue
			}

			// TODO - What is the last number? Looks constant per-game
			notification := fmt.Sprintf("%d,%d,%d", notificationID, pid, 1425562994)
			key := fmt.Sprintf("notify/%011d", pid)

			_, putRequestErr := globals.S3PutRequest(bucket, key, notification)
			if putRequestErr != nil {
				globals.Logger.Critical(putRequestErr.Error())
			}
		}

		rmcResponse.SetSuccess(datastore.MethodCompletePostObjectV1, nil)
	}

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewHPPPacket(client, nil)

	responsePacket.SetPayload(rmcResponseBytes)

	globals.HPPServer.Send(responsePacket)
}
