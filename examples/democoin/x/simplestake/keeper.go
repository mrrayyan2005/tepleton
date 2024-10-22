package simplestake

import (
	"github.com/tepleton/tepleton/crypto"

	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
	"github.com/tepleton/tepleton-sdk/x/bank"
)

const stakingToken = "steak"

const moduleName = "simplestake"

// simple stake keeper
type Keeper struct {
	ck bank.Keeper

	key       sdk.StoreKey
	cdc       *wire.Codec
	codespace sdk.CodespaceType
}

func NewKeeper(key sdk.StoreKey, coinKeeper bank.Keeper, codespace sdk.CodespaceType) Keeper {
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	return Keeper{
		key:       key,
		cdc:       cdc,
		ck:        coinKeeper,
		codespace: codespace,
	}
}

func (k Keeper) getBondInfo(ctx sdk.Context, addr sdk.Address) bondInfo {
	store := ctx.KVStore(k.key)
	bz := store.Get(addr)
	if bz == nil {
		return bondInfo{}
	}
	var bi bondInfo
	err := k.cdc.UnmarshalBinary(bz, &bi)
	if err != nil {
		panic(err)
	}
	return bi
}

func (k Keeper) setBondInfo(ctx sdk.Context, addr sdk.Address, bi bondInfo) {
	store := ctx.KVStore(k.key)
	bz, err := k.cdc.MarshalBinary(bi)
	if err != nil {
		panic(err)
	}
	store.Set(addr, bz)
}

func (k Keeper) deleteBondInfo(ctx sdk.Context, addr sdk.Address) {
	store := ctx.KVStore(k.key)
	store.Delete(addr)
}

// register a bond with the keeper
func (k Keeper) Bond(ctx sdk.Context, addr sdk.Address, pubKey crypto.PubKey, stake sdk.Coin) (int64, sdk.Error) {
	if stake.Denom != stakingToken {
		return 0, ErrIncorrectStakingToken(k.codespace)
	}

	_, _, err := k.ck.SubtractCoins(ctx, addr, []sdk.Coin{stake})
	if err != nil {
		return 0, err
	}

	bi := k.getBondInfo(ctx, addr)
	if bi.isEmpty() {
		bi = bondInfo{
			PubKey: pubKey,
			Power:  0,
		}
	}

	bi.Power = bi.Power + stake.Amount.Int64()

	k.setBondInfo(ctx, addr, bi)
	return bi.Power, nil
}

// register an unbond with the keeper
func (k Keeper) Unbond(ctx sdk.Context, addr sdk.Address) (crypto.PubKey, int64, sdk.Error) {
	bi := k.getBondInfo(ctx, addr)
	if bi.isEmpty() {
		return nil, 0, ErrInvalidUnbond(k.codespace)
	}
	k.deleteBondInfo(ctx, addr)

	returnedBond := sdk.NewCoin(stakingToken, bi.Power)

	_, _, err := k.ck.AddCoins(ctx, addr, []sdk.Coin{returnedBond})
	if err != nil {
		return bi.PubKey, bi.Power, err
	}

	return bi.PubKey, bi.Power, nil
}

// FOR TESTING PURPOSES -------------------------------------------------

func (k Keeper) bondWithoutCoins(ctx sdk.Context, addr sdk.Address, pubKey crypto.PubKey, stake sdk.Coin) (int64, sdk.Error) {
	if stake.Denom != stakingToken {
		return 0, ErrIncorrectStakingToken(k.codespace)
	}

	bi := k.getBondInfo(ctx, addr)
	if bi.isEmpty() {
		bi = bondInfo{
			PubKey: pubKey,
			Power:  0,
		}
	}

	bi.Power = bi.Power + stake.Amount.Int64()

	k.setBondInfo(ctx, addr, bi)
	return bi.Power, nil
}

func (k Keeper) unbondWithoutCoins(ctx sdk.Context, addr sdk.Address) (crypto.PubKey, int64, sdk.Error) {
	bi := k.getBondInfo(ctx, addr)
	if bi.isEmpty() {
		return nil, 0, ErrInvalidUnbond(k.codespace)
	}
	k.deleteBondInfo(ctx, addr)

	return bi.PubKey, bi.Power, nil
}
