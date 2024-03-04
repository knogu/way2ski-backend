INSERT INTO legs (
    line_name,
    departure_station,
    departure_hour,
    departure_minute,
    arrival_station,
    arrival_hour,
    arrival_minute,
    is_holiday
) VALUES
      /* 休日・行き */
      ('舞子シャトルバス', '越後湯沢', 7, 50, '舞子スノーリゾート', 8, 20, 1),
      ('舞子シャトルバス', '越後湯沢', 8, 10, '舞子スノーリゾート', 8, 40, 1),
      ('舞子シャトルバス', '越後湯沢', 8, 30, '舞子スノーリゾート', 9, 0, 1),
      ('舞子シャトルバス', '越後湯沢', 8, 50, '舞子スノーリゾート', 9, 20, 1),
      ('舞子シャトルバス', '越後湯沢', 9, 10, '舞子スノーリゾート', 9, 40, 1),
      ('舞子シャトルバス', '越後湯沢', 9, 30, '舞子スノーリゾート', 10, 0, 1),
      ('舞子シャトルバス', '越後湯沢', 9, 50, '舞子スノーリゾート', 10, 20, 1),
      ('舞子シャトルバス', '越後湯沢', 10, 10, '舞子スノーリゾート', 10, 40, 1),
      ('舞子シャトルバス', '越後湯沢', 11, 0, '舞子スノーリゾート', 11, 30, 1),
      ('舞子シャトルバス', '越後湯沢', 12, 0, '舞子スノーリゾート', 12, 30, 1),

      /* 休日・帰り */
      ('舞子シャトルバス', '舞子スノーリゾート', 13, 0, '越後湯沢', 13, 30, 1),
      ('舞子シャトルバス', '舞子スノーリゾート', 14, 0, '越後湯沢', 14, 30, 1),
      ('舞子シャトルバス', '舞子スノーリゾート', 15, 0, '越後湯沢', 15, 30, 1),
      ('舞子シャトルバス', '舞子スノーリゾート', 15, 30, '越後湯沢', 16, 0, 1),
      ('舞子シャトルバス', '舞子スノーリゾート', 16, 0, '越後湯沢', 16, 30, 1),
      ('舞子シャトルバス', '舞子スノーリゾート', 16, 20, '越後湯沢', 16, 40, 1),
      ('舞子シャトルバス', '舞子スノーリゾート', 16, 40, '越後湯沢', 17, 0, 1),
      ('舞子シャトルバス', '舞子スノーリゾート', 17, 0, '越後湯沢', 17, 30, 1),
      ('舞子シャトルバス', '舞子スノーリゾート', 17, 20, '越後湯沢', 17, 40, 1),
      ('舞子シャトルバス', '舞子スノーリゾート', 17, 40, '越後湯沢', 18, 0, 1),
      ('舞子シャトルバス', '舞子スノーリゾート', 18, 0, '越後湯沢', 18, 30, 1),
      ('舞子シャトルバス', '舞子スノーリゾート', 18, 20, '越後湯沢', 18, 40, 1),

    /* 平日・行き */
      ('舞子シャトルバス', '越後湯沢', 7, 50, '舞子スノーリゾート', 8, 20, 0),
      ('舞子シャトルバス', '越後湯沢', 8, 10, '舞子スノーリゾート', 8, 40, 0),
      ('舞子シャトルバス', '越後湯沢', 8, 30, '舞子スノーリゾート', 9, 0, 0),
      ('舞子シャトルバス', '越後湯沢', 8, 50, '舞子スノーリゾート', 9, 20, 0),
      ('舞子シャトルバス', '越後湯沢', 9, 10, '舞子スノーリゾート', 9, 40, 0),
      ('舞子シャトルバス', '越後湯沢', 9, 30, '舞子スノーリゾート', 10, 0, 0),
      ('舞子シャトルバス', '越後湯沢', 10, 10, '舞子スノーリゾート', 10, 40, 0),
      ('舞子シャトルバス', '越後湯沢', 11, 0, '舞子スノーリゾート', 11, 30, 0),
      ('舞子シャトルバス', '越後湯沢', 12, 0, '舞子スノーリゾート', 12, 30, 0),

      /* 平日・帰り */
      ('舞子シャトルバス', '舞子スノーリゾート', 13, 30, '越後湯沢', 14, 0, 0),
      ('舞子シャトルバス', '舞子スノーリゾート', 14, 30, '越後湯沢', 15, 0, 0),
      ('舞子シャトルバス', '舞子スノーリゾート', 15, 30, '越後湯沢', 16, 0, 0),
      ('舞子シャトルバス', '舞子スノーリゾート', 16, 0, '越後湯沢', 16, 30, 0),
      ('舞子シャトルバス', '舞子スノーリゾート', 16, 20, '越後湯沢', 16, 40, 0),
      ('舞子シャトルバス', '舞子スノーリゾート', 16, 40, '越後湯沢', 17, 0, 0),
      ('舞子シャトルバス', '舞子スノーリゾート', 17, 0, '越後湯沢', 17, 30, 0),
      ('舞子シャトルバス', '舞子スノーリゾート', 17, 20, '越後湯沢', 17, 40, 0),
      ('舞子シャトルバス', '舞子スノーリゾート', 17, 40, '越後湯沢', 18, 0, 0);