package types

type Tag struct {
	tableName struct{} `pg:"mch.tag,alias:tag"` // nolint:structcheck,unused
	ID        int      `json:"id"`
	Name      string   `json:"name"`
}

func GetTagsIntersection(t1, t2 []string) []string {
	mapT1, mapT2 := map[string]struct{}{}, map[string]struct{}{}
	var res []string

	for _, tag := range t1 {
		mapT1[tag] = struct{}{}
	}

	for _, tag := range t2 {
		mapT2[tag] = struct{}{}
	}

	for k := range mapT1 {
		if _, ok := mapT2[k]; ok {
			res = append(res, k)
		}
	}

	return res
}
