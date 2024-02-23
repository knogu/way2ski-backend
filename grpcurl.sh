grpcurl \
    -protoset <(buf build -o -) -plaintext \
    -d '{"hometown_station": "四ツ谷", "ski_resort": "かぐらスキー場", "is_holiday": true}' \
localhost:8080 way.v1.WayService/GetLines