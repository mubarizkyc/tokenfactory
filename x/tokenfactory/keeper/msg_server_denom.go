package keeper

import (
	"context"
	"fmt"

	"tokenfactory/x/tokenfactory/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateDenom(goCtx context.Context, msg *types.MsgCreateDenom) (*types.MsgCreateDenomResponse, error) {
	fmt.Println("====== CREATE DENOM FUNCTION CALLED ======")
	fmt.Printf("Attempting to create denom: %s\n", msg.Denom)
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the denom already exists
	_, isFound := k.GetDenom(ctx, msg.Denom)

	if isFound {
		fmt.Printf("Denom '%s' found: %v\n", msg.Denom, isFound)
		return nil, errorsmod.Wrapf(types.ErrDenomExists, "denom '%s' already exists", msg.Denom)
	}

	// Check if the ticker is unique
	denoms := k.GetAllDenom(ctx)
	for _, d := range denoms {
		if d.Ticker == msg.Ticker {
			return nil, errorsmod.Wrapf(types.ErrTickerExists, "ticker '%s' is already in use by denom '%s'", msg.Ticker, d.Denom)
		}
	}

	// Create the new denom
	var denom = types.Denom{
		Owner:              msg.Owner,
		Denom:              msg.Denom,
		Description:        msg.Description,
		Ticker:             msg.Ticker,
		Precision:          msg.Precision,
		Url:                msg.Url,
		MaxSupply:          msg.MaxSupply,
		Supply:             0, // Set initial supply to 0
		CanChangeMaxSupply: msg.CanChangeMaxSupply,
	}

	fmt.Printf("Denom '%s' created successfully\n", msg.Denom)
	fmt.Println("====== CREATE DENOM FUNCTION COMPLETED ======")
	k.SetDenom(ctx, denom)
	return &types.MsgCreateDenomResponse{}, nil
}

func (k msgServer) UpdateDenom(goCtx context.Context, msg *types.MsgUpdateDenom) (*types.MsgUpdateDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetDenom(
		ctx,
		msg.Denom,
	)
	if !isFound {
		return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "Denom to update not found")
	}

	// Checks if the the msg owner is the same as the current owner
	if msg.Owner != valFound.Owner {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if !valFound.CanChangeMaxSupply && valFound.MaxSupply != msg.MaxSupply {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "cannot change maxsupply")
	}
	if !valFound.CanChangeMaxSupply && msg.CanChangeMaxSupply {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "Cannot revert change maxsupply flag")
	}
	var denom = types.Denom{
		Owner:              msg.Owner,
		Denom:              msg.Denom,
		Description:        msg.Description,
		Ticker:             valFound.Ticker,
		Precision:          valFound.Precision,
		Url:                msg.Url,
		MaxSupply:          msg.MaxSupply,
		Supply:             valFound.Supply,
		CanChangeMaxSupply: msg.CanChangeMaxSupply,
	}

	k.SetDenom(ctx, denom)

	return &types.MsgUpdateDenomResponse{}, nil
}
