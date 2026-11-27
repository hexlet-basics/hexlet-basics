package domain

// Page is a window of rows, as an offset and a size.
type Page struct {
	Offset int
	Limit  int
}
