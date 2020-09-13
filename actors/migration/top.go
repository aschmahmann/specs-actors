package migration

import (
	"context"

	address "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	builtin0 "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/util/adt"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/specs-actors/v2/actors/builtin"
	"github.com/filecoin-project/specs-actors/v2/actors/states"
)

type StateMigration interface {
	// Loads an actor's state from an input store and writes new state to an output store.
	// Returns the new state head CID.
	MigrateState(ctx context.Context, store cbor.IpldStore, head cid.Cid) (cid.Cid, error)
}

type ActorMigration struct {
	OutCodeCID     cid.Cid
	StateMigration StateMigration
}

var migrations = map[cid.Cid]ActorMigration{ // nolint:varcheck,deadcode,unused
	builtin0.AccountActorCodeID: ActorMigration{
		OutCodeCID:     builtin.AccountActorCodeID,
		StateMigration: &accountMigrator{},
	},
	builtin0.CronActorCodeID: ActorMigration{
		OutCodeCID:     builtin.CronActorCodeID,
		StateMigration: &cronMigrator{},
	},
	builtin0.InitActorCodeID: ActorMigration{
		OutCodeCID:     builtin.InitActorCodeID,
		StateMigration: &initMigrator{},
	},
	builtin0.StorageMarketActorCodeID: ActorMigration{
		OutCodeCID:     builtin.StorageMarketActorCodeID,
		StateMigration: &marketMigrator{},
	},
	builtin0.StorageMinerActorCodeID: ActorMigration{
		OutCodeCID:     builtin.StorageMinerActorCodeID,
		StateMigration: &minerMigrator{},
	},
	builtin0.MultisigActorCodeID: ActorMigration{
		OutCodeCID:     builtin.MultisigActorCodeID,
		StateMigration: &multisigMigrator{},
	},
	builtin0.PaymentChannelActorCodeID: ActorMigration{
		OutCodeCID:     builtin.PaymentChannelActorCodeID,
		StateMigration: &paychMigrator{},
	},
	builtin0.StoragePowerActorCodeID: ActorMigration{
		OutCodeCID:     builtin.StoragePowerActorCodeID,
		StateMigration: &powerMigrator{},
	},
	builtin0.RewardActorCodeID: ActorMigration{
		OutCodeCID:     builtin.RewardActorCodeID,
		StateMigration: &rewardMigrator{},
	},
	builtin0.SystemActorCodeID: ActorMigration{
		OutCodeCID:     builtin.SystemActorCodeID,
		StateMigration: &systemMigrator{},
	},
	builtin0.VerifiedRegistryActorCodeID: ActorMigration{
		OutCodeCID:     builtin.VerifiedRegistryActorCodeID,
		StateMigration: &verifregMigrator{},
	},
}

// A phoenix tracks the unburning of funds
type phoenix struct {
	burntBalance abi.TokenAmount
}

func (p phoenix) load(ctx context.Context, actorsIn *states.TreeTop) error {
	burntFundsActor, err := actorsIn.GetActor(ctx, builtin0.BurntFundsActorAddr)
	if err != nil {
		return err
	}
	p.burntBalance = burntFundsActor.Balance
	return nil
}

func (p phoenix) curr() abi.TokenAmount {
	return p.burntBalance
}

func (p phoenix) transfer(amt abi.TokenAmount) error {
	p.burntBalance = big.Sub(p.burntBalance, amt)
	if p.burntBalance.LessThan(big.Zero()) {
		return xerrors.Errorf("migration programmer error burnt funds balance falls to %v, below zero", p.burntBalance)
	}
	return nil
}

func (p phoenix) flush(ctx context.Context, actorsIn, actorsOut *states.TreeTop) error {
	burntFundsActor, err := actorsIn.GetActor(ctx, builtin0.BurntFundsActorAddr)
	if err != nil {
		return err
	}
	burntFundsActor.Balance = p.burntBalance
	return actorsOut.SetActor(ctx, builtin.BurntFundsActorAddr, burntFundsActor)
}

// Migrates the filecoin state tree starting from the global state tree and upgrading all actor state.
func MigrateStateTree(ctx context.Context, store cbor.IpldStore, stateRootIn cid.Cid) (cid.Cid, error) {
	// first migrate the global state tree hamt to something v2 can work with
	// if this is very slow (likely) we can cherry-pick states.TreeTop to v0.9
	// to avoid this step.
	stateRootInTweaked, err := migrateHAMTRaw(ctx, store, stateRootIn)
	if err != nil {
		return cid.Undef, err
	}
	adtStore := adt.WrapStore(ctx, store)
	actorsIn, err := states.AsTreeTop(adtStore, stateRootInTweaked)
	if err != nil {
		return cid.Undef, err
	}
	var p phoenix
	if err := p.load(ctx, actorsIn); err != nil {
		return cid.Undef, err
	}

	stateRootOut, err := adt.MakeEmptyMap(adtStore).Root()
	if err != nil {
		return cid.Undef, err
	}
	actorsOut, err := states.AsTreeTop(adtStore, stateRootOut)

	// Iterate all actors in old state root
	// Set new state root actors as we go
	err = actorsIn.ForEach(ctx, func(addr address.Address, actorIn *states.Actor) error {
		migration := migrations[actorIn.Code]
		if actorIn.Code == builtin0.StorageMinerActorCodeID {
			// setup migration fields
			mm := migration.StateMigration.(*minerMigrator)
			mm.MinerBalance = actorIn.Balance
			mm.Transfer = big.Zero()
		}
		headOut, err := migration.StateMigration.MigrateState(ctx, store, actorIn.Head)
		if err != nil {
			return err
		}

		// set up new state root with the migrated state
		actorOut := states.Actor{
			Code:    migration.OutCodeCID,
			Head:    headOut,
			Nonce:   actorIn.Nonce,
			Balance: actorIn.Balance,
		}

		if actorIn.Code == builtin0.StorageMinerActorCodeID {
			// propagate transfer to miner actor
			mm := migration.StateMigration.(*minerMigrator)
			if err := p.transfer(mm.Transfer); err != nil {
				return err
			}

			actorOut.Balance = big.Add(actorOut.Balance, mm.Transfer)
		}
		return actorsOut.SetActor(ctx, addr, &actorOut)
	})
	if err != nil {
		return cid.Undef, err
	}

	// Track deductions to burntFunds actor's balance
	if err := p.flush(ctx, actorsIn, actorsOut); err != nil {
		return cid.Undef, err
	}

	return actorsOut.Root()
}