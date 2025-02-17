#
#Copyright Odin Solutions S.L. All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#
import os
from capability.evaluator.SimpleCapabilityEvaluator import SimpleCapabilityEvaluator


pyCapabilityLib_folderpath = str(os.getenv('pyCapabilityLib_folderpath') if os.getenv('pyCapabilityLib_folderpath') is not None else '')

# In this class we intend to develop a series of tests with a two-fold purpose:
# 	1º - Provide a straightforward code that shows how the library works
#  2º - Verify that indeed everything is working as expected
#  Note: Not to confuse with Unit testing which are not performed per se here. 
# 

def pyCapabilityEvaluator(device, action, resource, captoken, subject, headers, body) :

    try:
        #Evaluating external capability tokens.
                
        verdict = SimpleCapabilityEvaluator.validateCapabilityToken(pyCapabilityLib_folderpath + "local_dependencies/capManagerUMU.der", 
                action,resource,device,subject,captoken,headers,body)
                
        return (verdict)

    except Exception as e:
        print("pyCapabilityEvaluator.pyCapabilityEvaluator(ERROR): Could not validate capability token")
        print(e)
        return None

if __name__ == '__main__':

    device 	= "https://localhost:1027";
    action 	= "GET";
    resource = "/ngsi-ld/v1/entities/?type=http://example.org/vehicle/Vehicle";
    captoken = "{\"id\": \"dnc7dlq0vsb66j57h6knka1n0b\",\"ii\": 1673460939,\"is\": \"capabilitymanager@odins.es\",\"su\": \"usuariouno\",\"de\": \"https://localhost:1027\",\"si\":\"MEUCIQD+sKqcRFgpxp7lYswZ3jFu8ALVcniKZiYP0siTxa5GIwIgSjKXY4DeZqYSxaEjIGeRd5ei7bvKtSydT0K6MhfO3KM=\",\"ar\": [{\"ac\": \"GET\",\"re\": \"/ngsi-ld/v1/entities/?type=http://example.org/vehicle/Vehicle\"}],\"nb\": 1673461939,\"na\": 1673471939}";
    subject 	= "";
    headers	= "";
    body		= "";


    verdict = pyCapabilityEvaluator(device, action, resource, captoken, subject, headers, body)

    print (verdict)

