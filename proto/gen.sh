protoc -I. --go_out=plugins=grpc:../entity/entity_pb *.proto
protoc --grpc-gateway_out=logtostderr=true:../entity/entity_pb ./alliance_bot.proto
protoc  --swagger_out=logtostderr=true:../entity/entity_pb ./alliance_bot.proto