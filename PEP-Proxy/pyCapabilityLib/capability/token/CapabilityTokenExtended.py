#
#Copyright Odin Solutions S.L. All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

from .CapabilityToken import *
import json

class CapabilityTokenExtended(CapabilityToken) :

    rd = None

    ac = None

    re = None
    
    def __init__(self, id, ii, iss, su, de, si, ar, nb, na, rd):
        super().__init__(id, ii, iss, su, de, si, ar, nb, na)
        self.rd = str(rd)
        

    def getRd(self) :
        return str(self.rd)

    def setRd(self, rd) :
        self.rd = str(rd)

    def setAcRe(self, ac, re) :
        self.ac = str(ac)
        self.re = str(re)

    def  getCTSigned(self):

        try:
            return	json.dumps({
                    "id": super().getID(),
                    "ii": super().getIssueIns(),
                    "is": super().getIss(),
                    "su": super().getSub(),
                    "de": super().getDe(),
                    "si": super().getSig(),
                    "ar": [{"ac":str(self.ac), "re":str(self.re)}],
                    "nb": super().getNb(),
                    "na": super().getNa()
                }, indent=4)
        except Exception as e:
            print("CapabilityTokenExtended.getCTSigned(ERROR): Could not load capability token")
            print(e)
            return None
