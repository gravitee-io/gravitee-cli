# gio — Gravitee CLI
Labs repo for building a production-ready Gravitee CLI. We iterate here until the tool is solid, tested, and ready to ship.

Currently covers API Management (APIM). Access Management (AM) will be integrated next.

For now this is a straight 1:1 wrapper around the management APIs. More tooling (GitOps workflows, smart resolution, etc.) will come later.

```bash
make build    # build
make test     # tests
make lint     # lint
```
