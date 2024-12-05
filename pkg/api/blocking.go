package api

import (
	"encoding/json"
	"fmt"

	"github.com/gdanko/clearsky/globals"
	"github.com/sirupsen/logrus"
)

func GetBlocking(userId string, showBlockedUsers bool, batchOperationTimeout int, listMaxResults int, logger *logrus.Logger) (blockedList map[string]globals.BlockingUser, err error) {
	var (
		blockedListAll     = map[string]globals.BlockingUser{}
		blockingListPage   globals.BlockingListPage
		blockPageCursor    string
		body               []byte
		headers            map[string]string
		limitedBlockedList = map[string]globals.BlockingUser{}
		listRecordsLimit   = 100
		plcDirectoryEntry  globals.PlcDirectoryEntry
		serviceEndpoint    string
		totalRecords       int
		url                string
	)
	blockedListAll = make(map[string]globals.BlockingUser)

	url = fmt.Sprintf("https://plc.directory/%s", userId)
	body, err = FetchUrl("GET", url, map[string]string{}, nil, logger)
	if err != nil {
		return blockedList, err
	}
	plcDirectoryEntry = globals.PlcDirectoryEntry{}
	err = json.Unmarshal(body, &plcDirectoryEntry)
	if err != nil {
		return blockedList, nil
	}
	serviceEndpoint = plcDirectoryEntry.Service[0].ServiceEndpoint
	for {
		url = fmt.Sprintf("%s/xrpc/com.atproto.repo.listRecords?repo=%s&limit=%d&collection=app.bsky.graph.block&cursor=%s", serviceEndpoint, userId, listRecordsLimit, blockPageCursor)
		headers = map[string]string{
			"Content-Type": "appliation/json",
		}
		body, err = FetchUrl("GET", url, headers, nil, logger)
		if err != nil {
			return blockedList, err
		}
		blockingListPage = globals.BlockingListPage{}
		err = json.Unmarshal(body, &blockingListPage)
		if err != nil {
			return blockedList, nil
		}
		if len(blockingListPage.Records) > 0 {
			for _, blockingUser := range blockingListPage.Records {
				blockedListAll[blockingUser.Value.Subject] = globals.BlockingUser{DID: blockingUser.Value.Subject}
			}
			blockPageCursor = blockingListPage.Cursor
		} else {
			break
		}
	}

	if showBlockedUsers {
		totalRecords = len(blockedListAll)
		if listMaxResults < totalRecords {
			logger.Debugf("Limiting the number of records to %d because the --limit flag was used", listMaxResults)
			limitedBlockedList = make(map[string]globals.BlockingUser)
			for key, value := range blockedListAll {
				limitedBlockedList[key] = value
				if len(limitedBlockedList) == listMaxResults {
					blockedListAll = limitedBlockedList
					break
				}
			}
		} else {
			blockedList = blockedListAll
		}

		err = processUsersList(&blockedList, batchOperationTimeout, logger)
		if err != nil {
			return blockedList, err
		}
		return blockedList, nil
	}
	return blockedListAll, err
}
