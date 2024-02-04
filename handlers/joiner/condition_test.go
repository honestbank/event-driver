package joiner_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lukecold/event-driver/handlers/joiner"
)

type conditionTestCase struct {
	sources     []string
	condition   joiner.Condition
	shouldMatch bool
}

func TestMatchAll(t *testing.T) {
	matchBoth := joiner.MatchAll("source1", "source2")

	testCases := map[string]conditionTestCase{
		"all required sources present": {
			sources:     []string{"source1", "source2"},
			condition:   matchBoth,
			shouldMatch: true,
		},
		"additional sources present": {
			sources:     []string{"other", "source1", "source2"},
			condition:   matchBoth,
			shouldMatch: true,
		},
		"no required sources": {
			sources:     []string{"any"},
			condition:   joiner.MatchAll(),
			shouldMatch: true,
		},
		"required source(s) missing": {
			sources:     []string{"source1"},
			condition:   matchBoth,
			shouldMatch: false,
		},
		"empty input": {
			sources:     []string{},
			condition:   matchBoth,
			shouldMatch: false,
		},
		"null input": {
			sources:     nil,
			condition:   matchBoth,
			shouldMatch: false,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			isMatch := testCase.condition.Evaluate(testCase.sources)
			assert.Equal(t, testCase.shouldMatch, isMatch)
		})
	}
}

func TestMatchAny(t *testing.T) {
	matchEither := joiner.MatchAny("source1", "source2")

	testCases := map[string]conditionTestCase{
		"all sources present": {
			sources:     []string{"source1", "source2"},
			condition:   matchEither,
			shouldMatch: true,
		},
		"one of the sources present": {
			sources:     []string{"source2"},
			condition:   matchEither,
			shouldMatch: true,
		},
		"additional sources present": {
			sources:     []string{"other", "source1", "source2"},
			condition:   matchEither,
			shouldMatch: true,
		},
		"no sources to match": {
			sources:     []string{"any"},
			condition:   joiner.MatchAny(),
			shouldMatch: true,
		},
		"empty input": {
			sources:     []string{},
			condition:   matchEither,
			shouldMatch: false,
		},
		"null input": {
			sources:     nil,
			condition:   matchEither,
			shouldMatch: false,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			isMatch := testCase.condition.Evaluate(testCase.sources)
			assert.Equal(t, testCase.shouldMatch, isMatch)
		})
	}
}

func TestMatchNone(t *testing.T) {
	matchNeither := joiner.MatchNone("source1", "source2")

	testCases := map[string]conditionTestCase{
		"one unexpected sources present": {
			sources:     []string{"source1", "other"},
			condition:   matchNeither,
			shouldMatch: false,
		},
		"all unexpected sources present": {
			sources:     []string{"source1", "source2", "other"},
			condition:   matchNeither,
			shouldMatch: false,
		},
		"no sources to exclude": {
			sources:     []string{"any"},
			condition:   joiner.MatchNone(),
			shouldMatch: true,
		},
		"none of the unexpected sources present": {
			sources:     []string{"other", "more"},
			condition:   matchNeither,
			shouldMatch: true,
		},
		"empty input": {
			sources:     []string{},
			condition:   matchNeither,
			shouldMatch: true,
		},
		"null input": {
			sources:     nil,
			condition:   matchNeither,
			shouldMatch: true,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			isMatch := testCase.condition.Evaluate(testCase.sources)
			assert.Equal(t, testCase.shouldMatch, isMatch)
		})
	}
}

func TestLogicOperations(t *testing.T) {
	passCondition := joiner.MatchAny()
	failCondition := joiner.MatchAll("not-exist")
	var sources []string
	assert.True(t, passCondition.Evaluate(sources))
	assert.False(t, failCondition.Evaluate(sources))

	testCases := map[string]conditionTestCase{
		"true AND true = true": {
			condition:   passCondition.And(passCondition),
			shouldMatch: true,
		},
		"true AND false = false": {
			condition:   passCondition.And(failCondition),
			shouldMatch: false,
		},
		"false AND true = false": {
			condition:   failCondition.And(passCondition),
			shouldMatch: false,
		},
		"false AND false = false": {
			condition:   failCondition.And(failCondition),
			shouldMatch: false,
		},
		"true AND true AND true = true": {
			condition:   passCondition.And(passCondition, passCondition),
			shouldMatch: true,
		},
		"true AND true AND false = false": {
			condition:   passCondition.And(passCondition, failCondition),
			shouldMatch: false,
		},
		"true OR true = true": {
			condition:   passCondition.Or(passCondition),
			shouldMatch: true,
		},
		"true OR false = true": {
			condition:   passCondition.Or(failCondition),
			shouldMatch: true,
		},
		"false OR true = true": {
			condition:   failCondition.Or(passCondition),
			shouldMatch: true,
		},
		"false OR false = false": {
			condition:   failCondition.Or(failCondition),
			shouldMatch: false,
		},
		"true OR true OR true = true": {
			condition:   passCondition.Or(passCondition, passCondition),
			shouldMatch: true,
		},
		"true OR true OR false = true": {
			condition:   passCondition.Or(passCondition, failCondition),
			shouldMatch: true,
		},
		"false OR false OR false = false": {
			condition:   failCondition.Or(failCondition, failCondition),
			shouldMatch: false,
		},
		"true XOR true = false": {
			condition:   passCondition.XOr(passCondition),
			shouldMatch: false,
		},
		"true XOR false = true": {
			condition:   passCondition.XOr(failCondition),
			shouldMatch: true,
		},
		"false XOR true = true": {
			condition:   failCondition.XOr(passCondition),
			shouldMatch: true,
		},
		"false XOR false = false": {
			condition:   failCondition.XOr(failCondition),
			shouldMatch: false,
		},
		"false AND true OR true = true": {
			condition:   failCondition.And(passCondition).Or(passCondition),
			shouldMatch: true,
		},
		"false AND (true OR true) = false": {
			condition:   failCondition.And(passCondition.Or(passCondition)),
			shouldMatch: false,
		},
		"false AND true OR true XOR true = false": {
			condition:   failCondition.And(passCondition).Or(passCondition).XOr(passCondition),
			shouldMatch: false,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			isMatch := testCase.condition.Evaluate(testCase.sources)
			assert.Equal(t, testCase.shouldMatch, isMatch)
		})
	}
}
