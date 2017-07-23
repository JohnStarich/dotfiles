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


def weather(latitude: float, longitude: float) -> dict:
    response = requests.get(
        'https://query.yahooapis.com/v1/public/yql',
        params={
            'q': """
            select *
            from weather.forecast
            where woeid in (
                SELECT woeid FROM geo.places WHERE text="({lat},{lon})"
            )
            """.format(lat=latitude, lon=longitude),
            'format': 'json',
        },
        timeout=5,
    )
    json = response.json()['query']
    if json['results'] is None:
        raise Exception('No weather found for this location')
    results = json['results']['channel']
    return {
        'code': int(results['item']['condition']['code']),
        'coordinates': {
            'latitude': float(results['item']['lat']),
            'longitude': float(results['item']['long']),
        },
        'humidity': float(results['atmosphere']['humidity']),
        'pressure': float(results['atmosphere']['pressure']),
        'temperature': float(results['item']['condition']['temp']),
        'wind': {
            'chill': float(results['wind']['chill']),
            'direction': float(results['wind']['direction']),
            'speed': float(results['wind']['speed']),
        },
        'units': results['units'],
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
