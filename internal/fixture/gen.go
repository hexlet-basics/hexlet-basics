package fixture

// labels_gen.go gives every label in fixtures/ a Label variable of its own, so
// a test names a row with a compile-checked identifier instead of a string.
// `mise run gen:fixtures` regenerates it, and gen:check fails when it is stale.
//go:generate go run ../../cmd/fixturelabels -dir ../../fixtures -out labels_gen.go
