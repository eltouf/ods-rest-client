package odsrestclient

//Result : result of an API call
type Result struct {
	Nhits      uint
	Parameters struct {
		Timezone string
		Rows     uint
		Format   string
		Staged   bool
	}
}
