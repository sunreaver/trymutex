## 带有Try功能的sync Mutex

### 示例

```
func main() {
    m := NewTryMutexWithSyncMutex(&sync.Mutex{})
    m.Lock()      // 获取锁
    m.TryLock()   // 将不会等待锁，而是直接返回false
    m.TryUnLock() // 将会unlock，且返回true
    m.TryUnLock() // 不会panic，而是返回false
    m.IsLock()    // false
}
```

### 注意

`TryMutex` 的 `TryUnLock` 方法只会尝试解锁 TryMutex 实例自己锁定的锁

例如以下情况

```
func main() {
    mutex := &sync.Mutex{}
    m1 := NewTryMutexWithSyncMutex(mutex)
    m2 := NewTryMutexWithSyncMutex(mutex)
    m1.Lock()      // 获取锁
    m2.TryLock()   // 将不会等待锁，而是直接返回false
    m2.TryUnLock() // 将不会unlock，且返回false
    m1.TryUnLock() // 将会unlock，且返回true
}
```
