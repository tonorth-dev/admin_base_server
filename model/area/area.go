package area

type Division struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Province struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type City struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Province string `json:"province"`
}

type County struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}
