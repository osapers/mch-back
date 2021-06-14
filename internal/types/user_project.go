package types

type UserToProject struct {
	tableName struct{} `pg:"mch.user_project,alias:up"` // nolint:structcheck,unused
	UserID    string   `json:"user_id" pg:"type:uuid,pk"`
	ProjectID string   `json:"project_id" pg:"type:uuid,pk"`
	Viewed    bool     `json:"viewed"`
	Applied   bool     `json:"applied"`
	Confirmed bool     `json:"confirmed"`
	Rejected  bool     `json:"rejected"`
}
