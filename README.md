# tiny-cache-go 💽

**This package is designed for quick embedding and easy use of Go application cache**

Example:

```
c := cache.New() // New cache storage
c.Set("key1", "Hi", 5*time.Minute)
c.Set("key2", []byte("bye!"), 1*time.Minute)

fmt.Println(c.Get("key1"))
fmt.Println(c.Get("key2"))

var s string = c.Get("key1").String()
var b, err = c.Get("key2").Byte()

var s string = c.GetD("key1").String() // Get key value and delete in memory
```
