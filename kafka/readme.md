//查看consumer group列表，使用--list参数
kafka-consumer-groups.sh --bootstrap-server localhost:9092 --list

//查看特定consumer group 详情，使用--group与--describe参数
kafka-consumer-groups.sh  --bootstrap-server localhost:9092 --group grp1 --describe