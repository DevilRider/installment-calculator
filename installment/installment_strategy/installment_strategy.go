package installment_strategy

type Strategy string

const (
	Annuity  Strategy = "Annuity"  // 等额本息
	Linear   Strategy = "Linear"   // 等额本金
	Outright Strategy = "Outright" // 到期还本付息
)
