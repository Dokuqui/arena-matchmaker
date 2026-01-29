import grpc
import time
import random
import matchmaker_pb2
import matchmaker_pb2_grpc


def run_game_loop():
    channel = grpc.insecure_channel("localhost:50051")
    stub = matchmaker_pb2_grpc.MatchmakerServiceStub(channel)

    p1 = "Warrior_01"
    p2 = "Mage_88"

    print(f"🎮 {p1} and {p2} queuing...")
    stub.FindMatch(matchmaker_pb2.FindMatchRequest(player_id=p1, region="EU"))
    stub.FindMatch(matchmaker_pb2.FindMatchRequest(player_id=p2, region="EU"))

    print("⚔️  Match Found! Playing game... (Simulating 2s)")
    time.sleep(2)

    winner, loser = (p1, p2) if random.random() > 0.5 else (p2, p1)

    print(f"🏆 Game Over! Winner: {winner}")

    resp = stub.ReportResult(
        matchmaker_pb2.ReportResultRequest(
            match_id="match_123", winner_id=winner, loser_id=loser
        )
    )

    print("📈 MMR Update:")
    print(f"   {winner}: {resp.winner_new_mmr} (+)")
    print(f"   {loser}: {resp.loser_new_mmr} (-)")


if __name__ == "__main__":
    run_game_loop()
