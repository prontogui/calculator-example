module calculator

go 1.26

require github.com/prontogui/golib/v2 v2.0.0-00010101000000-000000000000

// 2.0 development: golib v2 has no release tag yet, so point at the local
// checkout instead of a version. Switch to a tagged version
// (github.com/prontogui/golib/v2 v2.x.y) once one is published.
replace github.com/prontogui/golib/v2 => ../golib

require (
	github.com/fxamacker/cbor/v2 v2.9.2 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260803160001-6ac0973c030d // indirect
	google.golang.org/grpc v1.83.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
