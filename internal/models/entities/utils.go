package entities

func Flat(ens []*DBEntity) []map[string]interface{} {
	m := make([]map[string]interface{}, len(ens))

	for i, e := range ens {
		m[i] = e.Flat()
	}

	return m
}
