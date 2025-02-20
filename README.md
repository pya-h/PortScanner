# Simple Port Scanner

This is a simple port scanner that scans a range of ports on a given host; Trying to be as fast as possible.

## Usage

```bash
go run portscanner.go -a <address> -f <first_port> -l <last_port>
```

## Example

```bash
go run portscanner.go -h 127.0.0.1 -first 1 -last 1000 -workers 550
```


## Note:
    The higher the count of workers(300 by default), the faster your program should execute. 
        But if you add too many workers, your results could become unreliable

## Default Values:
    - Host: localhost
    - First Port: 1
    - Last Port: 1024
    - Workers Count: 300

