package scrapper

type TradeDirection int

const (
	Short TradeDirection = iota + 1
	Long
)

func (pd TradeDirection) String() string {
	if pd == Short {
		return "SHORT"
	}
	return "LONG"
}
