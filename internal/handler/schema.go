package handler

type uriJSON struct {
	URI string `json:"url,omitempty"`
	Res string `json:"result"`
}

type jBatch struct {
	Id     string `json:"correlation_id"`
	Origin string `json:"original_url"`
}

type jBatchRes struct {
	Id    string `json:"correlation_id"`
	Short string `json:"short_url"`
}
