package podcasts

type Owner struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Copyright string `json:"copyright"`
}
