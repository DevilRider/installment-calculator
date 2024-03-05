package readjustment_strategy

type Strategy string

const (
	ReduceMonthlyRepayment Strategy = "ReduceMonthlyRepayment"
	ShortenRepaymentTenors Strategy = "ShortenRepaymentTenors" // 月供不变仅对 仅等额本息生效
)
