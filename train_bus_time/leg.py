from dataclasses import dataclass


@dataclass
class Leg:
    line_name: str
    departure_station: str
    departure_hour: int
    departure_minute: int
    arrival_station: str
    arrival_hour: int
    arrival_minute: int
    is_holiday: bool
    is_slow: bool

    def departure_in_minutes(self):
        return self.departure_hour * 60 + self.departure_minute

    def arrival_in_minutes(self):
        return self.arrival_hour * 60 + self. arrival_minute
