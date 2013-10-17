Go-Roomba
===
A go library for interacting with iRobot Roomba or Create robots following the Open Interface (OI) specification.

Details
---
The code of the library is remotely inspired by `pyrobot` library by damonkohler@gmail.com (Damon Kohler)

Also, I have no idea how to write and/or launch Go code, but 

    git clone https://github.com/xa4a/go-roomba.git
    cd go-roomba
    export GOPATH=`pwd`
    go build roomba
    go install
    go run roomba_run.go

Should output:

    2012/10/29 16:41:40 Failed to open serial port: /dev/cu.usbserial-FTTL3AW0
    2012/10/29 16:41:40 Making roomba failed
    exit status 1
   
And if you have Roomba connected to the specified port (`/dev/cu.usbserial-FTTL3AW0` above) it may move forward a bit.
