package stream

import js "github.com/nats-io/nats.go/jetstream"

type AddQueryRule struct {
	GoodID    string `json:"good_id"`
	Operator  string `json:"operator"`
	Threshold int    `json:"threshold"`
}

const QueryRuleSubject = "notification.rules"

var QueryRuleStreamConfig = js.StreamConfig{
	Name:      "NOTIFICATION_RULES",
	Subjects:  []string{QueryRuleSubject},
	Retention: js.LimitsPolicy,
	Storage:   js.FileStorage,
}
