package globals

import (
	"context"

	pb_account "github.com/PretendoNetwork/grpc-go/account"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-go/v2/globals"
	"google.golang.org/grpc/metadata"
)

func PasswordFromPID(pid types.PID) (string, uint32) {
	ctx := metadata.NewOutgoingContext(context.Background(), GRPCAccountCommonMetadata)

	response, err := GRPCAccountClient.GetNEXData(ctx, &pb_account.GetNEXDataRequest{Pid: uint32(pid)})
	if err != nil {
		globals.Logger.Error(err.Error())
		return "", nex.ResultCodes.RendezVous.InvalidUsername
	}

	// * We only allow tester accounts for now
	if response.AccessLevel < 1 {
		globals.Logger.Errorf("PID %d is not a tester!", response.Pid)
		return "", nex.ResultCodes.RendezVous.AccountDisabled
	}

	return response.Password, 0
}
