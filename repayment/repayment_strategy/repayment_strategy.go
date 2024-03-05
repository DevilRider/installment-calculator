package repayment_strategy

type Strategy string

const (
	Payoff           Strategy = "Payoff"
	PartialRepayment Strategy = "PartialRepayment"
)
