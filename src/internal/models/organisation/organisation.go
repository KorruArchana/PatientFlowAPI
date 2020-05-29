package organisation

// // Organisation holds the organisation records data
// type Organisation struct {
// 	ID         string `json:"id"`
// 	Name       string `json:"name"`
// 	SiteNumber string `json:"siteNumber"`
// }

// Organisation holds the organisation records data
type Organisation struct {
	SystemType       string `json:"systemType"`
	OrganisationName string `json:"organisationName"`
	PK               string `json:"pk"`
	SK               string `json:"sk"`
}
