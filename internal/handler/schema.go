package handler

type uriJSON struct {
	URI string `json:"url,omitempty"`
	Res string `json:"result"`
}

type jBatch struct {
	ID     string `json:"correlation_id,omitempty"`
	Origin string `json:"original_url,omitempty"`
	Short  string `json:"short_url,omitempty"`
}
