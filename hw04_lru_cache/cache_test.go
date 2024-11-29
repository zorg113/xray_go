package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type ComplexValue struct {
	name string
	val  int
}

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(5)
		{
			key1 := ComplexValue{"Ivan", 500}
			key2 := ComplexValue{"Egor", 600}
			key3 := ComplexValue{"Dan", 620}
			key4 := ComplexValue{"Alex", 650}
			key5 := ComplexValue{"Rik", 700}

			_ = c.Set("key_1", key1)
			_ = c.Set("key_2", key2)
			_ = c.Set("key_3", key3)
			_ = c.Set("key_4", key4)
			_ = c.Set("key_5", key5)
		}
		val, ok := c.Get("key_5")
		require.True(t, ok)
		require.Equal(t, ComplexValue{"Rik", 700}, val)

		val, ok = c.Get("key_3")
		require.True(t, ok)
		require.Equal(t, ComplexValue{"Dan", 620}, val)

		val, ok = c.Get("key_2")
		require.True(t, ok)
		require.Equal(t, ComplexValue{"Egor", 600}, val)

		val, ok = c.Get("key_1")
		require.True(t, ok)
		require.Equal(t, ComplexValue{"Ivan", 500}, val)

		key6 := ComplexValue{"Zorg", 799}
		_ = c.Set("key_6", key6)

		val, ok = c.Get("key_4")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear", func(t *testing.T) {
		c := NewCache(5)
		{
			key1 := ComplexValue{"Ivan", 500}
			key2 := ComplexValue{"Egor", 600}
			key3 := ComplexValue{"Dan", 620}
			key4 := ComplexValue{"Alex", 650}
			key5 := ComplexValue{"Rik", 700}

			_ = c.Set("key_1", key1)
			_ = c.Set("key_2", key2)
			_ = c.Set("key_3", key3)
			_ = c.Set("key_4", key4)
			_ = c.Set("key_5", key5)
		}
		c.Clear()
		val, ok := c.Get("key_1")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
