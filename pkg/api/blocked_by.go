package api

import (
	"encoding/json"
	"fmt"

	"github.com/gdanko/clearsky/globals"
	"github.com/sirupsen/logrus"
)

func GetBlockedBy(targetCredentials globals.BlueSkyCredentials, showBlockedByUsers bool, batchOperationTimeout int, listMaxResults int, logger *logrus.Logger) (blockingList map[string]globals.BlockingUser, err error) {
	var (
		blockingListAll     = map[string]globals.BlockingUser{}
		blockedByListPage   globals.BlockedByListPage
		body                []byte
		limitedBlockingList = map[string]globals.BlockingUser{}
		maxPages            = 1000
		totalRecords        int
		url                 string
	)
	blockingList = map[string]globals.BlockingUser{}
	blockingListAll = make(map[string]globals.BlockingUser)

	url = fmt.Sprintf("https://api.clearsky.services/api/v1/anon/single-blocklist/%s", targetCredentials.UserDid)
	body, err = FetchUrl("GET", url, map[string]string{}, nil, logger)
	if err != nil {
		return blockingList, err
	}
	blockedByListPage = globals.BlockedByListPage{}
	err = json.Unmarshal(body, &blockedByListPage)
	if err != nil {
		return blockingList, nil
	}
	for _, blockingUser := range blockedByListPage.Data.Blocklist {
		blockingListAll[blockingUser.DID] = blockingUser
	}

	if len(blockedByListPage.Data.Blocklist) >= 100 {
		for i := 2; i <= maxPages; i++ {
			url = fmt.Sprintf("https://api.clearsky.services/api/v1/anon/single-blocklist/%s/%d", targetCredentials.UserDid, i)
			body, err = FetchUrl("GET", url, map[string]string{}, nil, logger)
			if err != nil {
				return blockingList, err
			}
			blockedByListPage = globals.BlockedByListPage{}
			err = json.Unmarshal(body, &blockedByListPage)
			if err != nil {
				return blockingList, nil
			}
			if len(blockedByListPage.Data.Blocklist) > 0 {
				for _, blockingUser := range blockedByListPage.Data.Blocklist {
					blockingListAll[blockingUser.DID] = blockingUser
				}
			} else {
				break
			}
		}
	}

	if showBlockedByUsers {
		totalRecords = len(blockingListAll)
		if listMaxResults < totalRecords {
			logger.Debugf("Limiting the number of records to %d because the --limit flag was used", listMaxResults)
			limitedBlockingList = make(map[string]globals.BlockingUser)
			for key, value := range blockingListAll {
				limitedBlockingList[key] = value
				if len(limitedBlockingList) == listMaxResults {
					blockingList = limitedBlockingList
					break
				}
			}
		} else {
			blockingList = blockingListAll
		}

		err = processUsersList(&blockingList, batchOperationTimeout, logger)
		if err != nil {
			return blockingList, err
		}
		return blockingList, nil
	}
	return blockingListAll, nil
}
