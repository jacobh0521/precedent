# PRECEDENT 📑

> A daily cosmic-customs logic game where players smuggle absurd cargo past a mutating rulebook.

---

## 🌌 Core Mechanics

1. **Absurd Manifest:** Each daily puzzle presents an absurd parcel (e.g., *"A Live Giraffe"*, *"Tomorrow's Weather"*, or *"Three Regrets"*) bearing distinct traits (`living`, `tall`, `organic`, `intangible`, `volatile`, `emotional`).
2. **Active Rulebook:** 5 deterministic statutes ban or restrict specific combinations of traits, containers, adjectives, or classifications—with deliberate loopholes.
3. **Structured Declaration:** The player selects:
   - **Classification:** `furniture`, `souvenir`, `diplomatic_pouch`, `spare_part`, `art_installation`, `agricultural_sample`
   - **Modifiers (up to 2):** `folded`, `ceremonial`, `dehydrated`, `rental`, `unopened`, `metaphorical`
   - **Container:** `lead_crate`, `mesh_bag`, `glass_jar`, `velvet_envelope`, `flatpack_box`, `thermos`
4. **Bureaucratic Audit:** Form submission checks the declaration against every active law and stamps a verdict (`CLEARED` or `CONFISCATED`) with detailed bureaucratic citations.
5. **The Twist (Law Mutation):** Whenever a declaration clears all active laws, the customs engine records the winning loophole and codifies an **Emergency Precedent Statute** (e.g., `AMEND-006-E`) into the rulebook for the next day. Players can click **"Simulate Next Day"** to immediately test subsequent days against the laws they spawned.

---

## 🛠 Tech Stack

- **Backend:** Go 1.24+ (`net/http` + `github.com/go-chi/chi/v5`)
- **Database:** Pure-Go SQLite (`modernc.org/sqlite` — zero CGO/gcc required)
- **Frontend:** Go `html/template` + **Tailwind CSS** (CDN) + **HTMX** (CDN) (zero Node build step)
- **Deployment:** Multi-stage Dockerfile targeting Google Cloud Run (`precedent-515201`)

---

## 🚀 Running Locally

### Prerequisites
- Go 1.22+ installed

### Start the Server
```bash
# Run server directly
go run ./cmd/server
```
The game will be available at [http://localhost:8080](http://localhost:8080).

### Run Test Suite
```bash
go test -v ./...
```

### Run with Docker / Deploy to Cloud Run
```bash
# Build container image
docker build -t precedent:latest .

# Run locally on port 8080
docker run -p 8080:8080 precedent:latest

# Deploy to Cloud Run
gcloud run deploy precedent-515201 \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```
