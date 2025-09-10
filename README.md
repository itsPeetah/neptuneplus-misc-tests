# "prime-numbers caller"

This is just miscellaneous code that I used to test my work on my thesis.

## Notes

- the profiling images are used with a separate profiling repo
- can still be used outside of a k8s cluster by overwriting the env vars in the container
- image for prime-numbers: `itspeetah/np-prime-numbers-go`

## Usage

- `/entrypoint?mode=<MODE>&count=<COUNT>&upperBound=<UPPER_BOUND>`
  - mode: string `seq` or `par` (should the requests be made sequentially or concurrently)
  - count: int >= 0 (numbers of calls to prime-numbers)
  - upperBound: int >= 0 (upper bound for prime numbers, e.g. /prime/<UPPER_BOUND>)
