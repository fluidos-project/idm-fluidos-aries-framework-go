#
#Copyright Odin Solutions S.L. All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

from SimpleAccessRight import *

class AccessRight(SimpleAccessRight) :
	
	# Indicates with a flag if the all conditions have to be evaluated positively
	# to perform the solicited action ( AND with the integer 0) or if the case of any
	# of conditions are met the action can be performed (OR with the integer 1)

	f = None

	# A list of Condition to be met.
	
	co = None 

	def __init__(self):
        super().__init__()

	
	def  getFlag(self) :
        return int(self.f)

	def setFlag(self, f) :
        self.f = str(f)
	
	def  getConditions(self) :
        return self.co

	def setConditions(self, co) :
        self.co = co
