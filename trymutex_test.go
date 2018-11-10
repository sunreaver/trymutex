package trymutex

import (
	"sync"
	"sync/atomic"
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

			for i := 0; i < 100; i++ {
				var unlockCount int32
				var lockCount int32
				m.TryUnLock()
				c.So(m.IsLocked(), ShouldBeFalse)
				for index := 0; index < 10000; index++ {
					wg.Add(1)
					go func() {
						if m.TryUnLock() {
							atomic.AddInt32(&unlockCount, 1)
						}
						wg.Done()
					}()
					wg.Add(1)
					go func() {
						if m.TryLock() {
							atomic.AddInt32(&lockCount, 1)
						}
						wg.Done()
					}()
				}
				wg.Wait()
				if m.TryUnLock() {
					unlockCount++
				}
				// 起始状态是unlock
				// 结束状态是unlock
				c.So(unlockCount, ShouldEqual, lockCount)
				c.So(m.IsLocked(), ShouldBeFalse)
			}
		})
	})
}
