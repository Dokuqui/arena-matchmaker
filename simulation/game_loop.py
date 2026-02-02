import grpc
import time
import matchmaker_pb2
import matchmaker_pb2_grpc

def run_party_loop():
    channel = grpc.insecure_channel('localhost:50051')
    stub = matchmaker_pb2_grpc.MatchmakerServiceStub(channel)

    print("🎮 Queueing Team A (2 Solo Players)...")
    stub.FindMatch(matchmaker_pb2.FindMatchRequest(
        player_ids=["Warrior_01"], mmr=1200, region="EU"
    ))
    stub.FindMatch(matchmaker_pb2.FindMatchRequest(
        player_ids=["Mage_88"], mmr=1210, region="EU"
    ))

    print("🎮 Queueing Team B (Party of 2)...")
    stub.FindMatch(matchmaker_pb2.FindMatchRequest(
        player_ids=["Rogue_1", "Rogue_2"], mmr=1500, region="EU"
    ))

    print("⏳ Waiting for Matchmaker to group them (4 players needed)...")
    time.sleep(3)

    print("🏆 Match Over! Reporting Results...")
    stub.ReportResult(matchmaker_pb2.ReportResultRequest(
        match_id="match_123",
        winner_ids=["Warrior_01", "Mage_88"],
        loser_ids=["Rogue_1", "Rogue_2"]
    ))
    print("✅ Result Reported.")

if __name__ == '__main__':
    run_party_loop()
    