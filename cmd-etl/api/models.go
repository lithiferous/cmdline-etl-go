package api

type SnapshotRow struct {
	UserId      int     `csv:"user_id"`
	Name        string  `csv:"name"`
	StoreName   string  `csv:"store_name"`
	CreditLimit float64 `csv:"credit_limit"`
}
