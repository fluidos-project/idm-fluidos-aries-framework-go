#
#Copyright Odin Solutions S.L. All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

from ..token.CapabilityToken import *
from .CapabilityEvaluator import *
from .CapabilityVerifierCode import *

class SimpleCapabilityEvaluator :

    def validateCapabilityToken (CACertificatePATH, action, resource, device, subject, capabilityToken, headers, body) :

        try:

            cev = CapabilityEvaluator(CACertificatePATH)
            ct  = CapabilityToken.getCapabilityTokenFromString(capabilityToken)

            code = cev.evaluateCapabilityToken(ct, action,resource,device,subject,headers,body)
            
            result = "CODE: " + str(code)
            
            return str(result)

        except Exception as e:
            print("SimpleCapabilityEvaluator.validateCapabilityToken(ERROR): Could not validate capability token")
            print(e)
            return "CODE: " + str(CapabilityVerifierCode.TOKEN_NOT_VALID.name)