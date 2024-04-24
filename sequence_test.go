package factorybot_test

import (
	"fmt"
	"sync"
	"testing"

	. "factorybot"

	"github.com/stretchr/testify/require"
)

func TestSequence(t *testing.T) {
	t.Run("N", func(t *testing.T) {
		s := NewSequence()
		require.Equal(t, 1, s.N())
		require.Equal(t, 2, s.N())
	})

	t.Run("Rewind", func(t *testing.T) {
		s := NewSequence()
		require.Equal(t, 1, s.N())
		require.Equal(t, 2, s.N())
		s.Rewind()
		require.Equal(t, 1, s.N())
	})

	t.Run("One", func(t *testing.T) {
		s := NewSequence()
		require.Equal(t, 1, s.One().(int))
	})

	t.Run("One/Complex", func(t *testing.T) {
		s := NewSequence(func(n int) interface{} {
			return fmt.Sprintf("address%d@example.com", n)
		})

		require.Equal(t, "address1@example.com", s.One().(string))
	})

	t.Run("Race", func(t *testing.T) {
		s := NewSequence()
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(2)
			go func() {
				defer wg.Done()
				_ = s.One()
			}()
			go func() {
				defer wg.Done()
				_ = s.N()
			}()
		}
		wg.Wait()
	})
}
