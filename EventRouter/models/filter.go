package models

// swagger:enum EventType
type EventType string

const (
	COMMIT_EVENT EventType = "COMMIT_EVENT"
	COMMIT_ENTRY EventType = "COMMIT_ENTRY"
	ANCHOR_EVENT EventType = "ANCHOR_EVENT"
	REVEAL_ENTRY EventType = "REVEAL_ENTRY"
	NODE_MESSAGE EventType = "NODE_MESSAGE"
)

//  Define a filter with GraphQL
// swagger:model GraphQL
type GraphQL string

// Define a Filter on an EventType to filter the event. This allows to reduce the network traffic
// swagger:model Filter
type Filter struct {
	// Filtering with graph ql
	// required: false
	Filtering GraphQL
}