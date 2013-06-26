cache
=====
A simple string cache with string values.
Will maybe be made generic in some way some day.

It has two implementations depending on your need but you probably want this:

```go
c := cache.NewConcurrentCache()
```

This cache multiplexes the keys over 16 normal caches depending on the hash value of the key.
There are a number of uninvestigated issues and some known such as the inability to handle
empty keys. Not sure that it should but it should at least handle it in the Set method in that case.
Hash function is blatantly stolen from Java's HashMap and its use with the keys are naive at best.

It scales significantly better than a normal type cache since the buckets are not accessed using any locks.

Run benchmarks as:

```shell
go test -v -bench=".*"
```

If you don't see at least 4-5 times improvement in higher levels of parallellism then something is wrong.