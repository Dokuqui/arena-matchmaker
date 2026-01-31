# Arena Matchmaker

A high-performance, concurrent matchmaking service built with **Go**, **gRPC**, and **Redis**.

Designed to handle high-throughput queuing using Redis Sorted Sets and Go concurrency patterns.

## 🏗 Architecture

* **Core:** Golang (gRPC Server)
* **Database:** Redis
  * *Sorted Sets:* For sub-millisecond matchmaking queues.
  * *Hashes:* For persistent player stats (Elo/MMR).
* **Protocol:** Protobufs (Strict schemas)
* **Infrastructure:** Docker & Docker Compose (Microservices)

## How to Run

1. **Start the Stack:**
   ```bash
   docker-compose up --build
   ```

2. **Run Stress Test:**
   ```bash
   python3 simulation/swarm.py
   ```

3. **Run Full Game Loop (Match + Elo Update)**
   ```bash
   python3 simulation/game_loop.py
   ```

## Performance

* **Latency:** < 10ms per request under load.
* **Algorithm:** Skill-based matching (MMR) with time-based expansion logic.

## 📊 Observability

The system includes a built-in Prometheus metrics exporter to track real-time performance.

* **Metrics Port:** `:2112`
* **Key Metrics:**
  * `arena_queue_depth`: Real-time count of players waiting.
  * `arena_matches_total`: Cumulative counter of matches formed.
  * `arena_match_latency_seconds`: Histogram of wait times.

![Metrics Dashboard](docs/metrics_demo.png)

* **Ranked System:** Implements a full **Elo Rating System** (K-Factor 32).

  * Tracks player progression persistently across matches.
  * Handles post-match reporting and rating updates via gRPC.

## 🏗 Architecture (Microservices)

The system is split into independent services for scalability:
1.  **Frontend Service:** Handles high-throughput gRPC connections, validation, and pushing to Redis. Scalable to N replicas.
2.  **Matchmaker Worker:** Singleton background worker that processes the Redis queue and executes the matchmaking logic (Elo).

