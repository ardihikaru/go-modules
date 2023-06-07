package query

// FilterListParams defines the captured filter query parameters
type FilterListParams struct {
	Ids []string `json:"id"`
}

// FilterQueryParams defines the captured filter query parameters
type FilterQueryParams struct {
	Keyword string `json:"q"`
}

// maps valid query order
const (
	ASC  string = "ASC"
	DESC string = "DESC"
)

// GetOrderMap returns a boolean value to verify if the order valid or not
func GetOrderMap() map[string]bool {
	return map[string]bool{
		ASC:  true,
		DESC: true,
	}
}
