package types

const (
	RoleUnknown     string = "unknown"
	RoleOwner       string = "owner"
	RoleParticipant string = "participant"
	RoleApplicant   string = "applicant"
)

type Project struct {
	tableName        struct{}        `pg:"mch.project,alias:p"` // nolint:structcheck,unused
	ID               string          `json:"id" pg:"type:uuid"`
	Name             string          `json:"name"`
	Image            string          `json:"image"`
	Industry         string          `json:"industry"`
	OwnerID          string          `json:"owner_id"`
	Owner            ProjectOwner    `json:"owner" pg:"-"`
	Members          []ProjectMember `json:"members" pg:"-"`
	Description      string          `json:"description"`
	ReadinessStage   string          `json:"readiness_stage"`
	Tags             []string        `json:"tags" pg:",array"`
	TagsIntersection []string        `json:"tags_intersection" pg:"-"`
	LaunchDate       int64           `json:"launch_date"`
}

type ProjectOwner struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Image     string `json:"image"`
}

type ProjectMember struct {
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

func (p *Project) GetRole(userID string, appliedProjects map[string]*UserToProject) string {
	if p.OwnerID == userID {
		return RoleOwner
	}

	ap, ok := appliedProjects[p.ID]

	if ok && ap.Confirmed {
		return RoleParticipant
	}

	if ok && ap.Applied {
		return RoleApplicant
	}

	return RoleUnknown
}

func (p *Project) ToJson() *Project {
	project := *p
	project.Industry = ProjectIndustries[p.Industry]

	return &project
}

var ProjectReadinessStages = map[string]string{
	"idea":    "идея",
	"demo":    "демонстрационный образец",
	"product": "продукт",
	"scaling": "масштабирование",
}

var ProjectIndustries = map[string]string{
	"space":                 "космос",
	"aviation":              "авиация",
	"agriculture":           "сельское хозяйство",
	"biotech":               "биотех",
	"nuclearTech":           "ядертех",
	"building":              "строительство",
	"fashion":               "мода",
	"energetics":            "энергетика",
	"finance":               "финансы",
	"food":                  "питание",
	"games":                 "игры",
	"medicine":              "медицина",
	"law":                   "юриспруденция",
	"shipping":              "судоходство",
	"nanotechnology":        "нанотехнологии",
	"oilAndGas":             "нефть и газ",
	"trade":                 "торговля",
	"safety":                "безопасность",
	"it":                    "IT",
	"education":             "образование",
	"ecology":               "экология",
	"lifestyle":             "образ жизни",
	"mechanicalEngineering": "машиностроение",
	"microelectronics":      "микроэлектроника",
	"other":                 "другое",
}
