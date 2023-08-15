package query

import "encoding/json"

// FilterListParams defines the captured filter query parameters
type FilterListParams struct {
	Ids []string `json:"id"`
}

// FilterQueryParams defines the captured filter query parameters
type FilterQueryParams struct {
	Keyword    string               `json:"q"`
	Filter     *map[string]string   `json:"filter"`
	FilterList *map[string][]string `json:"filter_list"`
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

// GetFilterQuery extracts filter from the url query
func GetFilterQuery(filter string, destType interface{}) error {
	// in some cases, it will ignore it
	if filter == "{}" {
		return nil
	}

	// extracts to designated variable
	err := json.Unmarshal([]byte(filter), destType)
	if err != nil {
		return err
	}

	return nil
}
