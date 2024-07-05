package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/tokenfactory module sentinel errors
var (
	ErrInvalidSigner = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample        = sdkerrors.Register(ModuleName, 1101, "sample error")
	ErrDenomExists   = sdkerrors.Register(ModuleName, 1102, "denomination already exists")
	ErrTickerExists  = sdkerrors.Register(ModuleName, 1103, "ticker already exists")
)
