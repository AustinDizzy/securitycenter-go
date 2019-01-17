package sc

// Zone https://docs.tenable.com/sccv/api/Scan-Zone.html
type Zone struct {
	ID          string `json:"id" sc:"id"`
	Name        string `json:"name" sc:"name"`
	Description string `json:"description" sc:"description"`
	Status      string `json:"status" sc:"status"`
}
