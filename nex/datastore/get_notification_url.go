package nex_datastore

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PretendoNetwork/pokemon-rumble-world/globals"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetNotificationURL(err error, client *nex.Client, callID uint32, param *datastore_types.DataStoreGetNotificationURLParam) {
	rmcResponse := nex.NewRMCResponse(0, callID)

	var url, key string
	var presignGetErr error
	if err != nil {
		globals.Logger.Error(err.Error())
		rmcResponse.SetError(nex.Errors.DataStore.Unknown)
	} else {
		bucket := os.Getenv("PN_PRW_S3_BUCKET")
		key = fmt.Sprintf("notify/%011d", client.PID())

		var request *v4.PresignedHTTPRequest
		request, presignGetErr = globals.S3PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(7 * 24 * int64(time.Hour)) // * 1 week
		})

		if presignGetErr != nil {
			globals.Logger.Error(presignGetErr.Error())
			rmcResponse.SetError(nex.Errors.DataStore.Unknown)
		} else {
			url = request.URL
		}
	}

	if err == nil && presignGetErr == nil {
		info := datastore_types.NewDataStoreReqGetNotificationURLInfo()
		info.URL, info.Query, _ = strings.Cut(url, key) // * Split URL and query
		info.Key = key

		rmcResponseStream := nex.NewStreamOut(globals.HPPServer)

		rmcResponseStream.WriteStructure(info)

		rmcResponseBody := rmcResponseStream.Bytes()

		rmcResponse.SetSuccess(datastore.MethodGetNotificationURL, rmcResponseBody)
	}

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewHPPPacket(client, nil)

	responsePacket.SetPayload(rmcResponseBytes)

	globals.HPPServer.Send(responsePacket)
}
