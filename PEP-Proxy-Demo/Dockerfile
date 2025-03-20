#
#Copyright Odin Solutions S.L. All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

FROM python:3.11.2-slim

#FROM ubuntu:20.04

RUN apt-get update -y
#RUN apt-get update --fix-missing --> (try if the previous RUN apt-get update -y  (line 9) doesn't work)

##Install python3
#RUN apt-get install software-properties-common -y
##RUN apt-get update --> (try if the previous RUN apt-get update -y  (line 9) doesn't work)
#RUN add-apt-repository ppa:deadsnakes/ppa
#RUN apt-get install python3 -y
#
##Install pip3
#RUN apt-get install python3-pip -y
#
#Install java JRE
#RUN apt-get install openjdk-11-jre -y
##RUN pip3 install install-jdk ??

##or 
##RUN apt install default-jre -y --> It works but is not used to have control about the version.
#
##Install utils
#RUN apt-get install net-tools -y
#RUN apt-get install vim -y
#RUN apt-get install curl -y
#


# Establish workdir
WORKDIR /opt/PEP-Proxy

#Install requirements
COPY requirements.txt /opt/PEP-Proxy/
RUN pip3 install -r /opt/PEP-Proxy/requirements.txt

RUN pip3 install cryptography

RUN pip3 install "PyJWT==2.3.0"

#Transfer source and neccesary files
COPY API /opt/PEP-Proxy/API
COPY certs /opt/PEP-Proxy/certs
COPY PEP-Proxy.py UtilsPEP.py config.cfg cipher_configuration.json /opt/PEP-Proxy/
COPY pyCapabilityLib ./pyCapabilityLib

RUN touch out.log

#Broker protocol. Admitted values: "http","https"
#ENV target_protocol=http
#Broker host.
#ENV target_host=broker
#ENV target_port=9090

#Validate Capability token using blockchain: Admitted values: "0: No use; 1:Use"
#ENV blockchain_usevalidation=1

#BlockChain protocol. Admitted values: "http","https"
#ENV blockchain_protocol=http
#BlockChain host.
#ENV blockchain_host=blockchain
#ENV blockchain_port=8000

#ENV PEP_ENDPOINT=https://pephost:1027

#capability_lib - VARIABLES
ENV pyCapabilityLib_folderpath=pyCapabilityLib/
#Define pyCapabilityLib_verifyMethod parameter to choose the method to verify the signature the Capability Token: Admitted values:
# - certsFile --> to use python code based in the logic of the original java .jar executable file. (capabilityEvaluartor) (RECOMMENDED - fastest option)
# - jarCPABE --> to use CP-ABE (.jar executable JAVA) (slowest option -Change if certsFiles is not working until the issue is solved)
ENV pyCapabilityLib_verifyMethod=certsFile

# application's default port
EXPOSE 1027

# Launch app
CMD [ "python3", "/opt/PEP-Proxy/PEP-Proxy.py" ]
