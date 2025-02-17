from .ElementoXACML import ElementoXACML

class ElementoPolicySetDefaults(ElementoXACML):
    TIPO_POLICYSETDEFAULTS = "PolicySetDefaults"

    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_POLICYSETDEFAULTS)
        super().setAtributos(ht)

    def getID(self) -> str:
        return ""
    
    def getAllowedChild(self) -> list:
        return ["XPathVersion"]
        
    def getAllObligatory(self) -> None:
        return None
    
    def getMaxNumChild(self, e) -> int:
        if e.getTipo() == "XPathVersion":
            return 1
        else:
            return super().getMaxNumChild(e)
        
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string