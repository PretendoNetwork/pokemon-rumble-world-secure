package nex_datastore

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/PretendoNetwork/pokemon-rumble-world/globals"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func PrepareGetObjectV1(err error, client *nex.Client, callID uint32, param *datastore_types.DataStorePrepareGetParamV1) {
	rmcResponse := nex.NewRMCResponse(0, callID)

	if err != nil {
		globals.Logger.Error(err.Error())
		rmcResponse.SetError(nex.Errors.DataStore.Unknown)
	}

	var pReqGetInfo *datastore_types.DataStoreReqGetInfoV1
	var bucket, key string
	var sizeErr error
	if err == nil {
		bucket = os.Getenv("PN_PRW_S3_BUCKET")
		key = fmt.Sprintf("data/%011d", param.DataID)

		var size uint64
		size, sizeErr = globals.S3ObjectSize(bucket, key)
		if sizeErr != nil {
			globals.Logger.Error(sizeErr.Error())
			rmcResponse.SetError(nex.Errors.DataStore.Unknown)
		} else {
			pReqGetInfo = datastore_types.NewDataStoreReqGetInfoV1()
			pReqGetInfo.Size = uint32(size)
		}
	}

	var presignGetErr error
	if err == nil && sizeErr == nil {
		var request *v4.PresignedHTTPRequest
		request, presignGetErr = globals.S3PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(15 * int64(time.Minute))
		})

		if presignGetErr != nil {
			globals.Logger.Error(presignGetErr.Error())
			rmcResponse.SetError(nex.Errors.DataStore.Unknown)
		} else {
			pReqGetInfo.URL = request.URL
		}
	}

	if err == nil && sizeErr == nil && presignGetErr == nil {
		rmcResponseStream := nex.NewStreamOut(globals.HPPServer)

		rmcResponseStream.WriteStructure(pReqGetInfo)

		rmcResponseBody := rmcResponseStream.Bytes()

		rmcResponse.SetSuccess(datastore.MethodPrepareGetObjectV1, rmcResponseBody)
	}

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewHPPPacket(client, nil)

	responsePacket.SetPayload(rmcResponseBytes)

	globals.HPPServer.Send(responsePacket)
}
