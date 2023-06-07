package httputils

// GetQueryParams defines parameters for GetUsers.
type GetQueryParams struct {
	// Maximum number of docs to return
	Limit int64 `json:"limit,omitempty"`

	// Number of docs to skip
	Offset int64 `json:"offset,omitempty"`

	// order target field
	Order string `json:"order,omitempty"`

	// sort target by field
	Sort string `json:"sort,omitempty"`

	// sort target by field
	Search string `json:"search,omitempty"`

	// filter target by field
	Filter string `json:"filter,omitempty"`
}
