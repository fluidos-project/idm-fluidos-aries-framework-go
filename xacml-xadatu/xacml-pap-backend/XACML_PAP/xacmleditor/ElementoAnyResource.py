from .ElementoXACML import ElementoXACML

class ElementoAnyResource(ElementoXACML):
    TIPO_ANY_RESOURCE = "AnyResource"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_ANY_RESOURCE)
        super().setAtributos(ht)

    def getID(self) -> str:
        return ""

    def getAllowedChild(self) -> None:
        return None
    
    def getAllObligatory(self) -> None:
        return None
    
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string