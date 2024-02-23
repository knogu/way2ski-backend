package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	pb "way2ski-backend/proto"
)

const (
	ECHIGOYUZAWA = "越後湯沢"
	TOKYO        = "東京"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedWayServiceServer
}

type Leg struct {
	LineName         string
	DepartureStation string
	DepartureHour    int
	DepartureMinute  int
	ArrivalStation   string
	ArrivalHour      int
	ArrivalMinute    int
	IsHoliday        int
}

func convLegToRespType(leg Leg) pb.OneCandidateInLeg {
	return pb.OneCandidateInLeg{
		LineName:         leg.LineName,
		DepartureStation: leg.DepartureStation,
		DepartureHour:    uint32(leg.DepartureHour),
		DepartureMinute:  uint32(leg.DepartureMinute),
		ArrivalStation:   leg.ArrivalStation,
		ArrivalHour:      uint32(leg.ArrivalHour),
		ArrivalMinute:    uint32(leg.ArrivalMinute),
	}
}

func getLegsFromDb(departureStation string, arrivalStation string, isHoliday bool) pb.CandidatesInOneLeg {
	rows, err := db.Query(`
		SELECT * FROM legs 
		WHERE departure_station = ? 
		AND arrival_station = ? 
		AND is_holiday = ?
	`, departureStation, arrivalStation, isHoliday)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var legsRespType []*pb.OneCandidateInLeg
	for rows.Next() {
		var leg Leg
		err = rows.Scan(&leg.LineName, &leg.DepartureStation, &leg.DepartureHour, &leg.DepartureMinute, &leg.ArrivalStation, &leg.ArrivalHour, &leg.ArrivalMinute, &leg.IsHoliday)
		if err != nil {
			log.Fatal(err)
		}
		v := convLegToRespType(leg)
		legsRespType = append(legsRespType, &v)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return pb.CandidatesInOneLeg{
		CandidatesInOneLeg: legsRespType,
	}
}

func getLegs(curStation string, arrivalStation string, cur2nextStation map[string]string, isHoliday bool) pb.AllLegsInOneWay {
	println("GetLegs called")
	var ret []*pb.CandidatesInOneLeg
	for {
		println(curStation)
		if curStation == arrivalStation {
			break
		}
		nextStation, ok := cur2nextStation[curStation]
		if !ok {
			panic("err")
		}
		v := getLegsFromDb(curStation, nextStation, isHoliday)
		ret = append(ret, &v)
		curStation = nextStation
	}
	println("GetLegs returning...")
	return pb.AllLegsInOneWay{
		AllLegsInOneWay: ret,
	}
}

func getLegsHome(in *pb.Params) *pb.AllLegsInOneWay {
	println("GetLegsHome called")
	cur2nextStationHome := map[string]string{
		in.SkiResort: ECHIGOYUZAWA,
		ECHIGOYUZAWA: TOKYO,
		TOKYO:        in.HometownStation,
	}
	v := getLegs(in.SkiResort, in.HometownStation, cur2nextStationHome, in.IsHoliday)
	return &v
}

func getLegsToSki(in *pb.Params) *pb.AllLegsInOneWay {
	println("GetLegToSki called")
	cur2nextStationToSki := map[string]string{
		in.HometownStation: TOKYO,
		TOKYO:              ECHIGOYUZAWA,
		ECHIGOYUZAWA:       in.SkiResort,
	}
	v := getLegs(in.HometownStation, in.SkiResort, cur2nextStationToSki, in.IsHoliday)
	return &v
}

func (s *server) GetLines(ctx context.Context, in *pb.Params) (*pb.Lines, error) {
	println("GetLines called")
	in.HometownStation = "四ツ谷"
	in.SkiResort = "かぐらスキー場"
	in.IsHoliday = true
	return &pb.Lines{
		AllLegsToSki: getLegsToSki(in),
		AllLegsHome:  getLegsHome(in),
	}, nil
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./train_bus_time/train_bus_time.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterWayServiceServer(s, &server{})
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
