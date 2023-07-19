# signalk-server-go

Start of a [signal K](https://signalk.org/specification/1.7.0/doc/) server implemented in go 

First only a n2k can source will be implemented with a subset of [canboat](https://signalk.org/specification/1.7.0/doc/) supported PGN's 


## Dependencies

go 1.20

## Get started

Get the webapps [Freeboard-sk](https://github.com/SignalK/freeboard-sk) and [Instrumentpanel](https://github.com/SignalK/instrumentpanel)

```
$ make webapps
```

```
$ make build
```
or
```
$ make buildarm
```
For a ARMv6 target

Run the server
```
$ ./signalk-server-go --file-source samples/nemo-n2k.txt
```

Point a browser at [http://localhost:3000/](https://localhost:3000/)

![Screenshot](img/screenshot.jpg)


For live data with a can device with socketcan support

```
./signalk-server-go --source can0
```

If you have an AIS connected to the N2K network

```
./signalk-server-go --mmsi <mmsi number> --source can0
```