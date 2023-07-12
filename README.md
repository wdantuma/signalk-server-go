# signalk-server-go

Start of a [signal K](https://signalk.org/specification/1.7.0/doc/) server in go 

First only a n2k can source will be implemented with a subset of [canboat](https://signalk.org/specification/1.7.0/doc/) supported PGN's 


## Dependencies

go 1.20

## Get started

Get the webapps [Freeboard-sk](https://github.com/SignalK/freeboard-sk) and [Instrumentpanel](https://github.com/SignalK/instrumentpanel)

```
make webapps
```

```
make run 
```

Point a browser at [http://localhost:3000/](https://localhost:3000/)

![Screenshot](img/screenshot.jpg)