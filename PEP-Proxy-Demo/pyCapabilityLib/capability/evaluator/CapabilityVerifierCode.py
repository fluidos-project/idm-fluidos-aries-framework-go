#
#Copyright Odin Solutions S.L. All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

from enum import Enum


class CapabilityVerifierCode(Enum) :
    AUTHORIZED = 0
    TOKEN_NOT_VALID = 1
    ACTION_NOT_PERMITTED = 2
    RESOURCE_NOT_PERMITTED = 3
    DEVICE_NOT_PERMITTED = 4
    SIGNATURE_NOT_VALID = 5
    OUTATIME = 6