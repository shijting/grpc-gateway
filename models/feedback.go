package models

type Feedback struct {
	tableName   struct{} `pg:"feedback"`
	Id          int32    `pg:",pk`
	Content     string   `pg:",notnull"`
	PhoneNumber string   `pg:",notnull"`
}

func NewFeedback(content string, phoneNumber string) *Feedback {
	return &Feedback{Content: content, PhoneNumber: phoneNumber}
}
