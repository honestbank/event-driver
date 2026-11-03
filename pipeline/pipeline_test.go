package pipeline_test

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/event-driver/event"
	"github.com/honestbank/event-driver/handlers"
	"github.com/honestbank/event-driver/pipeline"
)

type testHandler struct {
	processTime  time.Duration
	err          error
	rethrowError func(error) error
}

func createHandler(processTime time.Duration) handlers.Handler {
	return &testHandler{
		processTime:  processTime,
		err:          nil,
		rethrowError: func(err error) error { return err },
	}
}

func createRethrowErrorHandler(processTime time.Duration, rethrowError func(err error) error) handlers.Handler {
	return &testHandler{
		processTime:  processTime,
		err:          nil,
		rethrowError: rethrowError,
	}
}

func createFailedHandler(processTime time.Duration, err error) handlers.Handler {
	return &testHandler{
		processTime:  processTime,
		err:          err,
		rethrowError: func(err error) error { return err },
	}
}

func (h *testHandler) Process(ctx context.Context, in *event.Message, next handlers.CallNext) error {
	time.Sleep(h.processTime)
	if h.err != nil {
		return h.err
	}
	err := next.Call(ctx, in)

	return h.rethrowError(err)
}

func TestPipelineFail(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		testPipeline := pipeline.New().
			WithNextHandler(createHandler(time.Nanosecond)).
			WithNextHandler(createHandler(time.Nanosecond))

		err := testPipeline.Process(context.Background(), nil)
		assert.NoError(t, err)
	})

	t.Run("stop processing when handler failed", func(t *testing.T) {
		expectedError := errors.New("fail")
		testPipeline := pipeline.New().
			WithNextHandler(createHandler(time.Nanosecond)).
			WithNextHandler(createFailedHandler(time.Nanosecond, expectedError)).
			WithNextHandler(createFailedHandler(time.Nanosecond, errors.New("other error")))

		err := testPipeline.Process(context.Background(), nil)
		assert.Equal(t, expectedError, err)
	})

	t.Run("rethrow error from downstream", func(t *testing.T) {
		expectedError := errors.New("rethrown error")
		testPipeline := pipeline.New().
			WithNextHandler(createRethrowErrorHandler(time.Nanosecond, func(err error) error {
				return expectedError
			})).
			WithNextHandler(createFailedHandler(time.Nanosecond, errors.New("original error")))

		err := testPipeline.Process(context.Background(), nil)
		assert.Equal(t, expectedError, err)
	})
}

func TestPipelineTimeout(t *testing.T) {
	ctx := context.Background()
	testPipeline := pipeline.New().
		WithNextHandler(createHandler(10 * time.Millisecond)).
		WithNextHandler(createHandler(100 * time.Millisecond))

	type TestCase struct {
		createCtx     func() (context.Context, context.CancelFunc)
		expectedError error
	}
	testCases := map[string]TestCase{
		"timeout at index 0": {
			createCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(ctx, time.Millisecond)
			},
			expectedError: errors.New("pipeline timed out"),
		},
		"timeout at index 1": {
			createCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(ctx, 101*time.Millisecond)
			},
			expectedError: errors.New("pipeline timed out"),
		},
		"finish within timeout": {
			createCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(ctx, 200*time.Millisecond)
			},
			expectedError: nil,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			ctxWithTimeout, _ := testCase.createCtx()
			actualError := testPipeline.Process(ctxWithTimeout, nil)
			assert.Equal(t, testCase.expectedError, actualError)
		})
	}
}

type blockingHandler struct {
	release chan struct{}
}

func (h *blockingHandler) Process(_ context.Context, _ *event.Message, _ handlers.CallNext) error {
	<-h.release

	return nil
}

func TestPipelineTimeoutDoesNotLeakGoroutine(t *testing.T) {
	handler := &blockingHandler{release: make(chan struct{})}
	testPipeline := pipeline.New().WithNextHandler(handler)
	before := runtime.NumGoroutine()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	err := testPipeline.Process(ctx, nil)
	assert.Equal(t, errors.New("pipeline timed out"), err)

	// the handler ignores the context, so it ends only after release
	close(handler.release)
	deadline := time.Now().Add(time.Second)
	for runtime.NumGoroutine() > before && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	assert.LessOrEqual(t, runtime.NumGoroutine(), before)
}
