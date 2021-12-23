/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/3/4 11:18).
 */
package im_mysql_msg_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	pbMsg "Open_IM/pkg/proto/chat"
	"Open_IM/pkg/utils"
	"database/sql"
	"time"
)

// ChatLog Chat information table structure
type ChatLog struct {
	MsgId            string         `gorm:"primary_key"`               // Chat history primary key ID
	SendID           string         `gorm:"column:send_id"`            // Send ID
	RecvID           string         `gorm:"column:recv_id"`            //Receive ID
	SendTime         time.Time      `gorm:"column:send_time"`          // Send time
	SessionType      int32          `gorm:"column:session_type"`       // Session type
	ContentType      int32          `gorm:"column:content_type"`       // Message content type
	MsgFrom          int32          `gorm:"column:msg_from"`           // Source, user, system
	Content          string         `gorm:"column:content"`            // Chat content
	SenderPlatformID int32          `gorm:"column:sender_platform_id"` //The sender's platform ID
	Remark           sql.NullString `gorm:"column:remark"`             // remark
}

func InsertMessageToChatLog(msg pbMsg.MsgDataToMQ) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	chatLog := ChatLog{
		MsgId:            msg.MsgData.ServerMsgID,
		SendID:           msg.MsgData.SendID,
		SendTime:         utils.UnixNanoSecondToTime(msg.MsgData.SendTime),
		SessionType:      msg.MsgData.SessionType,
		ContentType:      msg.MsgData.ContentType,
		MsgFrom:          msg.MsgData.MsgFrom,
		Content:          string(msg.MsgData.Content),
		SenderPlatformID: msg.MsgData.SenderPlatformID,
	}
	switch msg.MsgData.SessionType {
	case constant.GroupChatType:
		chatLog.RecvID = msg.MsgData.GroupID
	case constant.SingleChatType:
		chatLog.RecvID = msg.MsgData.RecvID
	}

	return dbConn.Table("chat_log").Create(chatLog).Error
}
