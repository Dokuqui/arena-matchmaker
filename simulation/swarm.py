import asyncio
import grpc
import random
import time
import matchmaker_pb2
import matchmaker_pb2_grpc

# CONFIG
TOTAL_PLAYERS = 500
CONCURRENCY = 50
SERVER_ADDRESS = "localhost:50051"


async def run_bot(bot_id):
    """Simulates a single player joining and waiting."""
    try:
        async with grpc.aio.insecure_channel(SERVER_ADDRESS) as channel:
            stub = matchmaker_pb2_grpc.MatchmakerServiceStub(channel)

            my_mmr = random.randint(1000, 2000)

            start_time = time.time()

            response = await stub.FindMatch(
                matchmaker_pb2.FindMatchRequest(
                    player_id=f"Bot-{bot_id}", mmr=my_mmr, region="EU"
                )
            )

            # In a real game, we would now POLL the server until our ticket says "MATCH_FOUND".
            # Since our Go server creates the match internally but doesn't notify back yet,
            # we consider the "Request Accepted" as success for this load test.

            latency = (time.time() - start_time) * 1000
            return latency

    except Exception as e:
        print(f"❌ Bot-{bot_id} Failed: {e}")
        return None


async def main():
    print(f"🚀 Launching SWARM: {TOTAL_PLAYERS} bots targeting {SERVER_ADDRESS}")

    start_global = time.time()
    latencies = []

    for i in range(0, TOTAL_PLAYERS, CONCURRENCY):
        batch = []
        for j in range(CONCURRENCY):
            if i + j < TOTAL_PLAYERS:
                batch.append(run_bot(i + j))

        results = await asyncio.gather(*batch)
        for res in results:
            if res is not None:
                latencies.append(res)

        print(f"   Batch {i // CONCURRENCY + 1} sent...")
        await asyncio.sleep(0.1)

    duration = time.time() - start_global

    print("\n--- 📊 STRESS TEST RESULTS ---")
    print(f"Total Requests: {len(latencies)} / {TOTAL_PLAYERS}")
    print(f"Total Time:     {duration:.2f} seconds")
    print(f"Avg API Latency: {sum(latencies) / len(latencies):.2f} ms")
    print("-------------------------------")
    print("Check your Go Server logs. Is it churning through matches?")


if __name__ == "__main__":
    asyncio.run(main())
