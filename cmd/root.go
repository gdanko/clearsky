package cmd

import (
	"database/sql"

	"github.com/gdanko/clearsky/globals"
	"github.com/gdanko/clearsky/pkg/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	accountName           string
	batchOperationTimeout int
	blueSkyCredentials    globals.BlueSkyCredentials
	db                    *sql.DB
	debugFlag             bool
	defaultLogLevel       = "info"
	err                   error
	listMaxResults        int
	logger                *logrus.Logger
	logLevel              logrus.Level
	logLevelStr           string
	logLevelMap           = map[string]logrus.Level{
		"panic": logrus.PanicLevel,
		"fatal": logrus.FatalLevel,
		"error": logrus.ErrorLevel,
		"warn":  logrus.WarnLevel,
		"info":  logrus.InfoLevel,
		"debug": logrus.DebugLevel,
		"trace": logrus.TraceLevel,
	}
	nocolorFlag        bool
	serviceEndpoint    string
	showBlockedUsers   bool
	showBlockedByUsers bool
	showListNames      bool
	displayName        string
	userId             string
	versionFull        bool
	rootCmd            = &cobra.Command{
		Use:   "clearsky",
		Short: "clearsky is a command line interface for the clearsky.services API. Written by Juicy Sharts (@juicysharts.bsky.social)",
		Long:  "clearsky is a command line interface for the clearsky.services API. Written by Juicy Sharts (@juicysharts.bsky.social)",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	blueSkyCredentials, err = api.Authenticate()
	if err != nil {
		panic(err)
	}
	globals.SetCredentials(blueSkyCredentials)
	GetPersistenFlags(rootCmd)
}
