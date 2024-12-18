#!/bin/bash

checkService() {
	OK=0

	if [ $# -eq 2 ]; then
		I=0

		while [ $I -lt 3 ]; do
			request_cmd=$(curl --max-time 10 -s -i -o /dev/null -w "%{http_code}" -X GET "{{Protocol}}://{{IP}}:$1" --header 'Accept: application/json')
			#echo $request_cmd
			if [ $request_cmd -eq 200 ]; then
				# ok
				OK=1
				I=3

#				echo "checkPEPProxy: OK / port $1 container $2"
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
		echo `date` " : ERROR - {{Protocol}}://{{IP}}:$1 - Not working"
		cd {{projectPath}} / PEP-Proxy && /usr/local/bin/docker-compose restart $2
	fi
}

# PEP-Proxy
PORT=1027
CONTAINER="pepproxy"

checkService $PORT $CONTAINER

exit 0

