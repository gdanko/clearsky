package api

import (
	"errors"
	"os"

	"github.com/gdanko/clearsky/globals"
	"github.com/gdanko/clearsky/util"
	"github.com/sirupsen/logrus"
)

func readCredentials() (blueSkyHandle, blueSkyPassword string, err error) {
	blueSkyHandle = os.Getenv("BLUESKY_HANDLE")
	blueSkyPassword = os.Getenv("BLUESKY_PASSWORD")
	if blueSkyHandle == "" {
		return blueSkyHandle, blueSkyPassword, errors.New("Please export BLUESKY_HANDLE")
	}
	if blueSkyPassword == "" {
		return blueSkyHandle, blueSkyPassword, errors.New("Please export BLUESKY_PASSWORD")
	}
	return blueSkyHandle, blueSkyPassword, nil
}

// Read the credentials and store needed bits in a struct
func Authenticate() (blueSkyCredentials globals.BlueSkyCredentials, err error) {
	var (
		blueSkyHandle     string
		blueSkyPassword   string
		logger            *logrus.Logger
		plcDirectoryEntry globals.PlcDirectoryEntry
		serviceEndpoint   string
		sessionDocument   globals.SessionDocument
		userId            string
	)

	logger = util.ConfigureLogger(logrus.DebugLevel, false)

	// Read the credentials first
	blueSkyHandle, blueSkyPassword, err = readCredentials()
	if err != nil {
		return blueSkyCredentials, err
	}

	// Get the target's DID
	userId, err = GetUserDid(blueSkyHandle, logger)
	if err != nil {
		return blueSkyCredentials, err
	}

	// Get the PLC directory info
	plcDirectoryEntry, err = GetPlcDirectoryInfo(userId, logger)
	if err != nil {
		return blueSkyCredentials, err
	}
	serviceEndpoint = plcDirectoryEntry.Service[0].ServiceEndpoint

	// Get the BlueSky session document
	sessionDocument, err = CreateBlueSkySession(blueSkyHandle, blueSkyPassword, serviceEndpoint, logger)
	if err != nil {
		return blueSkyCredentials, err
	}
	blueSkyCredentials = globals.BlueSkyCredentials{UserDid: userId, ServiceEndpoint: serviceEndpoint, AccessJwtCookie: sessionDocument.AccessJwt}

	return blueSkyCredentials, nil
}
