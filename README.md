[![GoDoc](https://godoc.org/github.com/go-ng/xsync?status.svg)](https://pkg.go.dev/github.com/go-ng/xsync?tab=doc)


This package provides `BunchAllocator` which is close to `sync.Pool` by use cases, but it does not reuse memory to gain performance. Instead it allocates data by bunches (avoiding a lot of small allocated pieces). It is much less effective than `sync.Pool`, but does not require to put data back to the pool.

```
name                               old time/op    new time/op    delta
BunchAllocator/structSize1-16        7.62ns ± 5%    7.34ns ± 4%      ~     (p=0.310 n=5+5)
BunchAllocator/structSize8-16        10.5ns ± 2%     8.1ns ±15%   -22.84%  (p=0.008 n=5+5)
BunchAllocator/structSize32-16       17.0ns ± 7%    11.9ns ± 5%   -29.60%  (p=0.008 n=5+5)
BunchAllocator/structSize256-16      44.7ns ± 3%    33.1ns ±14%   -26.06%  (p=0.008 n=5+5)
BunchAllocator/structSize1024-16      156ns ±11%     190ns ± 5%   +21.78%  (p=0.008 n=5+5)
BunchAllocator/structSize65536-16    4.35µs ±27%    5.41µs ± 5%      ~     (p=0.095 n=5+5)

name                               old alloc/op   new alloc/op   delta
BunchAllocator/structSize8-16         8.00B ± 0%     9.00B ± 0%   +12.50%  (p=0.008 n=5+5)
BunchAllocator/structSize32-16        32.0B ± 0%     37.0B ± 0%   +15.62%  (p=0.008 n=5+5)
BunchAllocator/structSize256-16        256B ± 0%      288B ± 0%   +12.50%  (p=0.008 n=5+5)
BunchAllocator/structSize1024-16     1.02kB ± 0%    1.06kB ± 0%    +3.09%  (p=0.008 n=5+5)
BunchAllocator/structSize65536-16    65.5kB ± 0%    65.6kB ± 0%    +0.05%  (p=0.008 n=5+5)
```