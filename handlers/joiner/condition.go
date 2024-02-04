package joiner

// Condition evaluates the sources with the criteria.
type Condition interface {
	Evaluate(sources []string) bool
}

// pointer to struct condition implements the interface Condition.
type condition struct {
	evaluate func(sources []string) bool
}

func (c *condition) Evaluate(sources []string) bool {
	return c.evaluate(sources)
}

func (c *condition) And(conditions ...Condition) *condition {
	return &condition{
		evaluate: func(sources []string) bool {
			if !c.evaluate(sources) {
				return false
			}
			for _, cond := range conditions {
				if !cond.Evaluate(sources) {
					return false
				}
			}

			return true
		},
	}
}

func (c *condition) Or(conditions ...Condition) *condition {
	return &condition{
		evaluate: func(sources []string) bool {
			if c.evaluate(sources) {
				return true
			}
			for _, cond := range conditions {
				if cond.Evaluate(sources) {
					return true
				}
			}

			return false
		},
	}
}

func (c *condition) XOr(other Condition) *condition {
	return &condition{
		evaluate: func(sources []string) bool {
			return c.Evaluate(sources) != other.Evaluate(sources)
		},
	}
}

// MatchAll returns a Condition to verify that all required sources are present.
func MatchAll(requiredSources ...string) *condition {
	if len(requiredSources) == 0 {
		return &condition{evaluate: alwaysTrue}
	}

	return &condition{
		evaluate: func(sources []string) bool {
			if len(sources) < len(requiredSources) {
				return false
			}

			isSourcePresent := make(map[string]bool)
			for _, source := range sources {
				isSourcePresent[source] = true
			}
			for _, requiredSource := range requiredSources {
				if !isSourcePresent[requiredSource] {
					return false
				}
			}

			return true
		},
	}
}

// MatchAny returns a Condition to verify that any sources-to-match are present.
func MatchAny(sourcesToMatch ...string) *condition {
	if len(sourcesToMatch) == 0 {
		return &condition{evaluate: alwaysTrue}
	}
	isSourceMatched := make(map[string]bool)
	for _, source := range sourcesToMatch {
		isSourceMatched[source] = true
	}

	return &condition{
		evaluate: func(sources []string) bool {
			for _, source := range sources {
				if isSourceMatched[source] {
					return true
				}
			}

			return false
		},
	}
}

// MatchNone returns a Condition to verify that none of the sources-to-exclude are present.
// note that this condition would pass if input sources is empty or nil.
func MatchNone(sourcesToExclude ...string) *condition {
	if len(sourcesToExclude) == 0 {
		return &condition{evaluate: alwaysTrue}
	}
	isSourceUnexpected := make(map[string]bool)
	for _, source := range sourcesToExclude {
		isSourceUnexpected[source] = true
	}

	return &condition{
		evaluate: func(sources []string) bool {
			for _, source := range sources {
				if isSourceUnexpected[source] {
					return false
				}
			}

			return true
		},
	}
}

func alwaysTrue(_ []string) bool {
	return true
}
