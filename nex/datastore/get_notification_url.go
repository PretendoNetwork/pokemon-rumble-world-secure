package nex_datastore

import (
	"context"
	"fmt"
	"os"
	"strings"
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

func GetNotificationURL(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStoreGetNotificationURLParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.InvalidArgument, err.Error())
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	var url, key string
	bucket := os.Getenv("PN_PRW_S3_BUCKET")
	key = fmt.Sprintf("notify/%011d", connection.PID())

	var request *v4.PresignedHTTPRequest
	request, err = globals.S3PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(7 * 24 * int64(time.Hour)) // * 1 week
	})

	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	url = request.URL

	info := datastore_types.NewDataStoreReqGetNotificationURLInfo()
	infoURL, infoQuery, _ := strings.Cut(url, key) // * Split URL and query
	info.URL = types.NewString(infoURL)
	info.Query = types.NewString(infoQuery)
	info.Key = types.NewString(key)

	rmcResponseStream := nex.NewByteStreamOut(globals.HPPServer.LibraryVersions(), globals.HPPServer.ByteStreamSettings())

	info.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodGetNotificationURL
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
