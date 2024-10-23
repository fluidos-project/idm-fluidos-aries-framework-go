from .ElementoXACML import ElementoXACML

class ElementoPolicyDefaults(ElementoXACML):
    TIPO_POLICYDEFAULTS = "PolicyDefaults"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_POLICYDEFAULTS)
        super().setAtributos(ht)

    def getID(self) -> str:
        return ""
    
    def getAllowedChild(self) -> list:
        return ["XPathVersion"]
    
    def getAllObligatory(self) -> None:
        return None
    
    def getMaxNumChild(self, elementoXACML : ElementoXACML) -> int:
        if elementoXACML.getTipo() == "XPathVersion":
            return 1
        else:
            return super().getMaxNumChild(elementoXACML)
    
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string