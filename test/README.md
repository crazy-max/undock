# Integration tests

Run from the repository root with Docker Compose:

```console
docker compose -f test/compose.yaml run --build --rm tests
```

Coverage is written to `bin/coverage/integration/coverage.txt`. Set
`TEST_PUBLIC_REGISTRY=1` to include the public-registry tests.

Clean up afterward:

```console
docker compose -f test/compose.yaml down --volumes --remove-orphans
```
