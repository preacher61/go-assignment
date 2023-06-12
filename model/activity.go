package model

// Activity represents the response object
// of the activity api.
type Activity struct {
	Activity string `json:"activity"`
	Key      string `json:"key"`
}
