package joiner

// TODO: implement me
// requirement: be able to match (AND|OR|()) operations
type Condition interface {
	Match(sources []string) bool
}
