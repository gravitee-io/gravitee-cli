<!--
Copyright (C) 2015 The Gravitee team (http://gravitee.io)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Secure an API with OAuth2 / JWT, end to end

A single script that wires the full AM + APIM chain from the command line: an AM domain, an OAuth2 client, an API imported from JSON, a JWT plan validating tokens against AM's JWKS, a subscription, deploy, a live test, and observability. Zero UI.

This is a working example of **how `gctl` is built for scripting**: one CLI for both AM and APIM, composable commands that pass ids from one step to the next, and JSON resources imported from version control. The same script runs unchanged in CI. The patterns it uses are summarized in [the table at the end](#scripting-patterns-worth-reusing).

> **The URLs in this cookbook are fictitious.** `am.gateway.gravitee.io` and `apim.gateway.gravitee.io` stand in for your AM and APIM gateway addresses; replace them with your own.

## Prerequisites

- A running Gravitee stack (AM + APIM), with both gateways reachable.
- `gctl` installed and configured for both products (see the [main README](../README.md)).
- `jq` and `envsubst` available (standard on most systems; `envsubst` ships with GNU gettext).

## The trust chain

```
Client  ->  AM    : "here's my client_id + secret, give me a token"
AM      ->  Client: signed JWT (RS256, AM private key)
Client  ->  APIM  : "here's my token"
APIM    ->  AM JWKS: "is this signature valid?" (background)
APIM    ->  Client: 200 or 401
```

The gateway does two checks: the JWT signature (against AM's JWKS), and that the token's `client_id` has an **active subscription** to the plan. The second check is what lets you revoke access without invalidating the token: just close the subscription.

## The config files
`api.json` defines a v4 HTTP proxy API targeting an echo backend (it echoes headers, handy to confirm the token reaches the upstream):

```json
{
  "name": "orders-api",
  "apiVersion": "1.0.0",
  "definitionVersion": "V4",
  "type": "PROXY",
  "description": "Demo API secured with OAuth2/JWT via Gravitee AM",
  "listeners": [
    {
      "type": "HTTP",
      "paths": [{ "path": "/orders" }],
      "entrypoints": [
        { "type": "http-proxy", "qos": "AUTO", "configuration": {} }
      ]
    }
  ],
  "endpointGroups": [
    {
      "name": "default-group",
      "type": "http-proxy",
      "endpoints": [
        {
          "name": "default",
          "type": "http-proxy",
          "inheritConfiguration": false,
          "configuration": { "target": "https://api.gravitee.io/echo" },
          "secondary": false
        }
      ]
    }
  ],
  "flowExecution": { "mode": "DEFAULT", "matchRequired": false },
  "flows": [],
  "analytics": { "enabled": true }
}
```

`plan-jwt.json` is the security policy: the APIM gateway validates incoming JWTs against AM's JWKS endpoint. `validation: AUTO` skips manual subscription approval.

```json
{
  "definitionVersion": "V4",
  "name": "JWT via AM",
  "description": "JWT plan. Validates tokens issued by the AM domain demo-orders against its JWKS endpoint.",
  "mode": "STANDARD",
  "validation": "AUTO",
  "security": {
    "type": "JWT",
    "configuration": {
      "signature": "RSA_RS256",
      "publicKeyResolver": "JWKS_URL",
      "resolverParameter": "https://am.gateway.gravitee.io/demo-orders/oidc/.well-known/jwks.json",
      "useSystemProxy": false,
      "extractClaims": true,
      "propagateAuthHeader": true,
      "jwksCacheDurationSeconds": 0
    }
  },
  "selectionRule": null,
  "flows": []
}
```

> The JWKS URL points to your AM gateway: the APIM gateway calls it server-to-server to fetch AM's public keys and verify the token signature.

`app.json` is the APIM application. Its `client_id` is a `${CLIENT_ID}` placeholder, filled in at import time so the real value never lands in Git:

```json
{
  "name": "orders-api-client",
  "description": "APIM application linked to the AM OAuth2 client",
  "settings": {
    "app": {
      "client_id": "${CLIENT_ID}"
    }
  }
}
```

## The script, step by step

### Step 1: AM domain

An AM **domain** is an isolated identity tenant with its own OAuth2 clients, scopes, and signing certificates. A fresh domain is inactive by default, so `enable` is required before the AM gateway will serve it.

```bash
export DOMAIN_ID=$(gctl am domain create --name "demo-orders" -o id)
gctl am domain enable "$DOMAIN_ID"
```

Note the `-o id` pattern: the command prints only the new domain's id, captured straight into a variable.

### Step 2: AM application (machine-to-machine client)

`--type service` maps to the `client_credentials` grant: no user, two machines trusting each other via a shared secret. Three calls: create to get the app id, `get` to read the `clientId`, `secret create` to generate a named secret (its value is only returned at creation time).

```bash
export APP_ID=$(gctl am app create \
  --domain "$DOMAIN_ID" \
  --name "orders-api-client" \
  --type service \
  -o id)
export CLIENT_ID=$(gctl am app get "$APP_ID" --domain "$DOMAIN_ID" -o json | jq -r .settings.oauth.clientId)
export CLIENT_SECRET=$(gctl am app secret create --domain "$DOMAIN_ID" --app-id "$APP_ID" --name "demo-secret" -o json | jq -r .secret)
```

### Step 3: import the API into APIM

The API is imported from `api.json` over stdin. Created in `STOPPED` state with no plan yet.

```bash
export API_ID=$(gctl apim api create -o id < api.json)
```

### Step 4: JWT plan

Create the plan from `plan-jwt.json` and publish it so it accepts subscriptions.

```bash
export PLAN_ID=$(gctl apim plan create --api "$API_ID" -o id < plan-jwt.json)
gctl apim plan publish "$PLAN_ID" --api "$API_ID" -q
```

### Step 5: link the two products (APIM application + subscription)

Why an app in APIM when one already exists in AM? Because AM manages **identity** and APIM manages **access**. The APIM app carries the subscription, and setting its `client_id` to the AM client's `client_id` is what binds the two.

This is where the placeholder gets filled. `envsubst` substitutes `${CLIENT_ID}` (captured in step 2) into the JSON before it reaches `gctl`, so the committed file stays free of environment-specific values:

```bash
export APIM_APP_ID=$(envsubst < app.json | gctl apim application create -o id)
gctl apim subscription create --api "$API_ID" --plan "$PLAN_ID" --app "$APIM_APP_ID"
```

> The same `envsubst` trick injects any secret or per-environment value (client secret, JWKS URL, backend target) at import time. For overriding non-placeholder fields, pipe through `jq` instead: `jq '.field = "value"' api.json | gctl apim api create`.

### Step 6: deploy and start

`deploy` pushes the config to the gateway (without it, changes stay in the Management API). `start` opens traffic; an API can be deployed but stopped for maintenance.

```bash
gctl apim api deploy "$API_ID" --label "demo-oauth2"
gctl apim api start "$API_ID"
sleep 3
```

### Step 7: test

Get a token from AM via `client_credentials`, then call the API. The gateway verifies the RS256 signature against AM's JWKS, checks the `client_id` has an active subscription, and lets it through.

```bash
TOKEN=$(curl -s -u "$CLIENT_ID:$CLIENT_SECRET" \
  -d "grant_type=client_credentials" \
  https://am.gateway.gravitee.io/demo-orders/oauth/token | jq -r .access_token)

# With a valid token -> 200
curl -H "Authorization: Bearer $TOKEN" https://apim.gateway.gravitee.io/orders

# Without a token -> 401
curl -i https://apim.gateway.gravitee.io/orders
```

### Step 8: cross-product observability

Query analytics from **both products in the same session**: APIM request logs and AM's token-issuance audit trail. From token issuance to API call, everything is observable without switching tools.

```bash
gctl apim log list --api /orders
gctl am audit list --domain "$DOMAIN_ID" --type TOKEN_CREATED
```

> Resource flags accept either the id or the API path, so `$API_ID` and `/orders` are interchangeable. The path form is what lets the teardown below run in a fresh shell, with no captured variables.

### Teardown

Idempotent cleanup, scriptable like the rest:

```bash
gctl apim api delete /orders --force
gctl apim application list --query orders-api-client -o id | xargs -n 1 gctl apim application delete
gctl am domain delete "$(gctl am domain list --query "demo-orders" -o json | jq -r '.data[0].id')"
```

## Scripting patterns worth reusing

| Pattern | What it does |
|---------|--------------|
| `ID=$(gctl ... -o id)` | Capture the created resource's id for the next command |
| `gctl ... -o json \| jq -r .field` | Pull a single field out of a response |
| `gctl ... create < resource.json` | Import a version-controlled resource from stdin |
| `envsubst < resource.json \| gctl ... create` | Inject env vars / secrets at import time |
| `gctl ... list --query NAME -o id \| xargs -n 1 gctl ... delete` | Look up by name, then act on the result |
| `-q` / `--force` | Run non-interactively in CI |
