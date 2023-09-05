# signalk-server-go

Start of a [signal K](https://signalk.org/specification/1.7.0/doc/) server implemented in go 

First only a n2k can source will be implemented with a subset of [canboat](https://signalk.org/specification/1.7.0/doc/) supported PGN's 


## Dependencies

go 1.20

## Get started

### signalk-server-go

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
$ ./build/signalk-server-go --file-source samples/nemo-n2k.txt
```

Point a browser at [http://localhost:3000/](https://localhost:3000/)

![Screenshot](img/screenshot.jpg)


For live data with a can device with socketcan support

```
./build//signalk-server-go --source can0
```

If you have an AIS connected to the N2K network

```
./build/signalk-server-go --mmsi <mmsi number> --source can0
```

More options

```
$ build/signalk-server-go --help
Usage of build/signalk-server-go:
  -debug
        Enable debugging
  -file-source value
        Path to candump file
  -mmsi string
        Vessel MMSI
  -port int
        Listen port (default 3000)
  -source value
        Source Can device
  -tls
        Enable tls
  -tlscert string
        Tls certificate file
  -tlskey string
        Tls key file
  -version
        Show version
  -webapp-path string
        Path to webapps (default "./static")
  -webapps
        Serve webapps (default true)
```

