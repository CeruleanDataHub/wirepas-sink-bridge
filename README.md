# wirepas-sink-bridge
Reads and decodes messages from a Wirepas sink and encodes them into JSON and prints them out, optionally also writing them out to a Unix socket.

## Docker
### Build image
```shell
    docker build -t wirepas-sink-bridge .
```

### Run image
```shell
    docker run --device=/dev/ttyUSB0 wirepas-sink-bridge
```
