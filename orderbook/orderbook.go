package orderbook

type Order struct {
	ID        int64
	UserID    int64
	Size      float64
	Bid       float64 // limit or market
	Limit     *Limit
	Timestamp int64
}

type Limit struct {
}
