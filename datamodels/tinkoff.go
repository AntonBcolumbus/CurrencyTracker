package datamodels

// TinkoffData ...
type TinkoffData struct {
	ResultCode string  `json:"resultCode,omitempty"`
	Payload    Payload `json:"payload,omitempty"`
	TrackingID string  `json:"trackingId,omitempty"`
}

// Payload ..
type Payload struct {
	LastUpdate LastUpdate `json:"lastUpdate,omitempty"`
	Rates      []Rate     `json:"rates,omitempty"`
}

// LastUpdate ...
type LastUpdate struct {
	Milliseconds uint64 `json:"milliseconds,omitempty"`
}

// Rate ...
type Rate struct {
	Category     string   `json:"category,omitempty"`
	FromCurrency Currency `json:"fromCurrency,omitempty"`
	ToCurrency   Currency `json:"toCurrency,omitempty"`
	Buy          float64  `json:"buy,omitempty"`
	Sell         float64  `json:"sell,omitempty"`
}

// Currency ...
type Currency struct {
	Name string `json:"name,omitempty"`
	Code int    `json:"code,omitempty"`
}
