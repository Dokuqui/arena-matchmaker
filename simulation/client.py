import grpc
import matchmaker_pb2
import matchmaker_pb2_grpc
import time

def run():
    channel = grpc.insecure_channel('localhost:50051')
    stub = matchmaker_pb2_grpc.MatchmakerServiceStub(channel)

    print("Client: Connecting to Matchmaker...")

    response = stub.FindMatch(matchmaker_pb2.FindMatchRequest(
        player_id="Warrior_01",
        mmr=1200,
        region="EU"
    ))

    print(f"Server Response: Ticket ID = {response.ticket_id}")

if __name__ == '__main__':
    run()
