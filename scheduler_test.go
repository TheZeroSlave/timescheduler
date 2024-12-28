package timescheduler

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestScheduler_Stops(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := NewScheduler(ctx)
	called := false
	s.Add(
		time.Now().Add(time.Second),
		func() {
			t.Log("since 1 second")
			called = true
		},
	)
	s.Add(
		time.Now().Add(time.Second*3),
		func() {
			require.Fail(t, "function should not be called, timer stop isn't working?")
		},
	)
	time.Sleep(time.Second * 2)
	cancel()
	time.Sleep(time.Second * 2)
	require.True(t, called)
}

func TestScheduler_AddRecent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()
	sinceOne := time.Now().Add(time.Second)
	sinceTwo := time.Now().Add(time.Second * 2)
	sinceThree := time.Now().Add(time.Second * 3)

	s := NewScheduler(ctx)
	calledOne := false
	calledTwo := false
	s.Add(
		sinceOne,
		func() {
			t.Log("since 1")
			calledOne = true
		},
	)
	s.Add(
		sinceThree,
		func() {
			t.Log("since 3")
			require.Fail(t, "shouldn't be called")
		},
	)

	time.Sleep(time.Second*1 + 100*time.Millisecond)
	require.True(t, calledOne)

	s.Add(
		sinceTwo,
		func() {
			t.Log("since 2")
			calledTwo = true
		},
	)
	time.Sleep(time.Second + 100*time.Millisecond)
	require.True(t, calledTwo)
}

func TestScheduler_Complete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()
	sinceOne := time.Now().Add(time.Second)
	sinceTwo := time.Now().Add(time.Second * 2)
	sinceThree := time.Now().Add(time.Second * 3)

	s := NewScheduler(ctx)
	wg := sync.WaitGroup{}
	wg.Add(3)
	calls := []int{}
	s.Add(
		sinceOne,
		func() {
			t.Log("since 1")
			calls = append(calls, 1)
			wg.Done()
		},
	)
	s.Add(
		sinceThree,
		func() {
			t.Log("since 3")
			calls = append(calls, 3)
			wg.Done()
		},
	)

	time.Sleep(time.Second*1 + 100*time.Millisecond)

	s.Add(
		sinceTwo,
		func() {
			t.Log("since 2")
			calls = append(calls, 2)
			wg.Done()
		},
	)
	wg.Wait()
	require.EqualValues(t, calls, []int{1, 2, 3})
}
