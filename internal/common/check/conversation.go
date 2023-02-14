package check

import (
	"Open_IM/pkg/common/config"
	discoveryRegistry "Open_IM/pkg/discoveryregistry"
	"Open_IM/pkg/proto/conversation"
	pbConversation "Open_IM/pkg/proto/conversation"
	"context"
	"google.golang.org/grpc"
)

type ConversationChecker struct {
	zk discoveryRegistry.SvcDiscoveryRegistry
}

func NewConversationChecker(zk discoveryRegistry.SvcDiscoveryRegistry) *ConversationChecker {
	return &ConversationChecker{zk: zk}
}

func (c *ConversationChecker) ModifyConversationField(ctx context.Context, req *pbConversation.ModifyConversationFieldReq) (resp *pbConversation.ModifyConversationFieldResp, err error) {
	cc, err := c.getConn()
	if err != nil {
		return nil, err
	}
	resp, err = conversation.NewConversationClient(cc).ModifyConversationField(ctx, req)
	return
}

func (c *ConversationChecker) getConn() (*grpc.ClientConn, error) {
	return c.zk.GetConn(config.Config.RpcRegisterName.OpenImConversationName)
}
