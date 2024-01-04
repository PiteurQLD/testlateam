package scrapper

import (
	"errors"
	"math"
)

var (
	ErrNoPreviousPosition = errors.New("no previous position amount")
)

type PositionType int

const (
	Opened PositionType = iota + 1
	Closed
	AddedTo
	PartiallyClosed
)

func (pt PositionType) String() string {
	switch pt {
	default:
		return ""
	case Opened:
		return "opened"
	case Closed:
		return "closed"
	case AddedTo:
		return "added to"
	case PartiallyClosed:
		return "partially closed"
	}
}

type Position struct {
	Type       PositionType
	Direction  TradeDirection
	Ticker     string
	EntryPrice float64
	MarkPrice  float64
	Amount     float64
	PrevAmount float64
	Leverage   int
	Pnl        float64
	Roe        float64
}

func DeterminePositionType(amt float64, prevAmt float64) PositionType {
	if prevAmt == 0 {
		return Opened
	}
	if prevAmt < amt {
		return AddedTo
	}
	if amt == 0 {
		return Closed
	}
	return 0
}
func newPosition(rp rawPosition) Position {
	dir := Long
	if rp.Amount < 0 {
		dir = Short
		rp.Amount = math.Abs(rp.Amount)
	}
	return Position{
		Direction:  dir,
		Ticker:     rp.Symbol,
		EntryPrice: rp.EntryPrice,
		MarkPrice:  rp.MarkPrice,
		Amount:     rp.Amount,
		Leverage:   rp.Leverage,
		Pnl:        rp.Pnl,
		Roe:        rp.Roe,
	}
}
