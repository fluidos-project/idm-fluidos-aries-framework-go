#
#Copyright Odin Solutions S.L. All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

# Create the class to evaluate a Capability Token.
from .CapabilityVerifierCode import *

import time
import re

from subprocess import Popen, PIPE

from cryptography import x509
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives import hashes
import base64
import os

pyCapabilityLib_folderpath = str(os.getenv('pyCapabilityLib_folderpath') if os.getenv('pyCapabilityLib_folderpath') is not None else '')
pyCapabilityLib_verifyMethod = str(os.getenv('pyCapabilityLib_verifyMethod') if os.getenv('pyCapabilityLib_verifyMethod') == 'jarCPABE' else 'certsFile')

class CapabilityEvaluator :

    #X509Certificate certificate;
    certificate = None

    # Constructor of the class that contains the certificate to verify the token.
    # 
    # @param certificatePath

    def __init__(self,certificatePath) :

        try:

            #https://cryptography.io/en/latest/x509/reference/#cryptography.x509.load_der_x509_certificate
            with open(certificatePath, "rb") as f:
                self.certificate = x509.load_der_x509_certificate(f.read())

        except Exception as e:
            print("CapabilityEvaluator.__init__(ERROR): The certificate in " + certificatePath + " could not be processed as a X.509 certificate")
            print(e)
            self.certificate = None

    # Given a capability token, the action (normally a REST command) to perform over a service, 
    # and the device and resource on which the action is intended to be performed.
    # @param capability_token
    # @param action
    # @param device_resource
    # @return CapabilityVerifierCode

    def evaluateCapabilityToken(self, capability_token, action, resource, device, subject, headers, body) : 

        try:

            #print("*************evaluateCapabilityToken************")
            if ((capability_token == None) or (action == None) or (resource == None) or (device == None)):
                #print("Error: The token is not Valid")
                return CapabilityVerifierCode.TOKEN_NOT_VALID.name
            
            if (self.isTimeValid(capability_token.getIssueIns(), capability_token.getNb(), capability_token.getNa()) == False):
                #print("Error: The time is not Valid")	
                return CapabilityVerifierCode.OUTATIME.name

            result = self.actionIsPermitted(capability_token, action, resource, device)

            if (result != CapabilityVerifierCode.AUTHORIZED.name):
                #print("Error: The action is not permitted")
                return result            
                    
            if pyCapabilityLib_verifyMethod == 'certsFile':
                signatureIsValid = self.isSignatureValid(capability_token)
            else:
                signatureIsValid = self.isSignatureValid_CPABE(capability_token)

            if (signatureIsValid == False):
                #print("Error: The signature is not Valid")
                #print("The action is NOT authorized")
                return CapabilityVerifierCode.SIGNATURE_NOT_VALID.name

            #print("The action is authorized")
            return CapabilityVerifierCode.AUTHORIZED.name

        except Exception as e:
            print("CapabilityEvaluator.evaluateCapabilityToken(ERROR): Could not validate capability token")
            print(e)
            return CapabilityVerifierCode.TOKEN_NOT_VALID.name

    #  Given a capability token, we evaluate if the action to be performed, matches the action requested on the service of the resource's device.
    #  
    # @param ct
    # @param action
    # @param device_resource
    # @return
    
    def actionIsPermitted(self, ct, actionReq, resourceReq, deviceReq):

        try:
            if (ct == None):
                return CapabilityVerifierCode.TOKEN_NOT_VALID.name

            ar = ct.getAr()

            for i in range(len(ar)):

                #print("Permitted Action: " + ar[i]["ac"] + " - Length: " + str(len(ar[i]["ac"])))
                #print("Intented  Action: " + actionReq + " - Length: " + str(len(actionReq)))
                #print("Permitted Resource: " + ar[i]["re"] + " - Length: " + str(len(ar[i]["re"])))
                #print("Intented  Resource: " + resourceReq  + " - Length: " + str(len(resourceReq)))
                #print("Permitted Device: " + ct.getDe() + " - Length: " + str(len(ct.getDe())))
                #print("Intented  Device: " + deviceReq  + " - Length: " + str(len(deviceReq)))

                if actionReq != ar[i]["ac"]:
                    #print("Action is not equal")
                    return CapabilityVerifierCode.ACTION_NOT_PERMITTED.name;

                if resourceReq != ar[i]["re"]:
                
                    #print( "ar[i]["re"]: " + ar[i]["re"]);
                    #print( "resourceReq: " + resourceReq);
                    
                    # Replace \\\\ --> \\
                    resourceCapToken = ar[i]["re"].replace('\\\\', '\\')
                    resourceCapToken = resourceCapToken.replace('?', '\\?');
                    #print( "resourceCapToken: " + resourceCapToken);

                    pattern1 = re.compile("^" + resourceCapToken + "$")
                    matcher1 = pattern1.match (resourceReq)

                    if ( matcher1 == None):

                        pattern2 = re.compile("^" + resourceCapToken + "$")
                        matcher2 = pattern2.match ("\\Q" + resourceReq + "\\E" )

                        if ( matcher2 == None): 
                            #print("Resource is not allowed")
                            return CapabilityVerifierCode.RESOURCE_NOT_PERMITTED.name

            if(deviceReq != ct.getDe()):

                #print( "ct.getDe(): " + ct.getDe());
                #print( "deviceReq: " + deviceReq);

                # Replace \\\\ --> \\
                deviceCapToken = ct.getDe().replace('\\\\', '\\')

                # print( "deviceCapToken: " + deviceCapToken);

                pattern1 = re.compile("^" + deviceCapToken + "$")
                matcher1 = pattern1.match (deviceReq)

                if ( matcher1 == None):

                    pattern2 = re.compile("^" +  deviceCapToken + "$")
                    matcher2 = pattern2.match ("\\Q" + deviceReq + "\\E" )

                    if ( matcher2 == None): 
                        #print("Device is not allowed")
                        return CapabilityVerifierCode.DEVICE_NOT_PERMITTED.name

            return CapabilityVerifierCode.AUTHORIZED.name
            
        except Exception as e:
            print("CapabilityEvaluator.actionIsPermitted(ERROR): Could not validate capability token")
            print(e)
            return CapabilityVerifierCode.TOKEN_NOT_VALID.name

    # This function returns if the time frame is valid for the token.
    # @param ii_field
    # @param nb_field
    # @param na_field
    # @return

    def isTimeValid(self, ii_field, nb_field, na_field) :

        try:
            currentTime = int(round(time.time()))

            return (ii_field < currentTime) and (ii_field <= nb_field) and (nb_field < na_field) and (ii_field < na_field) and (na_field >= currentTime);
        except Exception as e:
            print("CapabilityEvaluator.isTimeValid(ERROR): Could not validate time")
            print(e)
            return False

    # This function evaluates if the signature is valid
    # @param capability_token
    # @return

    def isSignatureValid(self, capability_token) :

        try:

            #print(base64.b64decode("MEUCIQD+sKqcRFgpxp7lYswZ3jFu8ALVcniKZiYP0siTxa5GIwIgSjKXY4DeZqYSxaEjIGeRd5ei7bvKtSydT0K6MhfO3KM="))
            #print(str.encode("CapabilityToken{id=dnc7dlq0vsb66j57h6knka1n0b, ii=1673460939, is=capabilitymanager@odins.es, su=usuariouno, de=https://localhost:1027, ar=; SimpleAccessRight{ac=GET, re=/ngsi-ld/v1/entities/?type=http://example.org/vehicle/Vehicle}, nb=1673461939, na=1673471939}"))
            #print(self.certificate.public_key().verify(
            #base64.b64decode("MEUCIQD+sKqcRFgpxp7lYswZ3jFu8ALVcniKZiYP0siTxa5GIwIgSjKXY4DeZqYSxaEjIGeRd5ei7bvKtSydT0K6MhfO3KM="),
            #str.encode("CapabilityToken{id=dnc7dlq0vsb66j57h6knka1n0b, ii=1673460939, is=capabilitymanager@odins.es, su=usuariouno, de=https://localhost:1027, ar=; SimpleAccessRight{ac=GET, re=/ngsi-ld/v1/entities/?type=http://example.org/vehicle/Vehicle}, nb=1673461939, na=1673471939}"),
            #ec.ECDSA(hashes.SHA256())))

            dataToCheck = self.getSignatureString(capability_token)
            dataToCheckEncode = str.encode(dataToCheck)

            #https://cryptography.io/en/latest/hazmat/primitives/asymmetric/ec/#cryptography.hazmat.primitives.asymmetric.ec.EllipticCurvePublicKey
            #https://cryptography.io/en/latest/hazmat/primitives/asymmetric/ec/
            if (self.certificate.public_key().verify(base64.b64decode(str(capability_token.getSig())), dataToCheckEncode, ec.ECDSA(hashes.SHA256())) == None):
                return True
            
            # JAVA VERSION TRY TO VERIFY USING THESE TO ALGORITHMS.
            #   - signatures2Check.add("SHA1withRSA");
            #   - signatures2Check.add("SHA256withECDSA");
            # UNTIL THIS MOMENT ONLY "SHA256withECDSA" IS SUPPORTE BY THIS PYTHON VERSION (ec.ECDSA(hashes.SHA256()))

            return False

        except Exception as e:
            print("CapabilityEvaluator.isSignatureValid(ERROR): Could not validate signature")
            print(e)
            return False

    def getstatusoutput(self,command):

        try:

            process = Popen(command, stdout=PIPE,stderr=PIPE)
            out, err = process.communicate()

            #print(out)
            #print(err)

            return (process.returncode, out)
        
        except Exception as e:
            print("CapabilityEvaluator.getstatusoutput(ERROR): Could not obtain data")
            print(e)
            return (1, None)


    def isSignatureValid_CPABE(self, capability_token) :

        #print("*****isSignatureValid: ")

        dataToCheck = self.getSignatureString(capability_token)
        #print("isSignatureValid.dataToCheck: " + dataToCheck);

        try:

            changeDir = False
            try:
                if (len(pyCapabilityLib_folderpath) > 0):
                    os.chdir(pyCapabilityLib_folderpath)
                    changeDir = True

                codeValue, outValue = self.getstatusoutput(["java", "-jar", "cpabe_decipher.jar", str(capability_token.getSig())])

                if (changeDir == True):
                    os.chdir('../')

            except Exception as e:
                print("CapabilityEvaluator.isSignatureValid_CPABE(ERROR-4): Could not validate the signature")
                if (changeDir == True):
                    os.chdir('../')
                return False

            if(codeValue == 0):
                #Assign decipher attribute value

                #print(dataToCheck)
                if (outValue.decode('utf8') == dataToCheck):
                    #print("isSignatureValid_CPABE completed: Signature is valid.")
                    return True
                else:
                    print("CapabilityEvaluator.isSignatureValid_CPABE(ERROR-2): Signature is not valid.")
                    return False

            else:
                print("CapabilityEvaluator.isSignatureValid_CPABE(ERROR-1): Could not validate the signature")
                return False

        except Exception as e:
            print("CapabilityEvaluator.isSignatureValid_CPABE(ERROR): Could not validate the signature")
            print(e)
            return False

    # This function returns a String with the content to check that the signature is correct.
    # @param capability_token
    # @return String
    
    def getSignatureString(self, capability_token) :

        try:
        
            dataToCheck = "CapabilityToken{id=" + str(capability_token.getID()) + ", ii=" + str(capability_token.getIssueIns()) + ", is=" + str(capability_token.getIss()) + ", su=" + str(capability_token.getSub()) + ", de=" + str(capability_token.getDe()) + ", ar="

            ar = capability_token.getAr()

            arString = ""

            for i in range(len(ar)):

                if len(arString) == 0:
                    arString = "{ac=" + str(ar[i]["ac"]) + ", re=" + str(ar[i]["re"]) + "}"
                else:
                    arString = arString + ", {ac=" + str(ar[i]["ac"]) + ", re=" + str(ar[i]["re"]) + "}"

            dataToCheck = dataToCheck + "[" + arString + "]"

            dataToCheck = dataToCheck + ", nb=" + str(capability_token.getNb()) + ", na=" + str(capability_token.getNa()) + "}"

            return str(dataToCheck)
        
        except Exception as e:
            print("CapabilityEvaluator.getSignatureString(ERROR): Could not getSignature String")
            print(e)
            return None
