package orderedmap_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/elliotchance/orderedmap/v3"
)

func TestRaceCondition(t *testing.T) {
	m := orderedmap.NewSyncOrderedMap[int, int]()
	wg := &sync.WaitGroup{}

	var asyncGet = func() {
		wg.Add(1)
		go func() {
			key := rand.Intn(100)
			m.Get(key)
			wg.Done()
		}()
	}

	var asyncSet = func() {
		wg.Add(1)
		go func() {
			key := rand.Intn(100)
			value := rand.Intn(100)
			m.Set(key, value)
			wg.Done()
		}()
	}

	var asyncDelete = func() {
		wg.Add(1)
		go func() {
			key := rand.Intn(100)
			m.Delete(key)
			wg.Done()
		}()
	}

	var asyncHas = func() {
		wg.Add(1)
		go func() {
			key := rand.Intn(100)
			m.Has(key)
			wg.Done()
		}()
	}

	var asyncReplaceKEy = func() {
		wg.Add(1)
		go func() {
			key := rand.Intn(100)
			newKey := rand.Intn(100)
			m.ReplaceKey(key, newKey)
			wg.Done()
		}()
	}

	var asyncGetOrDefault = func() {
		wg.Add(1)
		go func() {
			key := rand.Intn(100)
			def := rand.Intn(100)
			m.GetOrDefault(key, def)
			wg.Done()
		}()
	}

	var asyncLen = func() {
		wg.Add(1)
		go func() {
			m.Len()
			wg.Done()
		}()
	}

	var asyncCopy = func() {
		wg.Add(1)
		go func() {
			m.Copy()
			wg.Done()
		}()
	}

	var asyncGetElement = func() {
		wg.Add(1)
		go func() {
			key := rand.Intn(100)
			e := m.GetElement(key)
			if e != nil {
				fmt.Println(e.Value)
			}
			wg.Done()
		}()
	}

	for i := 0; i < 10000; i++ {
		asyncSet()
		asyncGet()
		asyncDelete()
		asyncHas()
		asyncLen()
		asyncReplaceKEy()
		asyncGetOrDefault()
		asyncCopy()
		asyncGetElement()
	}

	wg.Wait()
	fmt.Println("TestRaceCondition completed")
	fmt.Printf("SyncOrderedMap eventually has %v elements\n", m.Len())
}
