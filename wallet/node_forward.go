package wallet

import (
	"context"
	"fmt"

	"code.vegaprotocol.io/go-wallet/proto/api"

	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type nodeForward struct {
	log     *zap.Logger
	nodeCfg NodeConfig
	clt     api.TradingClient
	conn    *grpc.ClientConn
}

func NewNodeForward(log *zap.Logger, nodeConfig NodeConfig) (*nodeForward, error) {
	nodeAddr := fmt.Sprintf("%v:%v", nodeConfig.Host, nodeConfig.Port)
	conn, err := grpc.Dial(nodeAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := api.NewTradingClient(conn)
	return &nodeForward{
		log:     log,
		nodeCfg: nodeConfig,
		clt:     client,
		conn:    conn,
	}, nil
}

func (n *nodeForward) Stop() error {
	n.log.Info("closing grpc client", zap.String("address", fmt.Sprintf("%v:%v", n.nodeCfg.Host, n.nodeCfg.Port)))
	return n.conn.Close()
}

func (n *nodeForward) Send(ctx context.Context, tx *SignedBundle, ty api.SubmitTransactionRequest_Type) error {
	req := api.SubmitTransactionRequest{
		Tx:   tx.IntoProto(),
		Type: ty,
	}
	return backoff.Retry(
		func() error {
			resp, err := n.clt.SubmitTransaction(ctx, &req)
			if err != nil {
				return err
			}
			n.log.Debug("response from SubmitTransaction", zap.Bool("success", resp.Success))
			return nil
		},
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), n.nodeCfg.Retries),
	)
}
