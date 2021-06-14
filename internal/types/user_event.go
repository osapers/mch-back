package types

type UserToEvent struct {
	tableName struct{} `pg:"mch.user_event,alias:ue"` // nolint:structcheck,unused
	UserID    string   `json:"user_id" pg:"type:uuid,pk"`
	EventID   string   `json:"event_id" pg:"type:uuid,pk"`
}
