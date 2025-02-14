package nex_datastore

import (
	"fmt"
	"os"

	"github.com/PretendoNetwork/pokemon-rumble-world/database"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"

	"github.com/PretendoNetwork/nex-go/v2"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
)

func CompletePostObjectV1(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStoreCompletePostParamV1) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.InvalidArgument, err.Error())
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	bucket := os.Getenv("PN_PRW_S3_BUCKET")
	key := fmt.Sprintf("data/%011d", param.DataID)

	_, err = globals.S3HeadRequest(bucket, key)
	if err != nil {
		globals.Logger.Error(err.Error())
		// * Report the error to the client if it isn't aware of it
		if param.IsSuccess {
			return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
		}
	}

	// * Only send notifications if we could confirm the object exists
	if err == nil {
		friendList := globals.GetUserFriendPIDs(uint32(connection.PID()))
		for _, pid := range friendList {
			notificationID, notificationErr := database.InsertNotificationByDataIDAndPID(uint32(param.DataID), pid)
			if notificationErr != nil {
				globals.Logger.Critical(notificationErr.Error())
				continue
			}

			// TODO - What is the last number? Looks constant per-game
			notification := fmt.Sprintf("%d,%d,%d", notificationID, pid, 1425562994)
			key := fmt.Sprintf("notify/%011d", pid)

			_, err = globals.S3PutRequest(bucket, key, notification)
			if err != nil {
				globals.Logger.Critical(err.Error())
			}
		}
	}

	rmcResponse := nex.NewRMCSuccess(endpoint, nil)
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodCompletePostObjectV1
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
