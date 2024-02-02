# tiny-cache-go 💽

**This package is designed for quick embedding and easy use of Go application cache**

This package also automatically removes all garbage from the cache after the clear timer expires

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

You can also configure the cache to clear and expire cache times

Example:

```
c := cache.New()
c.SetDefaultCleanupInterval(1*time.Hour)
c.SetDefaultExpiration(1*time.Hour)
```

![Caching and performance optimization in golang.](https://media.licdn.com/dms/image/D4D12AQEWMwIgrMa0eg/article-cover_image-shrink_600_2000/0/1685216253061?e=2147483647&v=beta&t=dK37ND_rXXfi1dbNckpVb5-Q1H96QP7WIwbn9O7nYIc)
