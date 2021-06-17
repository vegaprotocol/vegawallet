package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	typespb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto"
	commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"
	oraclespb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/oracles/v1"
)

const (
	MaxDuration30DaysNs int64 = 2592000000000000
)

func CheckProposalSubmission(cmd *commandspb.ProposalSubmission) error {
	return checkProposalSubmission(cmd).ErrorOrNil()
}

func checkProposalSubmission(cmd *commandspb.ProposalSubmission) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("proposal_submission", ErrIsRequired)
	}

	if cmd.Terms == nil {
		return errs.FinalAddForProperty("proposal_submission.terms", ErrIsRequired)
	}

	if cmd.Terms.ClosingTimestamp <= 0 {
		errs.AddForProperty("proposal_submission.terms.closing_timestamp", ErrMustBePositive)
	}
	if cmd.Terms.EnactmentTimestamp <= 0 {
		errs.AddForProperty("proposal_submission.terms.enactment_timestamp", ErrMustBePositive)
	}
	if cmd.Terms.ValidationTimestamp < 0 {
		errs.AddForProperty("proposal_submission.terms.validation_timestamp", ErrMustBePositiveOrZero)
	}

	if cmd.Terms.ClosingTimestamp > cmd.Terms.EnactmentTimestamp {
		errs.AddForProperty("proposal_submission.terms.closing_timestamp",
			errors.New("cannot be after enactment time"),
		)
	}

	if cmd.Terms.ValidationTimestamp >= cmd.Terms.ClosingTimestamp {
		errs.AddForProperty("proposal_submission.terms.validation_timestamp",
			errors.New("cannot be after or equal to closing time"),
		)
	}

	errs.Merge(checkProposalChanges(cmd.Terms))

	return errs
}

func checkProposalChanges(terms *typespb.ProposalTerms) Errors {
	errs := NewErrors()

	if terms.Change == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change", ErrIsRequired)
	}

	switch c := terms.Change.(type) {
	case *typespb.ProposalTerms_NewMarket:
		errs.Merge(checkNewMarketChanges(c))
	case *typespb.ProposalTerms_UpdateNetworkParameter:
		errs.Merge(checkNetworkParameterUpdateChanges(c))
	case *typespb.ProposalTerms_NewAsset:
		errs.Merge(checkNewAssetChanges(c))
	default:
		return errs.FinalAddForProperty("proposal_submission.terms.change", ErrIsNotValid)
	}

	return errs
}

func checkNetworkParameterUpdateChanges(change *typespb.ProposalTerms_UpdateNetworkParameter) Errors {
	errs := NewErrors()

	if change.UpdateNetworkParameter == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.update_network_parameter", ErrIsRequired)
	}

	if change.UpdateNetworkParameter.Changes == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.update_network_parameter.changes", ErrIsRequired)
	}

	parameter := change.UpdateNetworkParameter.Changes

	if len(parameter.Key) == 0 {
		errs.AddForProperty("proposal_submission.terms.change.update_network_parameter.changes.key", ErrIsRequired)
	}

	if len(parameter.Value) == 0 {
		errs.AddForProperty("proposal_submission.terms.change.update_network_parameter.changes.value", ErrIsRequired)
	}

	return errs
}

func checkNewAssetChanges(change *typespb.ProposalTerms_NewAsset) Errors {
	errs := NewErrors()

	if change.NewAsset == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_asset", ErrIsRequired)
	}

	if change.NewAsset.Changes == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_asset.changes", ErrIsRequired)
	}

	if change.NewAsset.Changes.Source == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_asset.changes.source", ErrIsRequired)
	}

	if len(change.NewAsset.Changes.Name) == 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_asset.changes.name", ErrIsRequired)
	}
	if len(change.NewAsset.Changes.Symbol) == 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_asset.changes.symbol", ErrIsRequired)
	}
	if change.NewAsset.Changes.Decimals == 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_asset.changes.decimals", ErrIsRequired)
	}
	if len(change.NewAsset.Changes.TotalSupply) == 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_asset.changes.total_supply", ErrIsRequired)
	}

	totalSupply, err := strconv.ParseUint(change.NewAsset.Changes.TotalSupply, 10, 64)
	if err != nil {
		errs.AddForProperty("proposal_submission.terms.change.new_asset.changes.total_supply", ErrIsNotValidNumber)
	} else if totalSupply == 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_asset.changes.total_supply", ErrMustBePositive)
	}

	switch s := change.NewAsset.Changes.Source.(type) {
	case *typespb.AssetDetails_BuiltinAsset:
		errs.Merge(checkBuiltinAssetSource(s))
	case *typespb.AssetDetails_Erc20:
		errs.Merge(checkERC20AssetSource(s))
	default:
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_asset.changes.source", ErrIsNotValid)
	}

	return errs
}

func checkBuiltinAssetSource(s *typespb.AssetDetails_BuiltinAsset) Errors {
	errs := NewErrors()

	if s.BuiltinAsset == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_asset.changes.source.builtin_asset", ErrIsRequired)
	}

	asset := s.BuiltinAsset

	if len(asset.MaxFaucetAmountMint) == 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_asset.changes.source.builtin_asset.max_faucet_amount_mint", ErrIsRequired)
	}

	maxFaucetAmount, err := strconv.ParseUint(asset.MaxFaucetAmountMint, 10, 64)
	if err != nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_asset.changes.source.builtin_asset.max_faucet_amount_mint", ErrIsNotValidNumber)
	}

	if maxFaucetAmount == 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_asset.changes.source.builtin_asset.max_faucet_amount_mint", ErrMustBePositive)
	}

	return errs
}

func checkERC20AssetSource(s *typespb.AssetDetails_Erc20) Errors {
	errs := NewErrors()

	if s.Erc20 == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_asset.changes.source.erc20", ErrIsRequired)
	}

	asset := s.Erc20

	if len(asset.ContractAddress) == 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_asset.changes.source.erc20.contract_address", ErrIsRequired)
	}

	return errs
}

func checkNewMarketChanges(change *typespb.ProposalTerms_NewMarket) Errors {
	errs := NewErrors()

	if change.NewMarket == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_market", ErrIsRequired)
	}

	errs.Merge(checkLiquidityCommitment(change.NewMarket.LiquidityCommitment))

	if change.NewMarket.Changes == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_market.changes", ErrIsRequired)
	}

	changes := change.NewMarket.Changes

	if changes.DecimalPlaces <= 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.decimal_places", ErrMustBePositive)
	} else if changes.DecimalPlaces >= 150 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.decimal_places", ErrMustBeLessThan150)
	}

	errs.Merge(checkPriceMonitoring(changes.PriceMonitoringParameters))
	errs.Merge(checkLiquidityMonitoring(changes.LiquidityMonitoringParameters))
	errs.Merge(checkInstrument(changes.Instrument))
	errs.Merge(checkTradingMode(changes))
	errs.Merge(checkRiskParameters(changes))

	return errs
}

func checkPriceMonitoring(parameters *typespb.PriceMonitoringParameters) Errors {
	errs := NewErrors()

	if parameters == nil || len(parameters.Triggers) == 0 {
		return errs
	}

	for i, trigger := range parameters.Triggers {
		if trigger.Horizon <= 0 {
			errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers.%d.horizon", i), ErrMustBePositive)
		}
		if trigger.AuctionExtension <= 0 {
			errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers.%d.auction_extension", i), ErrMustBePositive)
		}
		if trigger.Probability <= 0 || trigger.Probability >= 1 {
			errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers.%d.probability", i),
				errors.New("should be between 0 (exclusive) and 1 (exclusive)"),
			)
		}
	}

	return errs
}

func checkLiquidityMonitoring(parameters *typespb.LiquidityMonitoringParameters) Errors {
	errs := NewErrors()

	if parameters == nil {
		return errs
	}

	if parameters.TriggeringRatio < 0 || parameters.TriggeringRatio > 1 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.liquidity_monitoring_parameters.triggering_ratio",
			errors.New("should be between 0 (inclusive) and 1 (inclusive)"),
		)
	}

	if parameters.TargetStakeParameters == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_market.changes.liquidity_monitoring_parameters.target_stake_parameters", ErrIsRequired)
	}

	if parameters.TargetStakeParameters.TimeWindow <= 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.liquidity_monitoring_parameters.target_stake_parameters.time_window", ErrMustBePositive)
	}
	if parameters.TargetStakeParameters.ScalingFactor <= 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.liquidity_monitoring_parameters.target_stake_parameters.scaling_factor", ErrMustBePositive)
	}

	return errs
}

func checkInstrument(instrument *typespb.InstrumentConfiguration) Errors {
	errs := NewErrors()

	if instrument == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_market.changes.instrument", ErrIsRequired)
	}

	if len(instrument.Name) == 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.instrument.name", ErrIsRequired)
	}
	if len(instrument.Code) == 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.instrument.code", ErrIsRequired)
	}

	if instrument.Product == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_market.changes.instrument.product", ErrIsRequired)
	}

	switch product := instrument.Product.(type) {
	case *typespb.InstrumentConfiguration_Future:
		errs.Merge(checkFuture(product.Future))
	default:
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_market.changes.instrument.product", ErrIsNotValid)
	}

	return errs
}

func checkFuture(future *typespb.FutureProduct) Errors {
	errs := NewErrors()

	if future == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_market.changes.instrument.product.future", ErrIsRequired)
	}

	if len(future.SettlementAsset) == 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.instrument.product.future.settlement_asset", ErrIsRequired)
	}
	if len(future.QuoteName) == 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.instrument.product.future.quote_name", ErrIsRequired)
	}

	if len(future.Maturity) == 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.instrument.product.future.maturity", ErrIsRequired)
	}
	_, err := time.Parse(time.RFC3339, future.Maturity)
	if err != nil {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.instrument.product.future.maturity", ErrMustBeValidDate)
	}

	errs.Merge(checkOracleSpec(future))

	return errs
}

func checkOracleSpec(future *typespb.FutureProduct) Errors {
	errs := NewErrors()

	if future.OracleSpec != nil {
		if len(future.OracleSpec.PubKeys) == 0 {
			errs.AddForProperty("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.pub_keys", ErrIsRequired)
		}
		for i, key := range future.OracleSpec.PubKeys {
			if len(strings.TrimSpace(key)) == 0 {
				errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.pub_keys.%d", i), ErrIsNotValid)
			}
		}
		if len(future.OracleSpec.Filters) == 0 {
			errs.AddForProperty("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters", ErrIsRequired)
		} else {
			for i, filter := range future.OracleSpec.Filters {
				if filter.Key == nil {
					errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.%d.key", i), ErrIsNotValid)
				} else {
					if len(filter.Key.Name) == 0 {
						errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.%d.key.name", i), ErrIsRequired)
					}
					if filter.Key.Type == oraclespb.PropertyKey_TYPE_UNSPECIFIED {
						errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.%d.key.type", i), ErrIsRequired)
					}
				}

				if len(filter.Conditions) != 0 {
					for j, condition := range filter.Conditions {
						if len(condition.Value) == 0 {
							errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.%d.conditions.%d.value", i, j), ErrIsRequired)
						}
						if condition.Operator == oraclespb.Condition_OPERATOR_UNSPECIFIED {
							errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.%d.conditions.%d.operator", i, j), ErrIsRequired)
						}
					}
				}
			}
		}
	} else {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec", ErrIsRequired)
	}

	if future.OracleSpecBinding != nil {
		if len(future.OracleSpecBinding.SettlementPriceProperty) == 0 {
			errs.AddForProperty("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec_binding.settlement_price_property", ErrIsRequired)
		}
	} else {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec_binding", ErrIsRequired)
	}

	return errs
}

func checkTradingMode(config *typespb.NewMarketConfiguration) Errors {
	errs := NewErrors()

	if config.TradingMode == nil {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.trading_mode", ErrIsRequired)
	}

	switch mode := config.TradingMode.(type) {
	case *typespb.NewMarketConfiguration_Continuous:
		errs.Merge(checkContinuousTradingMode(mode))
	case *typespb.NewMarketConfiguration_Discrete:
		errs.Merge(checkDiscreteTradingMode(mode))
	default:
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.trading_mode", ErrIsNotValid)
	}

	return errs
}

func checkContinuousTradingMode(mode *typespb.NewMarketConfiguration_Continuous) Errors {
	errs := NewErrors()

	if mode.Continuous == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_market.changes.trading_mode.continuous", ErrIsRequired)
	}

	return errs
}

func checkDiscreteTradingMode(mode *typespb.NewMarketConfiguration_Discrete) Errors {
	errs := NewErrors()

	if mode.Discrete == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_market.changes.trading_mode.discrete", ErrIsRequired)
	}

	if mode.Discrete.DurationNs <= 0 || mode.Discrete.DurationNs >= MaxDuration30DaysNs {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.trading_mode.discrete.duration_ns",
			fmt.Errorf(fmt.Sprintf("should be between 0 (excluded) and %d (excluded)", MaxDuration30DaysNs)))
	}

	return errs
}

func checkRiskParameters(config *typespb.NewMarketConfiguration) Errors {
	errs := NewErrors()

	if config.RiskParameters == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_market.changes.risk_parameters", ErrIsRequired)
	}

	switch parameters := config.RiskParameters.(type) {
	case *typespb.NewMarketConfiguration_Simple:
		errs.Merge(checkSimpleParameters(parameters))
	case *typespb.NewMarketConfiguration_LogNormal:
		errs.Merge(checkLogNormalRiskParameters(parameters))
	default:
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.risk_parameters", ErrIsNotValid)
	}

	return errs
}

func checkSimpleParameters(params *typespb.NewMarketConfiguration_Simple) Errors {
	errs := NewErrors()

	if params.Simple == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_market.changes.risk_parameters.simple", ErrIsRequired)
	}

	if params.Simple.MinMoveDown > 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.risk_parameters.simple.min_move_down", ErrMustBeNegativeOrZero)
	}

	if params.Simple.MaxMoveUp < 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.risk_parameters.simple.max_move_up", ErrMustBePositiveOrZero)
	}

	if params.Simple.ProbabilityOfTrading < 0 || params.Simple.ProbabilityOfTrading > 1 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.changes.risk_parameters.simple.probability_of_trading",
			fmt.Errorf("should be between 0 (inclusive) and 1 (inclusive)"),
		)
	}

	return errs
}

func checkLogNormalRiskParameters(params *typespb.NewMarketConfiguration_LogNormal) Errors {
	errs := NewErrors()

	if params.LogNormal == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_market.changes.risk_parameters.log_normal", ErrIsRequired)
	}

	if params.LogNormal.Params == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_market.changes.risk_parameters.log_normal.params", ErrIsRequired)
	}

	return errs
}

func checkLiquidityCommitment(commitment *typespb.NewMarketCommitment) Errors {
	errs := NewErrors()

	if commitment == nil {
		return errs.FinalAddForProperty("proposal_submission.terms.change.new_market.liquidity_commitment", ErrIsRequired)
	}

	if commitment.CommitmentAmount == 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.liquidity_commitment.commitment_amount", ErrMustBePositive)
	}
	if len(commitment.Fee) == 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.liquidity_commitment.fee", ErrIsRequired)
	}
	fee, err := strconv.ParseFloat(commitment.Fee, 64)
	if err != nil {
		errs.AddForProperty("proposal_submission.terms.change.new_market.liquidity_commitment.fee", ErrIsNotValidNumber)
	} else if fee < 0 {
		errs.AddForProperty("proposal_submission.terms.change.new_market.liquidity_commitment.fee", ErrMustBePositiveOrZero)
	}

	errs.Merge(checkShape(commitment.Buys, typespb.Side_SIDE_BUY))
	errs.Merge(checkShape(commitment.Sells, typespb.Side_SIDE_SELL))

	return errs
}

func checkShape(orders []*typespb.LiquidityOrder, side typespb.Side) Errors {
	errs := NewErrors()

	humanizedSide := "buys"
	if side == typespb.Side_SIDE_SELL {
		humanizedSide = "sells"
	}

	if len(orders) == 0 {
		return errs.FinalAddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_asset.liquidity_commitment.%s", humanizedSide), ErrIsRequired)
	}

	for i, order := range orders {
		if order.Reference == typespb.PeggedReference_PEGGED_REFERENCE_UNSPECIFIED {
			errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_asset.liquidity_commitment.%s.reference.%d", humanizedSide, i), ErrIsRequired)
		}
		if _, ok := typespb.PeggedReference_name[int32(order.Reference)]; !ok {
			errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_asset.liquidity_commitment.%s.reference.%d", humanizedSide, i), ErrIsNotValid)
		}

		if order.Proportion == 0 {
			errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_asset.liquidity_commitment.%s.proportion.%d", humanizedSide, i), ErrIsRequired)
		}

		if side == typespb.Side_SIDE_BUY {
			switch order.Reference {
			case typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK:
				errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_asset.liquidity_commitment.%s.reference.%d", humanizedSide, i),
					errors.New("cannot have a reference of type BEST_ASK when on BUY side"),
				)
			case typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID:
				if order.Offset > 0 {
					errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_asset.liquidity_commitment.%s.offset.%d", humanizedSide, i), ErrMustBeNegativeOrZero)
				}
			case typespb.PeggedReference_PEGGED_REFERENCE_MID:
				if order.Offset >= 0 {
					errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_asset.liquidity_commitment.%s.offset.%d", humanizedSide, i), ErrMustBeNegative)
				}
			}
			continue
		}

		switch order.Reference {
		case typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID:
			errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_asset.liquidity_commitment.%s.reference.%d", humanizedSide, i),
				errors.New("cannot have a reference of type BEST_BID when on SELL side"),
			)
		case typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK:
			if order.Offset < 0 {
				errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_asset.liquidity_commitment.%s.offset.%d", humanizedSide, i), ErrMustBePositiveOrZero)
			}
		case typespb.PeggedReference_PEGGED_REFERENCE_MID:
			if order.Offset <= 0 {
				errs.AddForProperty(fmt.Sprintf("proposal_submission.terms.change.new_asset.liquidity_commitment.%s.offset.%d", humanizedSide, i), ErrMustBePositive)
			}
		}

	}

	return errs
}
