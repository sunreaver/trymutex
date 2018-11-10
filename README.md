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