# Arena Matchmaker

A high-performance, concurrent matchmaking service built with **Go**, **gRPC**, and **Redis**.

Designed to handle high-throughput queuing using Redis Sorted Sets and Go concurrency patterns.

## Architecture

* **Core:** Golang (gRPC Server)
* **Database:** Redis (Sorted Sets for ELO queuing)
* **Protocol:** Protobufs (Strict schemas)
* **Infrastructure:** Docker & Docker Compose

## How to Run

1. **Start the Stack:**
   ```bash
   docker-compose up --build
   ```

2. **Run Stress Test:**
   ```bash
   python3 simulation/swarm.py
   ```

## Performance

* **Latency:** < 10ms per request under load.
* **Algorithm:** Skill-based matching (MMR) with time-based expansion logic.
