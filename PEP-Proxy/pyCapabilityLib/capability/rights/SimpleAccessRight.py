#
#Copyright Odin Solutions S.L. All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

# This is an abridged version of the Access Rights. 
# Flags and Conditions are omitted.
class SimpleAccessRight :
    # Indicates the Action to perform.
    #  Examples: POST, PUT, GET, DELETE
    ac = None
    # Indicates the Resource to access
    # coap(s)://myIoTdevice/myresource
    # 
    # In this case the resource would be <myresource>
    re = None

    def  getPermittedAction(self) :
        return str(self.ac)

    def setPermittedAction(self, ac) :
        self.ac = str(ac)

    def  getResource(self) :
        return str(self.re)

    def setResource(self, re) :
        self.re = str(re)

    def  toString(self) :
        #return str("; SimpleAccessRight{ac=" + str(self.ac) + ", re=" + str(self.re) + "}")
        return str("[{ac=" + str(self.ac) + ", re=" + str(self.re) + "}]")

    #  @see smartie.utils.MemoryRemovable#free()
    def free(self) :
        self.ac = None
        self.re = None