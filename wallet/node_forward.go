package wallet

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/cenkalti/backoff/v4"
	"github.com/vegaprotocol/api-clients/go/generated/code.vegaprotocol.io/vega/proto/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type nodeForward struct {
	log      *zap.Logger
	nodeCfgs NodesConfig
	clts     []api.TradingServiceClient
	conns    []*grpc.ClientConn
	next     uint64
}

func NewNodeForward(log *zap.Logger, nodeConfigs NodesConfig) (*nodeForward, error) {
	if len(nodeConfigs.Hosts) <= 0 {
		return nil, errors.New("no node specified for node forwarding")
	}

	var (
		clts  []api.TradingServiceClient
		conns []*grpc.ClientConn
	)
	for _, v := range nodeConfigs.Hosts {
		conn, err := grpc.Dial(v, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		conns = append(conns, conn)
		clts = append(clts, api.NewTradingServiceClient(conn))
	}

	return &nodeForward{
		log:      log,
		nodeCfgs: nodeConfigs,
		clts:     clts,
		conns:    conns,
	}, nil
}

func (n *nodeForward) Stop() error {
	for i, v := range n.nodeCfgs.Hosts {
		n.log.Info("closing grpc client", zap.String("address", v))
		if err := n.conns[i].Close(); err != nil {
			return err
		}
	}
	return nil
}

func (n *nodeForward) Send(ctx context.Context, tx *SignedBundle, ty api.SubmitTransactionRequest_Type) error {
	req := api.SubmitTransactionRequest{
		Tx:   tx.IntoProto(),
		Type: ty,
	}
	return backoff.Retry(
		func() error {
			clt := n.nextClt()
			resp, err := clt.SubmitTransaction(ctx, &req)
			if err != nil {
				return err
			}
			n.log.Debug("response from SubmitTransaction", zap.Bool("success", resp.Success))
			return nil
		},
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), n.nodeCfgs.Retries),
	)
}

func (r *nodeForward) nextClt() api.TradingServiceClient {
	n := atomic.AddUint64(&r.next, 1)
	r.log.Info("sending transaction to vega node",
		zap.String("host", r.nodeCfgs.Hosts[(int(n)-1)%len(r.clts)]))
	return r.clts[(int(n)-1)%len(r.clts)]
}
