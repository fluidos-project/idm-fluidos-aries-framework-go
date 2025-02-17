#
#Copyright Odin Solutions S.L. All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

from .CapabilityToken import *


class CapabilityTokenDelegated(CapabilityToken) :

    del = None

    # This extended version of the Capability Token adds a new value, "del"
    # to express that this token has been delegated from the original owner
    # of the token to another entity for their use.
    #  
    # This has several implications, 
    # 	- The owner has to expedite the delegated token.
    #  - 
    
    def __init__(self, del, id, ii, iss, su, de, si, ar, nb, na):
        super().__init__(del, id, ii, iss, su, de, si, ar, nb, na)
        self.del = str(del)

    def  getDel(self) :
        return str(self.del)

    def setDel(self, del) :
        self.del = str(del)