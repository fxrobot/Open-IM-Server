package cronTask

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

const cronTaskOperationID = "cronTaskOperationID-"

func StartCronTask() {
	log.NewInfo(utils.OperationIDGenerator(), "start cron task")
	c := cron.New()
	_, err := c.AddFunc(config.Config.Mongo.ChatRecordsClearTime, func() {
		operationID := getCronTaskOperationID()
		userIDList, err := im_mysql_model.SelectAllUserID()
		if err == nil {
			log.NewDebug(operationID, utils.GetSelfFuncName(), "userIDList: ", userIDList)
			for _, userID := range userIDList {
				if err := DeleteMongoMsgAndResetRedisSeq(operationID, userID); err != nil {
					log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), userID)
				}
			}
		} else {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
		}

		workingGroupIDList, err := im_mysql_model.GetGroupIDListByGroupType(constant.WorkingGroup)
		if err == nil {
			for _, groupID := range workingGroupIDList {
				userIDList, err = rocksCache.GetGroupMemberIDListFromCache(groupID)
				if err != nil {
					log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID)
					continue
				}
				log.NewDebug(operationID, utils.GetSelfFuncName(), "groupID:", groupID, "userIDList:", userIDList)
				if err := ResetUserGroupMinSeq(operationID, groupID, userIDList); err != nil {
					log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID, userIDList)
				}

			}
		} else {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
			return
		}
	})
	if err != nil {
		fmt.Println("start cron failed", err.Error())
		panic(err)
	}
	c.Start()
	fmt.Println("start cron task success")
	for {
		time.Sleep(time.Second)
	}
}

func getCronTaskOperationID() string {
	return cronTaskOperationID + utils.OperationIDGenerator()
}
