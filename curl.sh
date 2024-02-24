curl \
    --header "Content-Type: application/json" \
--data '{"hometown_station": "東京", "ski_resort": "かぐらスキー場", "is_holiday": true}'  \
http://localhost:8080/way.v1.WayService/GetLines
