#
#Copyright Odin Solutions S.L. All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

#FROM ubuntu:20.04
FROM python:3.11.2-slim

RUN apt-get update -y
# --> (try if the previous RUN apt-get update -y  (line 9) doesn't work)
#RUN apt-get update --fix-missing 

##Install python3
#RUN apt-get install software-properties-common -y

##--> (try if the previous RUN apt-get update -y  (line 9) doesn't work)
#RUN apt-get update 

#RUN add-apt-repository ppa:deadsnakes/ppa
#RUN apt-get install python3 -y

##Install pip3
#RUN apt-get install python3-pip -y

## Install Requests
RUN pip3 install requests

#Install utils
RUN apt-get install net-tools -y
RUN apt-get install vim -y
RUN apt-get install curl -y

RUN pip3 install "xmltodict==0.13.0"

# Establish workdir
WORKDIR /opt/PDP

#Transfer source and neccesary files
COPY PDP.py config.cfg ./
COPY certs ./certs

RUN touch out.log

# application's default port
EXPOSE 8080

# Launch app
CMD [ "python3", "/opt/PDP/PDP.py" ]
