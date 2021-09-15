package node

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"code.vegaprotocol.io/protos/vega/api"
	commandspb "code.vegaprotocol.io/protos/vega/commands/v1"

	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Forwarder struct {
	log      *zap.Logger
	nodeCfgs NodesConfig
	clts     []api.TradingServiceClient
	conns    []*grpc.ClientConn
	next     uint64
}

func NewForwarder(log *zap.Logger, nodeConfigs NodesConfig) (*Forwarder, error) {
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

	return &Forwarder{
		log:      log,
		nodeCfgs: nodeConfigs,
		clts:     clts,
		conns:    conns,
	}, nil
}

func (n *Forwarder) Stop() error {
	for i, v := range n.nodeCfgs.Hosts {
		n.log.Info("closing grpc client", zap.String("address", v))
		if err := n.conns[i].Close(); err != nil {
			return err
		}
	}
	return nil
}

func (n *Forwarder) HealthCheck(ctx context.Context) error {
	req := api.GetVegaTimeRequest{}
	return backoff.Retry(
		func() error {
			clt := n.nextClt()
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

func (n *Forwarder) LastBlockHeight(ctx context.Context) (uint64, error) {
	req := api.LastBlockHeightRequest{}
	var height uint64
	err := backoff.Retry(
		func() error {
			clt := n.nextClt()
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

	if err != nil {
		n.log.Error("could not get last block", zap.Error(err))
	} else {
		n.log.Debug("last block when sending transaction",
			zap.Time("request.time", time.Now()),
			zap.Uint64("block.height", height),
		)
	}

	return height, err
}

func (n *Forwarder) SendTx(ctx context.Context, tx *commandspb.Transaction, ty api.SubmitTransactionV2Request_Type) error {
	req := api.SubmitTransactionV2Request{
		Tx:   tx,
		Type: ty,
	}
	return backoff.Retry(
		func() error {
			clt := n.nextClt()
			resp, err := clt.SubmitTransactionV2(ctx, &req)
			if err != nil {
				n.log.Error("failed to send transaction", zap.Error(err))
				return err
			}
			n.log.Debug("response from SubmitTransactionV2", zap.Bool("success", resp.Success))
			return nil
		},
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), n.nodeCfgs.Retries),
	)
}

func (n *Forwarder) nextClt() api.TradingServiceClient {
	i := atomic.AddUint64(&n.next, 1)
	n.log.Info("sending transaction to Vega node",
		zap.String("host", n.nodeCfgs.Hosts[(int(i)-1)%len(n.clts)]))
	return n.clts[(int(i)-1)%len(n.clts)]
}
