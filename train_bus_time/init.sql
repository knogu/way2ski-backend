CREATE TABLE legs (
    line_name VARCHAR(30),

    departure_station VARCHAR(30),
    departure_hour INTEGER,
    departure_minute INTEGER,

    arrival_station VARCHAR(30),
    arrival_hour INTEGER,
    arrival_minute INTEGER,

    is_holiday INTEGER,
    UNIQUE (line_name, departure_station, departure_hour, departure_minute, arrival_station, arrival_hour, arrival_minute, is_holiday)
)
