package scrapper

import (
	"context"
	"fmt"
	"time"
)

func (u *User) SubscribePositions(ctx context.Context) (<-chan Position, <-chan error) {
	cp := make(chan Position)
	ce := make(chan error)
	go func() {
		defer close(cp)
		defer close(ce)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				res, err := u.GetOtherPosition(ctx)
				if err != nil {
					ce <- fmt.Errorf("failed to fetch positions: %w", err)
					time.Sleep(u.Delay())
					continue
				}
				if !res.Success {
					ce <- fmt.Errorf("failed to fetch positions, bad response message: %v", res.Message)
					time.Sleep(u.Delay())
					continue
				}

				u.handlePositions(res.Data.OtherPositionRetList, cp, ce)
				time.Sleep(u.Delay())
			}
		}
	}()
	return cp, ce
}

func (u *User) handlePositions(rps []rawPosition, cp chan<- Position, ce chan<- error) {
	used := make(map[string]struct{}, len(rps))
	for _, rp := range rps {
		p := newPosition(rp)
		used[p.Ticker] = struct{}{}
		pp := u.positions[p.Ticker]
		if pp.Amount == p.Amount {
			pp.MarkPrice = p.MarkPrice
			pp.Pnl = p.Pnl
			pp.Roe = p.Roe
			u.positions[p.Ticker] = pp
			continue
		}
		p.PrevAmount = pp.Amount
		p.Type = DeterminePositionType(p.Amount, pp.Amount)
		u.log.Printf("[%s] {send: %t} Position change: %d %s %f -> %f %s @ %f\n", u.UID, !u.firstFetch, p.Type, p.Direction, p.PrevAmount, p.Amount, p.Ticker, p.EntryPrice)
		if !u.firstFetch {
			cp <- p
		}
		u.positions[p.Ticker] = p
	}
	for h, p := range u.positions {
		if _, ok := used[h]; ok {
			continue
		}
		p.Type = Closed
		p.PrevAmount = p.Amount
		p.Amount = 0
		cp <- p
		delete(u.positions, h)
	}
	if u.firstFetch {
		u.firstFetch = false
	}
}
