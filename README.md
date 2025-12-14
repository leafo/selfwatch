
# selfwatch

Selfwatch is a program that monitors how you use your computer. It's inspired by
[selfspy](https://github.com/gurgeh/selfspy). It tracks the number of keys
pressed over time using the X11 RECORD extension on Linux.

![selfwatch screenshot](screenshot.png)

## Usage

```
> selfwatch --help

Usage of selfwatch:
  -config string
    	Path to json config file (default "selfwatch.json")
```

## Web Mode

Selfwatch includes a built-in web server that provides a dashboard for visualizing your typing activity. Start it with:

```
> selfwatch web [address]
```

The default address is `:8080`. You can specify a different address:

```
> selfwatch web :9000
> selfwatch web localhost:8080
```

The dashboard displays:

* **Last 24 hours** - Hourly breakdown of key presses, navigable by date
* **Last 30 days** - Daily key press totals
* **Yearly activity** - A contribution grid similar to GitHub's activity graph

## Config

The following options can be specified in the configuration json file:

* `DbName` - The name of the sqlite database to load to store data (default: `"selfwatch.json"`)
* `RemoteUrl` - A URL to flush key press counts to every `RemoteFlushDelay` seconds. Data is encoded as JSON and sent as a post request. It's formatted as an array of arrays: `[id, "YYYY:DD:MM HH:MM:SS", count]`
* `RemoteFlushDelay` - How long to wait between flushing key counts to remote server, default 60
* `SyncDelay` - How long to buffer key counts in memory before flushing to database (application switches will trigger immediate flush)
* `NewDayHour` - The hour (0-23) when a new day starts for statistics purposes (default: 4). Useful if you work late nights and want activity after midnight counted as part of the previous day

## About

Author: Leaf Corcoran (leafo) ([@moonscript](http://twitter.com/moonscript))  
Email: leafot@gmail.com  
Homepage: <http://leafo.net>  
License: MIT, Copyright (C) 2017 by Leaf Corcoran


