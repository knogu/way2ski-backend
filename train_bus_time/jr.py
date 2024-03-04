import sqlite3
from dataclasses import dataclass, astuple
import time
import requests
from bs4 import BeautifulSoup
import re

from train_bus_time.leg import Leg

conn = sqlite3.connect('./train_bus_time.db')

def extract_time(pattern, row):
    match = re.search(pattern, row.text)
    return match.group(1) if match else None


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
def get_legs(url: str, line_name: str, is_holiday: bool):
    rows = get_soup(url).find_all("tr", class_="time")
    station_to_last_departure_time = {}
    for row in rows:
        a_tag = row.find("a")
        if not a_tag:
            continue
        station_name = a_tag.text
        departure_time_str = extract_time(r"(.*)\s発", row)
        arrival_time_str = extract_time(r"(.*)\s着", row)
        if not arrival_time_str:
            arrival_time_str = departure_time_str
        if station_name == "東京":
            arr_h, arr_m = map(int, arrival_time_str.split(":"))
            for station_name, departure_time in station_to_last_departure_time.items():
                # TODO: skip if tokyo
                dep_h, dep_m = map(int, departure_time.split(":"))
                save(Leg(
                    departure_station=station_name,
                    departure_hour=dep_h,
                    departure_minute=dep_m,
                    arrival_station="東京",
                    arrival_hour=arr_h,
                    arrival_minute=arr_m,
                    line_name=line_name,
                    is_holiday=is_holiday
                ))
            station_to_last_departure_time = {"東京": departure_time_str}
        else:
            # 東京 -> station in the current row
            if "東京" in station_to_last_departure_time:
                dep_h, dep_m = map(int, station_to_last_departure_time["東京"].split(":"))
                arr_h, arr_m = map(int, arrival_time_str.split(":"))
                save(Leg(
                    departure_station="東京",
                    departure_hour=dep_h,
                    departure_minute=dep_m,
                    arrival_station=station_name,
                    arrival_hour=arr_h,
                    arrival_minute=arr_m,
                    line_name=line_name,
                    is_holiday=is_holiday
                ))
            assert station_name not in station_to_last_departure_time
            station_to_last_departure_time[station_name] = arrival_time_str


def get_legs_from_timetable(timetable_url: str, line_name: str, is_holiday: bool):
    soup = get_soup(timetable_url)
    all_tr = soup.find_all('tr')
    for data in all_tr:
        a_tags = data.find_all('a')
        rel_urls_for_details = [a.get('href') for a in a_tags]
        for run_detail_url in rel_urls_for_details:
            match = re.search(r'train.*', run_detail_url)
            abs_url = 'https://www.jreast-timetable.jp/2402/' + match.group()
            get_legs(abs_url, line_name, is_holiday)
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

# 京葉線・平日・下り
# get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt1500/1500010.html",
#                         r"東京", r".*", "京葉線", False)
# # 京葉線・土休日・下り
# get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt1500/1500011.html",
#                         r"東京", r".*", "京葉線", True)
#
# # 京葉線・平日・上り
# get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt1500/1500020.html",
#                         r".*", r"東京", "京葉線", False)
# # 京葉線・土休日・上り
# get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt1500/1500021.html",
#                         r".*", r"東京", "京葉線", True)

# 山手線内回り・平日
get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt1039/1039120.html", "山手線内回り", False)
# 山手線外回り・平日
get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt1039/1039110.html", "山手線外回り", False)

# 山手線内回り・土休日
get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt1039/1039121.html", "山手線内回り", True)
# 山手線外回り・土休日
get_legs_from_timetable("https://www.jreast-timetable.jp/2402/timetable/tt1039/1039111.html", "山手線外回り", True)

conn.close()
