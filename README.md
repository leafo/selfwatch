
# selfwatch

Selfwatch is a program that monitors how you use your computer. It's inspired by
[selfspy](https://github.com/gurgeh/selfspy).

This project is in it's very early stages, currently it only tracks number of
keys pressed over time, and only runs on Linux. In the future it will collect
more detailed information about applications used and the things typed.

It supports sending key counts to a remote server, which can be used to create
a graph of your activity: <http://leafo.net/#typing>

## Usage

```
> selfwatch --help

Usage of selfwatch:
  -config string
    	Path to json config file (default "selfwatch.json")
```

## Config

The following options can be specified in the configuration json file:

* `DbName` - The name of the sqlite database to load to store data (default: `"selfwatch.json"`)
* `RemoteUrl` - A URL to flush key press counts to every `RemoteFlushDelay` seconds. Data is encoded as JSON and sent as a post request. It's formatted as an array of arrays: `[id, "YYYY:DD:MM HH:MM:SS", count]`
* `RemoteFlushDelay` - How long to wait between flushing key counts to remote server, default 60
* `SyncDelay` - How long to buffer key counts in memory before flushing to database (application switches will trigger immediate flush)

## About

Author: Leaf Corcoran (leafo) ([@moonscript](http://twitter.com/moonscript))  
Email: leafot@gmail.com  
Homepage: <http://leafo.net>  
License: MIT, Copyright (C) 2017 by Leaf Corcoran


