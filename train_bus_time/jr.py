import sqlite3
from dataclasses import dataclass, astuple
import time
import requests
from bs4 import BeautifulSoup
import re


conn = sqlite3.connect('./train_bus_time.db')

def extract_time(pattern, row):
    match = re.search(pattern, row.text)
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

def save(leg: Leg):
    try:
        conn.execute("INSERT INTO legs VALUES (?,?,?,?,?,?,?,?)", astuple(leg))
    except sqlite3.IntegrityError:
        pass
    conn.commit()

def get_soup(url):
    time.sleep(10)
    print("requesting", url)
    response = requests.get(url)
    response.encoding = 'utf-8'
    return BeautifulSoup(response.text, 'html.parser')


# url: 列車詳細のurl
def get_legs(url: str, departure_station_pattern: str, arrival_station_pattern: str, line_name: str, is_holiday: bool, either_must_be = "", is_departure_arrival_check_loose = False):
    rows = get_soup(url).find_all("tr", class_="time")
    departure_times = {}  # station name -> departure time
    arrival_times = {}  # station name -> arrival time
    for row in rows:
        a_tag = row.find("a")
        if not a_tag:
            continue
        station_name = a_tag.text
        departure_match = re.search(departure_station_pattern, station_name)
        if departure_match:
            time_str = extract_time(r"(.*)\s発", row)
            if time_str:
                departure_times[station_name] = time_str
            elif is_departure_arrival_check_loose:
                if station_name not in departure_times:
                    print("INFO only arrival time found in yamanote-line", station_name)
            else:
                assert station_name == "東京" or station_name == "越後湯沢"
        arrival_match = re.search(arrival_station_pattern, station_name)
        if arrival_match:
            time_str = extract_time(r"(.*)\s着", row)
            if time_str:
                arrival_times[station_name] = time_str
            elif is_departure_arrival_check_loose:
                time_str = extract_time(r"(.*)\s発", row)
                if time_str:
                    arrival_times[station_name] = time_str
                else:
                    raise Exception("both of arrival and departure not found in", station_name)
            else:
                assert station_name == "東京" or station_name == "越後湯沢"
        # print(station_name, arrival_match, departure_match)

    for departure_station_name, departure_time in departure_times.items():
        for arrival_station_name, arrival_time in arrival_times.items():
            # print(departure_station_name, arrival_station_name)
            if either_must_be != "" and (departure_station_name != either_must_be and arrival_station_name != either_must_be):
                # print("either")
                continue
            if departure_station_name == arrival_station_name:
                continue
            dep_h, dep_m = map(int, departure_time.split(":"))
            arr_h, arr_m = map(int, arrival_time.split(":"))
            assert all(0 <= h for h in [dep_h, arr_h])
            if dep_h >= 24:
                print("INFO departure >= 24", departure_station_name, departure_time)
            if arr_h >= 24:
                print("INFO arrival >= 24", arrival_station_name, arrival_time)

            assert all(0 <= m < 60 for m in [dep_m, arr_m])
            leg = Leg(
                departure_station=departure_station_name,
                departure_hour=dep_h,
                departure_minute=dep_m,
                arrival_station=arrival_station_name,
                arrival_hour=arr_h,
                arrival_minute=arr_m,
                line_name=line_name,
                is_holiday=is_holiday
            )
            save(leg)


def get_legs_from_timetable(timetable_url: str, departure_station_pattern: str, arrival_station_pattern: str, line_name: str, is_holiday: bool, either_must_be = "", use_departure_time_for_arrival_if_not_found = False):
    soup = get_soup(timetable_url)
    all_tr = soup.find_all('tr')
    for data in all_tr:
        a_tags = data.find_all('a')
        rel_urls_for_details = [a.get('href') for a in a_tags]
        for run_detail_url in rel_urls_for_details:
            match = re.search(r'train.*', run_detail_url)
            abs_url = 'https://www.jreast-timetable.jp/2402/' + match.group()
            get_legs(abs_url, departure_station_pattern, arrival_station_pattern, line_name, is_holiday, either_must_be, use_departure_time_for_arrival_if_not_found)
            # break early while developing
            # break
        # break


url_to_tokyo_holidays = 'https://www.jreast-timetable.jp/2402/timetable/tt1647/1647041.html'
url_from_tokyo_holidays = 'https://www.jreast-timetable.jp/2402/timetable/tt1647/1647031.html'

url_to_tokyo_weekdays = 'https://www.jreast-timetable.jp/2402/timetable/tt1647/1647040.html'
url_from_tokyo_weekdays = 'https://www.jreast-timetable.jp/2402/timetable/tt1647/1647030.html'

# get_legs_from_timetable(url_to_tokyo_holidays, r".*", r"東京", "中央線", True),
# get_legs_from_timetable(url_from_tokyo_holidays, r"東京", r".*", "中央線", True),
# get_legs_from_timetable(url_to_tokyo_weekdays, r".*", r"東京", "中央線", False),
# get_legs_from_timetable(url_from_tokyo_weekdays, r"東京", r".*", "中央線", False)
# # 上越新幹線
# get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt0285/0285031.html",
#                         r"越後湯沢", r"東京", "上越新幹線", True),
# get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt1039/1039051.html",
#                         r"東京", r"越後湯沢", "上越新幹線", True),
# get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt0285/0285030.html",
#                         r"越後湯沢", r"東京", "上越新幹線", False),
# get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt1039/1039050.html",
#                              r"東京", r"越後湯沢", "上越新幹線", False),

# 山手線内回り・平日
# get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt1039/1039120.html",
#                         r".*", r".*", "山手線", False, "東京", True)
# 山手線外回り・平日
# get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt1039/1039110.html",
#                         r".*", r".*", "山手線", False, "東京", True),

# 山手線内回り・土休日
get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt1039/1039121.html",
                        r".*", r".*", "山手線内回り", True, "東京", True)
# 山手線外回り・土休日
get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt1039/1039111.html",
                        r".*", r".*", "山手線外回り", True, "東京", True)


conn.close()
