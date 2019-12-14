#!/usr/bin/env python3

import subprocess
import requests


def location() -> (float, float):
    try:
        process = subprocess.run(
            ['CoreLocationCLI', '-once', 'YES',
             '-format', '%latitude\n%longitude'],
            stdout=subprocess.PIPE,
            encoding='utf8',
            timeout=2,
        )
    except FileNotFoundError:
        print("CoreLocationCLI not installed")
        return (None, None)
    lines = str(process.stdout).splitlines()
    if len(lines) < 2:
        return (None, None)
    latitude = lines[0]
    longitude = lines[1]
    return (latitude, longitude)


def first_value(d: dict):
    return d['values'][0]['value']


def weather(latitude: float, longitude: float) -> dict:
    response = requests.get(
        'https://api.weather.gov/points/{lat},{lon}'
        .format(lat=latitude, lon=longitude),
        timeout=5,
    )
    forecast_url = response.json()['properties']['forecastGridData']
    response = requests.get(forecast_url, timeout=5)
    json = response.json()
    actual_lat, actual_lon = json['geometry']['coordinates'][0][0]
    props = json['properties']
    return {
        'code': 0,
        'coordinates': {
            'latitude': actual_lat,
            'longitude': actual_lon,
        },
        'humidity': first_value(props['relativeHumidity']) / 100.0,
        'apparentTemperature': first_value(props['apparentTemperature']),
        'temperature': first_value(props['temperature']),
        'wind': {
            'chill': first_value(props['windChill']),
            'direction': first_value(props['windDirection']),
            'speed': first_value(props['windSpeed']),
        },
        'units': {
            'temperature': 'C',
            'apparentTemperature': 'C',
        },
    }


last_location = (None, None)


def raw_weather_info():
    global last_location
    latitude, longitude = location()
    error = None
    if latitude is None or longitude is None:
        error = 'Error getting current location.'
        if last_location[0] is None or last_location[1] is None:
            return {'error': error}
        latitude = last_location[0]
        longitude = last_location[1]
    else:
        last_location = (latitude, longitude)
    try:
        info = weather(latitude, longitude)
        if error is not None:
            info['error'] = error
        return info
    except Exception as e:
        return {'error': 'Error parsing weather information: ' + str(e)}


if __name__ == '__main__':
    print(raw_weather_info())
