#!/bin/bash

checkService() {
	OK=0

	if [ $# -eq 6 ]; then
		I=0

		while [ $I -lt 3 ]; do
			#echo $1
			#echo $2
			#echo $3
			#echo $4

			request_cmd=$(curl --max-time 10 -s -i -k -o /dev/null -w "%{http_code}" -X GET "$1://$2:$3$4")
			#echo $request_cmd
			if [ $request_cmd -eq 200 ]; then
				# ok
				OK=1
				I=3

#				echo "checkService: OK / container: '$5' - {"protocol": '$1', 'host': '$2', 'port': $3, 'resource': '$4'"
			else
				let I=$I+1
				#echo $I
			fi
		done
	else
		return
	fi
	#echo $OK
	if [ $OK -ne 1 ]; then
		echo `date` " : ERROR - $1://$2:$3$4 - Not working"
		cd $5 && docker compose restart $6
	fi
}

# XACML
PROTOCOL="http"
HOST="localhost"
PORT=9092
URI="/pdp/test"
PROJECTPATH="/home/odins/deployment/securitycomponents/deploys/{{subfolder}}"
CONTAINER="xacml-pdp"

checkService $PROTOCOL $HOST $PORT $URI $PROJECTPATH $CONTAINER

exit 0

