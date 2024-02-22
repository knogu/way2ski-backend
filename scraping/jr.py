import pprint
from dataclasses import dataclass
import time
import requests
from bs4 import BeautifulSoup
import re


def extract_time(pattern, row):
    match = re.search(pattern, row.find("ul").text)
    return match.group(1) if match else None


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


def get_soup(url):
    response = requests.get(url)
    time.sleep(4)
    response.encoding = 'utf-8'
    return BeautifulSoup(response.text, 'html.parser')


# url: 列車詳細のurl
def get_legs(url: str, departure_station_pattern: str, arrival_station_pattern: str, line_name: str, is_holiday: bool) -> list[Leg]:
    rows = get_soup(url).find_all("tr", class_="time")
    legs: list[Leg] = []
    departure_times = {}  # station name -> departure time
    arrival_times = {}  # station name -> arrival time
    for row in rows:
        station_name = row.find("a").text
        departure_match = re.search(departure_station_pattern, station_name)
        if departure_match:
            departure_times[station_name] = extract_time(r"(.*)\s発", row)
        arrival_match = re.search(arrival_station_pattern, station_name)
        if arrival_match:
            arrival_times[station_name] = extract_time(r"(.*)\s着", row)

    for departure_station_name, departure_time in departure_times.items():
        for arrival_station_name, arrival_time in arrival_times.items():
            if departure_station_name == arrival_station_name:
                continue
            dep_h, dep_m = map(int, departure_time.split(":"))
            arr_h, arr_m = map(int, arrival_time.split(":"))
            assert all(0 <= h < 24 for h in [dep_h, arr_h])
            assert all(0 <= m < 60 for m in [dep_m, arr_m])
            legs.append(Leg(
                departure_station=departure_station_name,
                departure_hour=dep_h,
                departure_minute=dep_m,
                arrival_station=arrival_station_name,
                arrival_hour=arr_h,
                arrival_minute=arr_m,
                line_name=line_name,
                is_holiday=is_holiday
            ))

    return legs


def get_legs_from_timetable(timetable_url: str, departure_station_pattern: str, arrival_station_pattern: str, line_name: str, is_holiday: bool):
    soup = get_soup(timetable_url)
    all_tr = soup.find_all('tr')
    all_legs: list[Leg] = []
    for data in all_tr:
        a_tags = data.find_all('a')
        rel_urls_for_details = [a.get('href') for a in a_tags]
        for run_detail_url in rel_urls_for_details:
            match = re.search(r'train.*', run_detail_url)
            abs_url = 'https://www.jreast-timetable.jp/2402/' + match.group()
            all_legs.extend(get_legs(abs_url, departure_station_pattern, arrival_station_pattern, line_name, is_holiday))
            # break early while developing
            break
        break
    return all_legs


url_to_tokyo_holidays = 'https://www.jreast-timetable.jp/2402/timetable/tt1647/1647041.html'
url_from_tokyo_holidays = 'https://www.jreast-timetable.jp/2402/timetable/tt1647/1647031.html'

url_to_tokyo_weekdays = 'https://www.jreast-timetable.jp/2402/timetable/tt1647/1647040.html'
url_from_tokyo_weekdays = 'https://www.jreast-timetable.jp/2402/timetable/tt1647/1647030.html'

legs: list[Leg] = [
    *get_legs_from_timetable(url_to_tokyo_holidays, r".*", r"東京", "中央線", True),
    *get_legs_from_timetable(url_from_tokyo_holidays, r"東京", r".*", "中央線", True),
    *get_legs_from_timetable(url_to_tokyo_weekdays, r".*", r"東京", "中央線", False),
    *get_legs_from_timetable(url_from_tokyo_weekdays, r"東京", r".*", "中央線", False),
]

pprint.pprint(legs)
