#
#Copyright Odin Solutions S.L. All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

import time
import json

# Objects holding the fields necessary to form a token
class CapabilityToken :

    id = None
    ii = None
    iss = None
    su = None
    de = None
    si = None
    ar = None
    nb = None
    na = None

    parser = None
    gson = None
    listType = None
    
    # Constructor
     
    # id			ID
    # ii			Issuing timestamp
    # iss			Issuer
    # su			Subject
    # de			Device
    # si			Signature
    # ar			Access resources
    # nb			Not before
    # na			Not after
    def __init__(self, id, ii, iss, su, de, si, ar, nb, na) :
        self.id = str(id)
        self.ii = ii
        self.iss = str(iss)
        self.su = str(su)
        self.de = str(de)
        self.si = str(si)
        self.ar = ar
        self.nb = nb
        self.na = na

    # Try to parse a token
    # 
    # capability_token	The String received
    # return			The token (null if the format is incorrect)

    def getCapabilityTokenFromString(capability_token):
        
        try:

            if (capability_token == None):
                return None

            token_json = json.loads(capability_token)

            id = str(token_json["id"])
            ii = token_json["ii" ]
            iss = str(token_json["is"])
            su = str(token_json["su"])
            de = str(token_json["de"])
            si = str(token_json["si"])

            ar = token_json["ar"]

            nb = token_json["nb"]
            na = token_json["na"]

            if ((id == None) or (iss == None) or (su == None) or (de == None) or (si == None) or (ar == None) or (nb < ii) or (nb >= na)):
                return None
            
            return CapabilityToken(id, ii, iss, su, de, si, ar, nb, na)
        
        except Exception as e:
            print("CapabilityToken.getCapabilityTokenFromString(ERROR): Could not load capability token")
            print(e)
            return None

    # Getters
    def  getID(self) :
        return str(self.id)

    def  getIssueIns(self) :
        return self.ii
    
    def  getIss(self) :
        return str(self.iss)

    def  getSub(self) :
        return str(self.su)

    def  getDe(self) :
        return str(self.de)
    
    def  getSig(self) :
        return str(self.si)

    def  getAr(self) :
        return self.ar
    
    def  getNb(self) :
        return self.nb
    
    def  getNa(self) :
        return self.na


    def IsTimeValid(self) :
        try:

            ii_field = self.getIssueIns();
            nb_field = self.getNb();
            na_field = self.getNa();
            currentTime = int(round(time.time()))

            return ((ii_field < currentTime) and (ii_field <= nb_field) and (nb_field < na_field) and (ii_field < na_field) and (na_field >= currentTime))
        
        except Exception as e:
            print("CapabilityToken.IsTimeValid(ERROR): Could not validate")
            print(e)
            return False

    # Setters
    def setSig(self, sig) :
        self.si = str(sig)