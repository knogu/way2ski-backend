package main

import (
	"context"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/jmoiron/sqlx"
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

func getLegFromDb(departureStation string, arrivalStation string, isHoliday bool) way.Leg {
	rows, err := db.Query(`
		SELECT * FROM legs 
		WHERE departure_station = ? 
		AND arrival_station = ? 
		AND is_holiday = ?
		ORDER BY departure_hour, departure_minute
	`, departureStation, arrivalStation, isHoliday)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var runs []*way.Run
	for rows.Next() {
		var leg Leg
		err = rows.Scan(&leg.LineName, &leg.DepartureStation, &leg.DepartureHour, &leg.DepartureMinute, &leg.ArrivalStation, &leg.ArrivalHour, &leg.ArrivalMinute, &leg.IsHoliday)
		if err != nil {
			log.Fatal(err)
		}
		v := convLegToRespType(leg)
		runs = append(runs, &v)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return way.Leg{
		DepartureStation: departureStation,
		ArrivalStation:   arrivalStation,
		Runs:             runs,
	}
}

func getHometownStationsFromDb() []string {
	lineNames := []string{"中央線", "山手線"}

	query := `
    SELECT DISTINCT departure_station FROM legs
    WHERE line_name IN (?)`
	query, args, err := sqlx.In(query, lineNames)
	if err != nil {
		log.Fatalln(err)
	}

	query = db.Rebind(query)
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatalln(err)
	}

	var hometownStations []string
	for rows.Next() {
		var result string
		if err := rows.Scan(&result); err != nil {
			log.Fatal(err)
		}
		hometownStations = append(hometownStations, result)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return hometownStations
}

func getLegs(curStation string, arrivalStation string, cur2nextStation map[string]string, isHoliday bool) []*way.Leg {
	var ret []*way.Leg
	for {
		if curStation == arrivalStation {
			break
		}
		nextStation, ok := cur2nextStation[curStation]
		if !ok {
			panic("err")
		}
		v := getLegFromDb(curStation, nextStation, isHoliday)
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

func (s *WayServer) GetHometownStations(ctx context.Context, req *connect.Request[way.GetHometownStationsRequest]) (*connect.Response[way.GetHometownStationsResponse], error) {
	res := connect.NewResponse(&way.GetHometownStationsResponse{
		HometownStations: getHometownStationsFromDb(),
	})
	return res, nil
}

var db *sqlx.DB

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func main() {
	var err error
	db, err = sqlx.Connect("sqlite3", "./train_bus_time/train_bus_time.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := &WayServer{}
	mux := http.NewServeMux()
	path, handler := wayv1connect.NewWayServiceHandler(server)
	mux.Handle(path, handler)
	mux.Handle("/health", http.HandlerFunc(health))
	corsHandler := cors.AllowAll().Handler(h2c.NewHandler(mux, &http2.Server{}))
	println("start serving")
	http.ListenAndServe(
		"0.0.0.0:8080",
		corsHandler,
	)
}
