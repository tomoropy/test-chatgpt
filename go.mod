module github.com/tomoropy/test-chatgpt

go 1.19

require github.com/r3labs/sse/v2 v2.0.0-00010101000000-000000000000

require (
	golang.org/x/net v0.0.0-20191116160921-f9c825593386 // indirect
	gopkg.in/cenkalti/backoff.v1 v1.1.0 // indirect
)

replace github.com/r3labs/sse/v2 => ../../ghq/github.com/munisystem/sse
