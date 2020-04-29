# Bloom filter

```go
b := New()
b.AddString("qwer")
b.AddInt(9)
println(b.HasString("qwer")) // print: true
b.HasInt(9) // print: true
```