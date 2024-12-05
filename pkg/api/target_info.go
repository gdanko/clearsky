package api

import (
	"github.com/gdanko/clearsky/globals"
	"github.com/gdanko/clearsky/util"
	"github.com/sirupsen/logrus"
)

func GetTargetInfo(userId string, logger *logrus.Logger) (blueSkyCredentials globals.BlueSkyCredentials, err error) {
	var (
		plcDirectoryEntry globals.PlcDirectoryEntry
		serviceEndpoint   string
	)

	logger = util.ConfigureLogger(logrus.DebugLevel, false)

	// Get the PLC directory info
	plcDirectoryEntry, err = GetPlcDirectoryInfo(userId, logger)
	if err != nil {
		return blueSkyCredentials, err
	}
	serviceEndpoint = plcDirectoryEntry.Service[0].ServiceEndpoint

	blueSkyCredentials = globals.BlueSkyCredentials{UserDid: userId, ServiceEndpoint: serviceEndpoint}

	return blueSkyCredentials, nil
}
