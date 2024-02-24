package main

import (
	"context"
	"database/sql"
	"github.com/bufbuild/connect-go"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	way "way2ski-backend/gen/way/v1"
	"way2ski-backend/gen/way/v1/wayv1connect"
)

const (
	ECHIGOYUZAWA = "越後湯沢"
	TOKYO        = "東京"
)

type WayServer struct {
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

func convLegToRespType(leg Leg) way.Run {
	return way.Run{
		LineName:         leg.LineName,
		DepartureStation: leg.DepartureStation,
		DepartureHour:    uint32(leg.DepartureHour),
		DepartureMinute:  uint32(leg.DepartureMinute),
		ArrivalStation:   leg.ArrivalStation,
		ArrivalHour:      uint32(leg.ArrivalHour),
		ArrivalMinute:    uint32(leg.ArrivalMinute),
	}
}

func getLegsFromDb(departureStation string, arrivalStation string, isHoliday bool) way.Leg {
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

	var legsRespType []*way.Run
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
	return way.Leg{
		Runs: legsRespType,
	}
}

func getLegs(curStation string, arrivalStation string, cur2nextStation map[string]string, isHoliday bool) []*way.Leg {
	var ret []*way.Leg
	for {
		if curStation == arrivalStation {
			break
		}
		println(curStation)
		nextStation, ok := cur2nextStation[curStation]
		if !ok {
			panic("err")
		}
		v := getLegsFromDb(curStation, nextStation, isHoliday)
		ret = append(ret, &v)
		curStation = nextStation
	}
	return ret
}

func getLegsHome(req *connect.Request[way.GetLinesRequest]) []*way.Leg {
	cur2nextStationHome := map[string]string{
		req.Msg.SkiResort: ECHIGOYUZAWA,
		ECHIGOYUZAWA:      TOKYO,
		TOKYO:             req.Msg.HometownStation,
	}
	return getLegs(req.Msg.SkiResort, req.Msg.HometownStation, cur2nextStationHome, req.Msg.IsHoliday)
}

func getLegsToSki(req *connect.Request[way.GetLinesRequest]) []*way.Leg {
	cur2nextStationToSki := map[string]string{
		TOKYO:        ECHIGOYUZAWA,
		ECHIGOYUZAWA: req.Msg.SkiResort,
	}
	if req.Msg.HometownStation != TOKYO {
		cur2nextStationToSki[req.Msg.HometownStation] = TOKYO
	}
	return getLegs(req.Msg.HometownStation, req.Msg.SkiResort, cur2nextStationToSki, req.Msg.IsHoliday)
}

func (s *WayServer) GetLines(ctx context.Context, req *connect.Request[way.GetLinesRequest]) (*connect.Response[way.GetLinesResponse], error) {
	res := connect.NewResponse(&way.GetLinesResponse{
		AllLegsToSki: getLegsToSki(req),
		AllLegsHome:  getLegsHome(req),
	})
	return res, nil
}

var db *sql.DB

func main() {
	println("main called")
	var err error
	db, err = sql.Open("sqlite3", "./train_bus_time/train_bus_time.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := &WayServer{}
	mux := http.NewServeMux()
	path, handler := wayv1connect.NewWayServiceHandler(server)
	mux.Handle(path, handler)
	corsHandler := cors.AllowAll().Handler(h2c.NewHandler(mux, &http2.Server{}))
	println("start serving")
	http.ListenAndServe(
		"0.0.0.0:8080",
		corsHandler,
	)
}
