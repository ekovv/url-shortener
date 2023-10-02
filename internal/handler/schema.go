package handler

type uriJSON struct {
	URI string `json:"url,omitempty"`
	Res string `json:"result"`
}

type jBatch struct {
	ID     string `json:"correlation_id"`
	Origin string `json:"original_url,omitempty"`
	Short  string `json:"short_url,omitempty"`
}

type jBatchRes struct {
	ID    string `json:"correlation_id"`
	Short string `json:"short_url"`
}
