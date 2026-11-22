# Integration tests

Run from the repository root with Docker Compose:

```console
docker compose -f test/compose.yaml run --build --rm tests
```

Tests run through the vendored `go tool gotestsum`. Coverage, JSON, and JUnit
reports are written to `bin/testreports/integration`. Set `TEST_PUBLIC_REGISTRY=1`
to include the public-registry tests.

Clean up afterward:

```console
docker compose -f test/compose.yaml down --volumes --remove-orphans
```
