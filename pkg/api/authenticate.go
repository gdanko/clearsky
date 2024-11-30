package api

import (
	"encoding/json"
	"errors"
	"fmt"
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
		body              []byte
		blueSkyHandle     string
		blueSkyPassword   string
		credentials       globals.Credentials
		jsonBytes         []byte
		logger            *logrus.Logger
		plcDirectoryEntry globals.PlcDirectoryEntry
		serviceEndpoint   string
		sessionDocument   globals.SessionDocument
		url               string
		userDid           globals.UserDid
		userId            string
	)

	logger = util.ConfigureLogger(logrus.DebugLevel, false)

	// Read the credentials first
	blueSkyHandle, blueSkyPassword, err = readCredentials()
	if err != nil {
		panic(err)
	}

	// Ge the DID from the handle
	url = fmt.Sprintf("https://api.clearsky.services/api/v1/anon/get-did/%s", blueSkyHandle)
	body, err = FetchUrl("GET", url, logger, nil)
	if err != nil {
		return blueSkyCredentials, err
	}

	userDid = globals.UserDid{}
	err = json.Unmarshal(body, &userDid)
	if err != nil {
		return blueSkyCredentials, err
	}
	userId = userDid.Data.DidIdentifier

	// Get the PLC directory info
	url = fmt.Sprintf("https://plc.directory/%s", userId)
	body, err = FetchUrl("GET", url, logger, nil)
	if err != nil {
		return blueSkyCredentials, err
	}
	plcDirectoryEntry = globals.PlcDirectoryEntry{}
	err = json.Unmarshal(body, &plcDirectoryEntry)
	if err != nil {
		return blueSkyCredentials, err
	}
	serviceEndpoint = plcDirectoryEntry.Service[0].ServiceEndpoint

	// Create the session and get the cookie
	url = fmt.Sprintf("%s/xrpc/com.atproto.server.createSession", serviceEndpoint)
	credentials = globals.Credentials{Identifier: blueSkyHandle, Password: blueSkyPassword}
	jsonBytes, err = json.Marshal(credentials)
	if err != nil {
		return blueSkyCredentials, err
	}
	body, err = FetchUrl("POST", url, logger, jsonBytes)
	if err != nil {
		return blueSkyCredentials, err
	}
	sessionDocument = globals.SessionDocument{}
	err = json.Unmarshal(body, &sessionDocument)
	if err != nil {
		return blueSkyCredentials, err
	}
	blueSkyCredentials = globals.BlueSkyCredentials{UserDid: userId, ServiceEndpoint: serviceEndpoint, AccessJwtCookie: sessionDocument.AccessJwt}

	return blueSkyCredentials, nil
}
