elzzes-calculator-api has docs only. The team needs a common base before spec 001: a Go module, a net/http adapter, structured logging, config loading, and a working GET /healthz. No calculate route or arithmetic.

Instead of making a api with no ports at all, lets structure this as a lightweight Ports & Adapters layout so boundaries are replaceable (net/http → Fiber, slog → Zap, for example) without touching the calculator core.

Important: ports are defined by the layer that needs the capability, not by the infrastructure that provides it. So dont create a generic commons package. The desired structure would be somehow similar to this:

api/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── calculator/
│   │   ├── calculator.go
│   │   └── errors.go
│   │
│   ├── transport/
│   │   └── http/
│   │       ├── handler.go
│   │       ├── routes.go
│   │       ├── middleware/
│   │       │   └── logging.go
│   │        └── logger.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   └── logger/
│       └── slog.go


> After this, the llm wrote the code but I had to make some adjustments:
The panic-to-JSON-500 path in NewHandler is untested. The closure in handler.go that writes INTERNAL_ERROR can’t be exercised, because NewHandler hard-wires its routes and there’s no way to inject a panicking one. The middleware tests cover Recover in isolation, but not the real wiring. Consider splitting out something like wrap(router, logger) so a test can add a panicking route and assert the 500 envelope.

TestErrorCodeValues is tautological (errors_test.go). It builds a map from the constants and compares each to a hardcoded string, so it only fails if someone edits the test and the constant differently. The len(want) != 13 check counts the test’s own map. Considering that openapi.yaml is canonical, a useful version would parse that file and compare enum values.