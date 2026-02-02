import grpc
import time
import matchmaker_pb2
import matchmaker_pb2_grpc

def run_region_test():
    channel = grpc.insecure_channel('localhost:50051')
    stub = matchmaker_pb2_grpc.MatchmakerServiceStub(channel)

    print("🌍 --- Starting Region Sharding Test ---")

    print("\n🇪🇺 [Client] Sending EU-WEST Requests...")
    stub.FindMatch(matchmaker_pb2.FindMatchRequest(
        player_ids=["EU_Knight", "EU_Priest"], 
        mmr=1500, 
        region="EU-WEST"
    ))
    stub.FindMatch(matchmaker_pb2.FindMatchRequest(
        player_ids=["EU_Mage", "EU_Archer"], 
        mmr=1520, 
        region="EU-WEST"
    ))

    print("\n🇺🇸 [Client] Sending US-EAST Requests...")
    stub.FindMatch(matchmaker_pb2.FindMatchRequest(
        player_ids=["US_Marine", "US_Sniper"], 
        mmr=1200, 
        region="US-EAST"
    ))
    stub.FindMatch(matchmaker_pb2.FindMatchRequest(
        player_ids=["US_Medic", "US_Pilot"], 
        mmr=1210, 
        region="US-EAST"
    ))

    print("\n⏳ Waiting for Region Workers to process...")
    time.sleep(2) 
    print("✅ Done. Check your Docker logs to see the separation!")

if __name__ == '__main__':
    run_region_test()
    