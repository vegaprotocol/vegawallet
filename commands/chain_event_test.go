package commands_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/commands"
	"github.com/stretchr/testify/assert"
	typespb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto"
	commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"
)

func TestCheckChainEvent(t *testing.T) {
	t.Run("Submitting a nil chain event fails", testNilChainEventFails)
	t.Run("Submitting a chain event without event fails", testChainEventWithoutEventFails)
	t.Run("Submitting an ERC20 chain event without tx ID fails", testErc20ChainEventWithoutTxIDFails)
	t.Run("Submitting an ERC20 chain event without nonce succeeds", testErc20ChainEventWithoutNonceSucceeds)
	t.Run("Submitting a built-in chain event without tx ID succeeds", testBuiltInChainEventWithoutTxIDSucceeds)
	t.Run("Submitting a built-in chain event without nonce succeeds", testBuiltInChainEventWithoutNonceSucceeds)
}

func testNilChainEventFails(t *testing.T) {
	err := checkChainEvent(nil)

	assert.Contains(t, err.Get("chain_event"), commands.ErrIsRequired)
}

func testChainEventWithoutEventFails(t *testing.T) {
	event := newErc20ChainEvent()
	event.Event = nil

	err := checkChainEvent(event)

	assert.Contains(t, err.Get("chain_event.event"), commands.ErrIsRequired)
}

func testErc20ChainEventWithoutTxIDFails(t *testing.T) {
	event := newErc20ChainEvent()
	event.TxId = ""

	err := checkChainEvent(event)

	assert.Contains(t, err.Get("chain_event.tx_id"), commands.ErrIsRequired)
}

func testErc20ChainEventWithoutNonceSucceeds(t *testing.T) {
	event := newErc20ChainEvent()
	event.Nonce = 0

	err := checkChainEvent(event)

	assert.NotContains(t, err.Get("chain_event.nonce"), commands.ErrIsRequired)
}

func testBuiltInChainEventWithoutTxIDSucceeds(t *testing.T) {
	event := newBuiltInChainEvent()
	event.TxId = ""

	err := checkChainEvent(event)

	assert.NotContains(t, err.Get("chain_event.tx_id"), commands.ErrIsRequired)
}

func testBuiltInChainEventWithoutNonceSucceeds(t *testing.T) {
	event := newBuiltInChainEvent()
	event.Nonce = 0

	err := checkChainEvent(event)

	assert.NotContains(t, err.Get("chain_event.nonce"), commands.ErrIsRequired)
}

func checkChainEvent(cmd *commandspb.ChainEvent) commands.Errors {
	err := commands.CheckChainEvent(cmd)

	e, ok := err.(commands.Errors)
	if !ok {
		return commands.NewErrors()
	}

	return e
}

func newErc20ChainEvent() *commandspb.ChainEvent {
	return &commandspb.ChainEvent{
		TxId:  "my ID",
		Nonce: RandomPositiveU64(),
		Event: &commandspb.ChainEvent_Erc20{
			Erc20: &typespb.ERC20Event{
				Index:  0,
				Block:  0,
				Action: nil,
			},
		},
	}
}

func newBuiltInChainEvent() *commandspb.ChainEvent {
	return &commandspb.ChainEvent{
		TxId:  "my ID",
		Nonce: RandomPositiveU64(),
		Event: &commandspb.ChainEvent_Builtin{
			Builtin: &typespb.BuiltinAssetEvent{

			},
		},
	}
}
