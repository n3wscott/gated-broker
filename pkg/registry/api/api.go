package api

type LightRequest struct {
	Token     string  `json:"token"`
	Intensity float32 `json:"intensity"`
}
