 #!/bin/bash

    ps -ef | grep -v grep | grep "remail"
    # if not found - equals to 1, start it

    if [ $? -eq 1 ]
    then
        nohup ./remail > /home/go/github.com/reviashko/remail/logs_1 &
    fi