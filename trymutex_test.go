package trymutex

import (
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewTryMutexWithSyncMutex(t *testing.T) {
	Convey("test", t, func(c C) {
		Convey("new", func(c C) {
			m := NewTryMutexWithSyncMutex(&sync.Mutex{})
			c.So(m.IsLocked(), ShouldBeFalse)
			locked := m.TryLock()
			c.So(locked, ShouldBeTrue)
			locked = m.TryLock()
			c.So(locked, ShouldBeFalse)
			unlocked := m.TryUnLock()
			c.So(unlocked, ShouldBeTrue)
			c.So(m.IsLocked(), ShouldBeFalse)
			c.So(m.TryUnLock(), ShouldBeFalse)
			c.So(m.TryLock(), ShouldBeTrue)
			c.So(m.TryLock(), ShouldBeFalse)
			c.So(m.IsLocked(), ShouldBeTrue)

			wg := sync.WaitGroup{}
			for index := 0; index < 100; index++ {
				wg.Add(1)
				go func() {
					r := m.TryLock()
					c.So(r, ShouldBeFalse)
					wg.Done()
				}()
			}
			wg.Wait()

			var unlockCount int32
			for index := 0; index < 100; index++ {
				wg.Add(1)
				go func() {
					if m.TryUnLock() {
						unlockCount++
					}
					wg.Done()
				}()
			}
			c.So(unlockCount, ShouldEqual, 1)
			wg.Wait()

		})

		Convey("mutex", func(c C) {
			mutex := &sync.Mutex{}
			m1 := NewTryMutexWithSyncMutex(mutex)
			m2 := NewTryMutexWithSyncMutex(mutex)
			m1.Lock()
			c.So(m2.TryLock(), ShouldBeFalse)
			c.So(m2.TryUnLock(), ShouldBeFalse)
			c.So(m1.TryUnLock(), ShouldBeTrue)
		})
	})
}
