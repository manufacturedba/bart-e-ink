# bart-e-ink

North Berkeley BART Departures on ESP32-driven eInk

## Overview

This project is a simple display of the next southbound North Berkeley BART departures on an eInk display.

## Organization

- `display/` contains CircuitPython code for the eInk display
- `cmd/server` contains server code 

## Build

### Server

`go build ./cmd/server`

### Display

Populate the eInk device with the contents of `display/` directory including the `lib/` directory using the Adafruit CircuitPython library bundle.

You can find the bundle at https://circuitpython.org/libraries

## Disclaimer

This project currently requires manual fetching of GTFS data and assumes it is 
stored under `gtfs/` in the project directory. This is not included with the 
repository.

You can find a permalink to the GTFS data for BART at https://www.bart.gov/schedules/developers/gtfs

## Hardware

- [Adafruit HUZZAH32 – ESP32 Feather Board](https://www.adafruit.com/product/3405)
- [Adafruit 2.9" Tri-Color eInk Display FeatherWing](https://www.adafruit.com/product/4778)
  
Pins are defined according to the above hardware