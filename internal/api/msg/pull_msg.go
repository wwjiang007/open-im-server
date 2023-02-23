package msg

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/tokenverify"
	"OpenIM/pkg/proto/msg"
	sdkws "OpenIM/pkg/proto/sdkws"
	"OpenIM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsUserPullMsg struct {
	ReqIdentifier *int   `json:"reqIdentifier" binding:"required"`
	SendID        string `json:"sendID" binding:"required"`
	OperationID   string `json:"operationID" binding:"required"`
	Data          struct {
		SeqBegin *int64 `json:"seqBegin" binding:"required"`
		SeqEnd   *int64 `json:"seqEnd" binding:"required"`
	}
}

type paramsUserPullMsgBySeqList struct {
	ReqIdentifier int      `json:"reqIdentifier" binding:"required"`
	SendID        string   `json:"sendID" binding:"required"`
	OperationID   string   `json:"operationID" binding:"required"`
	SeqList       []uint32 `json:"seqList"`
}

func PullMsgBySeqList(c *gin.Context) {
	params := paramsUserPullMsgBySeqList{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	token := c.Request.Header.Get("token")
	if ok, err := tokenverify.VerifyToken(token, params.SendID); !ok {
		if err != nil {
			log.NewError(params.OperationID, utils.GetSelfFuncName(), err.Error(), token, params.SendID)
		}
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "token validate err"})
		return
	}
	pbData := sdkws.PullMessageBySeqListReq{}
	pbData.UserID = params.SendID
	pbData.OperationID = params.OperationID
	pbData.SeqList = params.SeqList

	grpcConn := rpc.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, pbData.OperationID)
	if grpcConn == nil {
		errMsg := pbData.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(pbData.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	msgClient := msg.NewMsgClient(grpcConn)
	reply, err := msgClient.PullMessageBySeqList(context.Background(), &pbData)
	if err != nil {
		log.Error(pbData.OperationID, "PullMessageBySeqList error", err.Error())
		return
	}
	log.NewInfo(pbData.OperationID, "rpc call success to PullMessageBySeqList", reply.String(), len(reply.List))
	c.JSON(http.StatusOK, gin.H{
		"errCode":       reply.ErrCode,
		"errMsg":        reply.ErrMsg,
		"reqIdentifier": params.ReqIdentifier,
		"data":          reply.List,
	})
}
