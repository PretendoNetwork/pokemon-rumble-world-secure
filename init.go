package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"strings"

	pb_account "github.com/PretendoNetwork/grpc-go/account"
	pb_friends "github.com/PretendoNetwork/grpc-go/friends"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/PretendoNetwork/pokemon-rumble-world/database"
	"github.com/PretendoNetwork/pokemon-rumble-world/globals"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func init() {
	globals.Logger = plogger.NewLogger()

	var err error

	err = godotenv.Load()
	if err != nil {
		globals.Logger.Warning("Error loading .env file")
	}

	postgresURI := os.Getenv("PN_PRW_POSTGRES_URI")
	authenticationServerPort := os.Getenv("PN_PRW_AUTHENTICATION_SERVER_PORT")
	secureServerHost := os.Getenv("PN_PRW_SECURE_SERVER_HOST")
	secureServerPort := os.Getenv("PN_PRW_SECURE_SERVER_PORT")
	hppServerPort := os.Getenv("PN_PRW_HPP_SERVER_PORT")
	accountGRPCHost := os.Getenv("PN_PRW_ACCOUNT_GRPC_HOST")
	accountGRPCPort := os.Getenv("PN_PRW_ACCOUNT_GRPC_PORT")
	accountGRPCAPIKey := os.Getenv("PN_PRW_ACCOUNT_GRPC_API_KEY")
	friendsGRPCHost := os.Getenv("PN_PRW_FRIENDS_GRPC_HOST")
	friendsGRPCPort := os.Getenv("PN_PRW_FRIENDS_GRPC_PORT")
	friendsGRPCAPIKey := os.Getenv("PN_PRW_FRIENDS_GRPC_API_KEY")
	s3Endpoint := os.Getenv("PN_PRW_S3_ENDPOINT")
	s3Region := os.Getenv("PN_PRW_S3_REGION")
	s3AccessKey := os.Getenv("PN_PRW_S3_ACCESS_KEY")
	s3AccessSecret := os.Getenv("PN_PRW_S3_ACCESS_SECRET")
	s3Bucket := os.Getenv("PN_PRW_S3_BUCKET")

	if strings.TrimSpace(postgresURI) == "" {
		globals.Logger.Error("PN_PRW_POSTGRES_URI environment variable not set")
		os.Exit(0)
	}

	kerberosPassword := make([]byte, 0x10)
	_, err = rand.Read(kerberosPassword)
	if err != nil {
		globals.Logger.Error("Error generating Kerberos password")
		os.Exit(0)
	}

	globals.KerberosPassword = string(kerberosPassword)

	globals.AuthenticationServerAccount = nex.NewAccount(types.NewPID(1), "Quazal Authentication", globals.KerberosPassword)
	globals.SecureServerAccount = nex.NewAccount(types.NewPID(2), "Quazal Rendez-Vous", globals.KerberosPassword)

	if strings.TrimSpace(authenticationServerPort) == "" {
		globals.Logger.Error("PN_PRW_AUTHENTICATION_SERVER_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(authenticationServerPort); err != nil {
		globals.Logger.Errorf("PN_PRW_AUTHENTICATION_SERVER_PORT is not a valid port. Expected 0-65535, got %s", authenticationServerPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_PRW_AUTHENTICATION_SERVER_PORT is not a valid port. Expected 0-65535, got %s", authenticationServerPort)
		os.Exit(0)
	}

	if strings.TrimSpace(secureServerHost) == "" {
		globals.Logger.Error("PN_PRW_SECURE_SERVER_HOST environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(secureServerPort) == "" {
		globals.Logger.Error("PN_PRW_SECURE_SERVER_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(secureServerPort); err != nil {
		globals.Logger.Errorf("PN_PRW_SECURE_SERVER_PORT is not a valid port. Expected 0-65535, got %s", secureServerPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_PRW_SECURE_SERVER_PORT is not a valid port. Expected 0-65535, got %s", secureServerPort)
		os.Exit(0)
	}

	if strings.TrimSpace(hppServerPort) == "" {
		globals.Logger.Error("PN_PRW_HPP_SERVER_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(hppServerPort); err != nil {
		globals.Logger.Errorf("PN_PRW_HPP_SERVER_PORT is not a valid port. Expected 0-65535, got %s", hppServerPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_PRW_HPP_SERVER_PORT is not a valid port. Expected 0-65535, got %s", hppServerPort)
		os.Exit(0)
	}

	if strings.TrimSpace(accountGRPCHost) == "" {
		globals.Logger.Error("PN_PRW_ACCOUNT_GRPC_HOST environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(accountGRPCPort) == "" {
		globals.Logger.Error("PN_PRW_ACCOUNT_GRPC_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(accountGRPCPort); err != nil {
		globals.Logger.Errorf("PN_PRW_ACCOUNT_GRPC_PORT is not a valid port. Expected 0-65535, got %s", accountGRPCPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_PRW_ACCOUNT_GRPC_PORT is not a valid port. Expected 0-65535, got %s", accountGRPCPort)
		os.Exit(0)
	}

	if strings.TrimSpace(accountGRPCAPIKey) == "" {
		globals.Logger.Warning("Insecure gRPC server detected. PN_PRW_ACCOUNT_GRPC_API_KEY environment variable not set")
	}

	globals.GRPCAccountClientConnection, err = grpc.Dial(fmt.Sprintf("%s:%s", accountGRPCHost, accountGRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		globals.Logger.Criticalf("Failed to connect to account gRPC server: %v", err)
		os.Exit(0)
	}

	globals.GRPCAccountClient = pb_account.NewAccountClient(globals.GRPCAccountClientConnection)
	globals.GRPCAccountCommonMetadata = metadata.Pairs(
		"X-API-Key", accountGRPCAPIKey,
	)

	if strings.TrimSpace(friendsGRPCHost) == "" {
		globals.Logger.Error("PN_PRW_FRIENDS_GRPC_HOST environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(friendsGRPCPort) == "" {
		globals.Logger.Error("PN_PRW_FRIENDS_GRPC_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(friendsGRPCPort); err != nil {
		globals.Logger.Errorf("PN_PRW_FRIENDS_GRPC_PORT is not a valid port. Expected 0-65535, got %s", friendsGRPCPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_PRW_FRIENDS_GRPC_PORT is not a valid port. Expected 0-65535, got %s", friendsGRPCPort)
		os.Exit(0)
	}

	if strings.TrimSpace(friendsGRPCAPIKey) == "" {
		globals.Logger.Warning("Insecure gRPC server detected. PN_PRW_FRIENDS_GRPC_API_KEY environment variable not set")
	}

	globals.GRPCFriendsClientConnection, err = grpc.Dial(fmt.Sprintf("%s:%s", friendsGRPCHost, friendsGRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		globals.Logger.Criticalf("Failed to connect to friends gRPC server: %v", err)
		os.Exit(0)
	}

	globals.GRPCFriendsClient = pb_friends.NewFriendsClient(globals.GRPCFriendsClientConnection)
	globals.GRPCFriendsCommonMetadata = metadata.Pairs(
		"X-API-Key", friendsGRPCAPIKey,
	)

	if strings.TrimSpace(s3Endpoint) == "" {
		globals.Logger.Error("PN_PRW_S3_ENDPOINT environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(s3Region) == "" {
		globals.Logger.Error("PN_PRW_S3_REGION environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(s3AccessKey) == "" {
		globals.Logger.Error("PN_PRW_S3_ACCESS_KEY environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(s3AccessSecret) == "" {
		globals.Logger.Error("PN_PRW_S3_ACCESS_SECRET environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(s3Bucket) == "" {
		globals.Logger.Error("PN_PRW_S3_BUCKET environment variable not set")
		os.Exit(0)
	}

	staticCredentials := credentials.NewStaticCredentialsProvider(s3AccessKey, s3AccessSecret, "")

	endpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: s3Endpoint,
			SigningRegion: s3Region,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(s3Region),
		config.WithCredentialsProvider(staticCredentials),
		config.WithEndpointResolverWithOptions(endpointResolver),
	)

	if err != nil {
		globals.Logger.Criticalf("Failed to create S3 config: %v", err)
		os.Exit(0)
	}

	globals.S3Client = s3.NewFromConfig(cfg)
	globals.S3PresignClient = s3.NewPresignClient(globals.S3Client)
	globals.S3PresignPostClient = globals.NewPresignClient(cfg)

	database.ConnectPostgres()
}
