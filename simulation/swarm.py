import grpc
import uuid
import random
import asyncio
import matchmaker_pb2
import matchmaker_pb2_grpc

# Config
TOTAL_PLAYERS = 1000
CONCURRENCY = 20


async def send_request(stub, region):
    p_id = f"User_{uuid.uuid4().hex[:8]}"
    mmr = random.randint(1000, 2000)

    try:
        await stub.FindMatch(
            matchmaker_pb2.FindMatchRequest(player_ids=[p_id], mmr=mmr, region=region)
        )
        return True
    except grpc.RpcError as e:
        print(f"❌ Failed: {e.code()}")
        return False


async def worker(channel, region):
    stub = matchmaker_pb2_grpc.MatchmakerServiceStub(channel)
    while True:
        await send_request(stub, region)
        await asyncio.sleep(random.uniform(0.01, 0.1))


async def main():
    async with grpc.aio.insecure_channel("localhost:50051") as channel:
        print(f"🚀 Launching Swarm: {TOTAL_PLAYERS} players...")

        tasks = []
        for _ in range(10):
            tasks.append(asyncio.create_task(worker(channel, "EU-WEST")))

        for _ in range(10):
            tasks.append(asyncio.create_task(worker(channel, "US-EAST")))

        await asyncio.sleep(10)

        for t in tasks:
            t.cancel()

        print("🛑 Swarm finished.")


if __name__ == "__main__":
    asyncio.run(main())
