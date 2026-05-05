# AM CLI — Gap Analysis & Merge Plan

> **Stan na 2026-04-30.** Porównanie branch `am` vs `main`. Celem jest doprowadzenie `main` do pełnej parności z branchem `am`, zachowując jednocześnie wszystko co zostało dodane na `main` (factor, flow, group, audit, token, health, whoami, export, import, copy).

---

## Tło

- **Branch `main`** — ma kompletne CRUD dla: domain (+ export/import/copy), application, user, idp (tylko list/get), factor, flow, group, certificate (tylko list/get/delete), audit, token, health, whoami. Używa warstwy serwisowej (`f.AM()` → `internal/am/`).
- **Branch `am`** — niezależna implementacja z innymi wzorcami architektonicznymi. Ma dodatkowe komendy (login, set domain, domain update, idp create/update/delete, role update/delete, scope update/delete, cert create/update, app settings), ale nie ma: factor, flow, group, audit, token, health, whoami, export/import/copy.

**Uwaga architektoniczna:** Branch `am` woła API bezpośrednio przez `f.Client` + `cmdutil.AMEnvPath/AMDomainPath`, omijając warstwę serwisową. Branch `main` idzie przez `f.AM()`. Decyzja przy merge: zostajemy przy warstwie serwisowej z `main` — to jest wzorzec ustalony w repo.

---

## Różnice infrastrukturalne do adoptowania

Te zmiany z brancha `am` trzeba zaaplikować na `main` — nie są komendami, ale są wymagane przez nowe komendy.

### 1. `internal/config/config.go`

Branch `am` dodaje do struktury `Context`:
```go
Type   string `yaml:"type,omitempty"`   // "am" lub "apim"
Domain string `yaml:"domain,omitempty"` // aktywna domena AM
```

Do struktury `Config`:
```go
CurrentContext string // zamiast Current
```

Metoda `SaveTo(path string) error` — zapis configu do pliku.

**Wpływ:** Wymagane przez `gio am set domain` i `gio am login`.

### 2. `internal/cmdutil/cmdutil.go`

Branch `am` dodaje:
- `RequireAMContext(f)` — sprawdza czy kontekst jest typu AM
- `RequireAMDomain(f)` — sprawdza czy jest ustawiona aktywna domena
- `AMEnvPath(f, path)` — buduje ścieżkę `/management/organizations/{org}/environments/{env}/{path}`
- `AMDomainPath(f, path)` — buduje ścieżkę z aktywną domeną z `f.Resolved.Domain`
- `CheckReadOnly(f, cmdName)` — błąd jeśli kontekst jest read-only

### 3. `internal/client/client.go`

Branch `am` dodaje path helpery:
- `AMEnvPath(org, env, path) string`
- `AMDomainPath(org, env, domain, path) string`

---

## Brakujące komendy do implementacji na `main`

### Priorytet 1 — proste, niezależne

#### `gio am login`
- **Plik:** `cmd/am/login.go`
- **Co robi:** Uwierzytelnia przez AM Management API (`/management/auth/token`) lub przyjmuje token bezpośrednio. Zapisuje kontekst do configu z `Type: "am"`.
- **Flagi:** `--url` (required), `--token`, `--username`, `--password`, `--context`, `--org`, `--env-id`
- **Gotowe na `am`:** `cmd/am/login.go` — do przepisania w stylu `main` (przez serwis zamiast `f.Client`)

#### `gio am set domain`
- **Plik:** `cmd/am/set.go`
- **Co robi:** Ustawia aktywną domenę w configu (`ctx.Domain = domainID`). Bez tego `RequireAMDomain` nie przejdzie.
- **Podkomendy:** `set domain <id>`, `set domain --clear`
- **Wymaga:** `config.SaveTo()`, `Domain` field w `Context`
- **Gotowe na `am`:** `cmd/am/set.go`

#### `gio am domain update`
- **Plik:** `cmd/am/domain/update.go`
- **Co robi:** PUT na domenę z flagami `--name`, `--description` lub `--file` (JSON)
- **Wymaga:** `am.Service.UpdateDomain()` + `internal/am/domain.go`
- **Gotowe na `am`:** `cmd/am/domain/update.go` (używa `f.Client` bezpośrednio — przepisać przez serwis)

### Priorytet 2 — pełny CRUD dla zasobów które mamy tylko list/get

#### `gio am idp create/update/delete`
- **Pliki:** `cmd/am/idp/create.go`, `update.go`, `delete.go`
- **Co robi:** Pełny CRUD przez `--file` (JSON body)
- **Wymaga:** rozszerzenia `IDPService` o `CreateIDP`, `UpdateIDP`, `DeleteIDP`
- **Gotowe na `am`:** wszystkie trzy pliki

#### `gio am role update/delete`
- **Pliki:** `cmd/am/role/update.go`, `delete.go`
- **Wymaga:** `RoleService.UpdateRole()`, `DeleteRole()`
- **Gotowe na `am`:** oba pliki

#### `gio am scope update/delete`
- **Pliki:** `cmd/am/scope/update.go`, `delete.go`
- **Wymaga:** `ScopeService.UpdateScope()`, `DeleteScope()`
- **Gotowe na `am`:** oba pliki

#### `gio am certificate create/update`
- **Pliki:** `cmd/am/certificate/create.go`, `update.go`
- **Wymaga:** `CertificateService.CreateCertificate()`, `UpdateCertificate()`
- **Gotowe na `am`:** oba pliki

### Priorytet 3 — nowe funkcjonalności

#### `gio am application settings`
- **Plik:** `cmd/am/application/settings.go` (lub `app/settings.go`)
- **Co robi:** Widok lub aktualizacja OAuth2 settings aplikacji. Flagi: `--grant-types`, `--response-types`, `--redirect-uris`, `--post-logout-uris`, `--token-lifetime`, `--refresh-token-lifetime`, `--id-token-lifetime`
- **Gotowe na `am`:** `cmd/am/app/settings.go` — do adaptacji do pakietu `application`

---

## Komendy tylko na `main` — NIE MA ich na `am` (zachować)

| Komenda | Status |
|---|---|
| `gio am factor list/get` | ✅ tylko main, zachować |
| `gio am flow list/get` | ✅ tylko main, zachować |
| `gio am group list/get/create/delete` | ✅ tylko main, zachować |
| `gio am audit list/get` | ✅ tylko main, zachować |
| `gio am token list/create/revoke` | ✅ tylko main, zachować |
| `gio am health` | ✅ tylko main, zachować |
| `gio am whoami` | ✅ tylko main, zachować |
| `gio am domain export/import/copy` | ✅ tylko main, zachować |

---

## Kwestia nazewnictwa pakietu: `app` vs `application`

Branch `am` używa `cmd/am/app/` (krótszy Use: `"app"`), `main` ma `cmd/am/application/` (Use: `"application"`).

**Decyzja:** Zostajemy przy `application` (spójne z APIM pattern). `settings.go` z brancha `am` trafia do `cmd/am/application/settings.go`.

---

## Zadania do wykonania (ordered)

- [ ] **Task A:** Adaptacja `internal/config` — dodać `Type`, `Domain` do `Context`; `SaveTo()` do `Config`
- [ ] **Task B:** Adaptacja `internal/cmdutil` — dodać `RequireAMContext`, `RequireAMDomain`, `CheckReadOnly`, `AMDomainPath`
- [ ] **Task C:** `gio am login` — `cmd/am/login.go` (przepisać przez config.SaveTo)
- [ ] **Task D:** `gio am set domain` — `cmd/am/set.go` (wymaga Task A+B)
- [ ] **Task E:** `gio am domain update` — rozszerzyć `DomainService` + `cmd/am/domain/update.go`
- [ ] **Task F:** IDP pełny CRUD — rozszerzyć `IDPService` + `create.go`, `update.go`, `delete.go`
- [ ] **Task G:** Role update/delete — rozszerzyć `RoleService` + pliki
- [ ] **Task H:** Scope update/delete — rozszerzyć `ScopeService` + pliki
- [ ] **Task I:** Certificate create/update — rozszerzyć `CertificateService` + pliki
- [ ] **Task J:** `gio am application settings` — przenieść i zaadaptować z `am/app/settings.go`

---

## Co zostaje poza zakresem (nie w am-cli gio-cli)

Komendy z TypeScript am-cli których nie ma ani na `am` ani na `main` branchu gio-cli, i które wymagają osobnej decyzji:

| Komenda TS | Opis | Priorytet |
|---|---|---|
| `gio am logs` | Polling audit logów real-time (--follow) | Niski |
| `gio am watch` | Live dashboard w terminalu | Niski |
| `gio am diff` | Porównanie domen między kontekstami | Średni |
| `gio am lint` | Security audit konfiguracji domeny | Niski |
| `gio am trace` | Trace ścieżki auth | Niski |
| `gio am support-dump` | Diagnostic dump | Niski |
| `gio am doctor` | Diagnostyka połączenia | Niski |
| `gio am test-oidc` | Test OIDC flows (discover/login/cc) | Średni |
