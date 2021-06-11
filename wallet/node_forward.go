package wallet

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/cenkalti/backoff/v4"
	"github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type nodeForward struct {
	log      *zap.Logger
	nodeCfgs NodesConfig
	clts     []api.TradingServiceClient
	cltDatas []api.TradingDataServiceClient
	conns    []*grpc.ClientConn
	next     uint64
}

func NewNodeForward(log *zap.Logger, nodeConfigs NodesConfig) (*nodeForward, error) {
	if len(nodeConfigs.Hosts) <= 0 {
		return nil, errors.New("no node specified for node forwarding")
	}

	var (
		clts     []api.TradingServiceClient
		cltDatas []api.TradingDataServiceClient
		conns    []*grpc.ClientConn
	)
	for _, v := range nodeConfigs.Hosts {
		conn, err := grpc.Dial(v, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		conns = append(conns, conn)
		clts = append(clts, api.NewTradingServiceClient(conn))
		cltDatas = append(cltDatas, api.NewTradingDataServiceClient(conn))
	}

	return &nodeForward{
		log:      log,
		nodeCfgs: nodeConfigs,
		clts:     clts,
		cltDatas: cltDatas,
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

func (n *nodeForward) HealthCheck(ctx context.Context) error {
	req := api.GetVegaTimeRequest{}
	return backoff.Retry(
		func() error {
			clt := n.nextCltData()
			resp, err := clt.GetVegaTime(ctx, &req)
			if err != nil {
				return err
			}
			n.log.Debug("response from GetVegaTime", zap.Int64("timestamp", resp.Timestamp))
			return nil
		},
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), n.nodeCfgs.Retries),
	)
}

func (n *nodeForward) LastBlockHeight(ctx context.Context) (uint64, error) {
	req := api.LastBlockHeightRequest{}
	var height uint64
	err := backoff.Retry(
		func() error {
			clt := n.nextCltData()
			resp, err := clt.LastBlockHeight(ctx, &req)
			if err != nil {
				n.log.Debug("could not get last block", zap.Error(err))
				return err
			}
			height = resp.Height
			return nil
		},
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), n.nodeCfgs.Retries),
	)

	return height, err
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

func (n *nodeForward) nextClt() api.TradingServiceClient {
	i := atomic.AddUint64(&n.next, 1)
	n.log.Info("sending transaction to vega node",
		zap.String("host", n.nodeCfgs.Hosts[(int(i)-1)%len(n.clts)]))
	return n.clts[(int(i)-1)%len(n.clts)]
}

func (n *nodeForward) nextCltData() api.TradingDataServiceClient {
	i := atomic.AddUint64(&n.next, 1)
	n.log.Info("sending healthcheck to vega node",
		zap.String("host", n.nodeCfgs.Hosts[(int(i)-1)%len(n.clts)]))
	return n.cltDatas[(int(i)-1)%len(n.clts)]
}
