import grpc
import matchmaker_pb2
import matchmaker_pb2_grpc
import time

def run():
    channel = grpc.insecure_channel('localhost:50051')
    stub = matchmaker_pb2_grpc.MatchmakerServiceStub(channel)

    print("🎮 Sending Player 1 (MMR 1200)...")
    stub.FindMatch(matchmaker_pb2.FindMatchRequest(
        player_id="Warrior_01", mmr=1200, region="EU"
    ))

    time.sleep(0.5)

    print("🎮 Sending Player 2 (MMR 1210)...")
    stub.FindMatch(matchmaker_pb2.FindMatchRequest(
        player_id="Mage_88", mmr=1210, region="EU"
    ))

    print("Requests Sent.")

if __name__ == '__main__':
    run()
