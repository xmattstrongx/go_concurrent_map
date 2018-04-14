| Benchmarks With The Lock Spaghetti          | Iterations    | Speed/Iteration  | Bytes alloc   | Allocs        |
| ------------------------------------------- | ------------- | ---------------- | ------------- | ------------- |
| BenchmarkPurgeWithLockSpaghetti1-8 |         	       1 |	1001295840 ns/op |	    3232 B/op |	      34 allocs/op |
| BenchmarkPurgeWithLockSpaghetti10-8 |        	       1 |	1001318668 ns/op |	    4656 B/op |	      36 allocs/op |
| BenchmarkPurgeWithLockSpaghetti100-8 |       	       1 |	1001301646 ns/op |	   24992 B/op |	      40 allocs/op |
| BenchmarkPurgeWithLockSpaghetti1000-8 |      	       1 |	1001327524 ns/op |	  358304 B/op |	      84 allocs/op |
| BenchmarkPurgeWithLockSpaghetti10000-8 |     	       1 |	1007438657 ns/op |	23220096 B/op |	    4082 allocs/op |
| BenchmarkPurgeWithLockSpaghetti1000000-8 |   	       1 |	1960973985 ns/op |	197053344 B/op |	   50263 allocs/op |


| Benchmarks With The Extra Delete          | Iterations    | Speed/Iteration  | Bytes alloc   | Allocs        |
| ------------------------------------------- | ------------- | ---------------- | ------------- | ------------- |
| BenchmarkPurgeWithExtraDelete1-8 |           	       1 |	1000161485 ns/op |	    2208 B/op |	      28 allocs/op |
| BenchmarkPurgeWithExtraDelete100-8 |         	       1 |	1000434512 ns/op |	    4160 B/op |	      31 allocs/op |
| BenchmarkPurgeWithExtraDelete1000-8 |        	       1 |	1000139811 ns/op |	   31296 B/op |	      49 allocs/op |
| BenchmarkPurgeWithExtraDelete10000-8 |       	       1 |	1001275805 ns/op |	  255200 B/op |	     117 allocs/op |
| BenchmarkPurgeWithExtraDelete100000-8 |      	       1 |	1001873836 ns/op |	 3450672 B/op |	     523 allocs/op |
| BenchmarkPurgeWithExtraDelete1000000-8 |     	       1 |	1006071585 ns/op |	28858672 B/op |	    8036 allocs/op |


go test -bench=. -benchmem -cpuprofile=cpu.out -memprofile=mem.out .
goos: darwin
goarch: amd64
pkg: github.com/xmattstrongx/go_concurrent_map

BenchmarkPurgeWithLockSpaghetti1-8         	       1	1001295840 ns/op	    3232 B/op	      34 allocs/op
BenchmarkPurgeWithLockSpaghetti10-8        	       1	1001318668 ns/op	    4656 B/op	      36 allocs/op
BenchmarkPurgeWithLockSpaghetti100-8       	       1	1001301646 ns/op	   24992 B/op	      40 allocs/op
BenchmarkPurgeWithLockSpaghetti1000-8      	       1	1001327524 ns/op	  358304 B/op	      84 allocs/op
BenchmarkPurgeWithLockSpaghetti10000-8     	       1	1007438657 ns/op	23220096 B/op	    4082 allocs/op
BenchmarkPurgeWithLockSpaghetti1000000-8   	       1	1960973985 ns/op	197053344 B/op	   50263 allocs/op
BenchmarkPurgeWithExtraDelete1-8           	       1	1000161485 ns/op	    2208 B/op	      28 allocs/op
BenchmarkPurgeWithExtraDelete100-8         	       1	1000434512 ns/op	    4160 B/op	      31 allocs/op
BenchmarkPurgeWithExtraDelete1000-8        	       1	1000139811 ns/op	   31296 B/op	      49 allocs/op
BenchmarkPurgeWithExtraDelete10000-8       	       1	1001275805 ns/op	  255200 B/op	     117 allocs/op
BenchmarkPurgeWithExtraDelete100000-8      	       1	1001873836 ns/op	 3450672 B/op	     523 allocs/op
BenchmarkPurgeWithExtraDelete1000000-8     	       1	1006071585 ns/op	28858672 B/op	    8036 allocs/op