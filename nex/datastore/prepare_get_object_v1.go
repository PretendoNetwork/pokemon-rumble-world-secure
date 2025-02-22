package nex_datastore

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/PretendoNetwork/pokemon-rumble-world/globals"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func PrepareGetObjectV1(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStorePrepareGetParamV1) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.InvalidArgument, err.Error())
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	bucket := os.Getenv("PN_PRW_S3_BUCKET")
	key := fmt.Sprintf("data/%011d", param.DataID)

	var size uint64
	size, err = globals.S3ObjectSize(bucket, key)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	pReqGetInfo := datastore_types.NewDataStoreReqGetInfoV1()
	pReqGetInfo.Size = types.NewUInt32(uint32(size))

	var request *v4.PresignedHTTPRequest
	request, err = globals.S3PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(15 * int64(time.Minute))
	})

	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	pReqGetInfo.URL = types.NewString(request.URL)

	rmcResponseStream := nex.NewByteStreamOut(globals.HPPServer.LibraryVersions(), globals.HPPServer.ByteStreamSettings())

	pReqGetInfo.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodPrepareGetObjectV1
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
