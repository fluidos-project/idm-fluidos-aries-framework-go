#
#Copyright Odin Solutions S.L. All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

class Condition :

    # The Condition type from :
    # https://tools.ietf.org/html/draft-li-core-conditional-observe-04#section-4
    # 
    # 	+-----------------------+------+-----------------+
    #   | Condition type (5 bit)| Id.  | Condition Value |
    #   +-----------------------+------+-----------------+
    #   | Cancellation          |  0   |       no        |
    #   +-----------------------+------+-----------------+
    #   | Time series           |  1   |       no        |
    #   +-----------------------+------+-----------------+
    #   | Minimum response time |  2   |       yes       |
    #   +-----------------------+------+-----------------+
    #   | Maximum response time |  3   |       yes       |
    #   +-----------------------+------+-----------------+
    #   | Step                  |  4   |       yes       |
    #   +-----------------------+------+-----------------+
    #   | AllValues<            |  5   |       yes       |
    #   +-----------------------+------+-----------------+
    #   | AllValues>            |  6   |       yes       |
    #   +-----------------------+------+-----------------+
    #   | Value=                |  7   |       yes       |
    #   +-----------------------+------+-----------------+
    #   | Value<>               |  8   |       yes       |
    #   +-----------------------+------+-----------------+
    #   | Periodic              |  9   |       yes       |
    #   +-----------------------+------+-----------------+

	t = None

    # The value of the condition.
    # 
    # Example: If we were to retrieve the temperature of a sensor, we could
    # express a range of values between which gathering this information is permitted. 
    v = None

    # The Unit of the of measure of the measure that will be represented by the resource.
    # Example: In the case of temperature ?C
    u = None

    # t The type according to draft-li-core-conditional-observe-04#section-4
    # v The value that has to be evaluated by the condition
    # u The unit of measurement of the value (v)

    def __init__(self, t,  v,  u) :
        self.t = int(t)
        self.v = int(v)
        self.setU(u)

    def  getType(self) :
        return int(self.t)

    def setType(self, t) :
        self.t = int(t)

    def  getValue(self) :
        return int(self.v)

    def setValue(self, v) :
        self.v = int(v)

    def  getU(self) :
        return str(self.u)

    def setU(self, u) :
        self.u = str(u)