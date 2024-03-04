import sqlite3
from dataclasses import asdict

from leg import Leg

conn = sqlite3.connect('./train_bus_time.db')

def mark_is_slow(conn, leg):
    leg_dict = asdict(leg)

    where_clause = ' AND '.join(f'{key} = :{key}' for key in leg_dict.keys() if key != "is_slow")

    conn.execute(f'UPDATE legs SET is_slow = 1 WHERE {where_clause}', leg_dict)
    conn.commit()

def mark_is_slow_from_rows(rows):
    legs = []
    for row in rows:
        legs.append(Leg(*row))

    legs.sort(key=lambda leg: leg.departure_in_minutes())
    slow_train_indices = set()
    for i, checked_leg in enumerate(legs):
        for j in range(i + 1, len(legs)):
            later_departure_leg = legs[j]
            if later_departure_leg.arrival_in_minutes() <= checked_leg.arrival_in_minutes():
                slow_train_indices.add(i)
                # print(checked_leg, "will be deleted. better leg:", later_departure_leg)
                break

    for idx in slow_train_indices:
        mark_is_slow(conn, legs[idx])


def handle(station_name):
    cursor = conn.cursor()
    params_list = [
        # (departure_station, arrival_station, is_holiday)
        (station_name, "東京", False),
        (station_name, "東京", True),
        ("東京", station_name, False),
        ("東京", station_name, True),
    ]
    for param in params_list:
        cursor.execute("select * from legs where departure_station=? and arrival_station=? and is_holiday=?", param)
        mark_is_slow_from_rows(cursor.fetchall())

def extract_station_names():
    cursor = conn.cursor()
    cursor.execute("select distinct departure_station from legs")
    res = []
    for row in cursor.fetchall():
        assert len(row) == 1
        res.append(row[0])
    return res

for station_name in extract_station_names():
    handle(station_name)
