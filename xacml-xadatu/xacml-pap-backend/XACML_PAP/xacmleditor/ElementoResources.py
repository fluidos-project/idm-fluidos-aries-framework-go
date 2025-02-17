from .ElementoXACML import ElementoXACML

class ElementoResources(ElementoXACML):
    TIPO_RESOURCES = "Resources"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_RESOURCES)
        super().setAtributos(ht)

    def getID(self) -> str:
        return ""
    
    def isUnico(self) -> bool:
        return True

    def getAllowedChild(self) -> list:
        return ["AnyResource", "Resource"]
    
    def getAllObligatory(self) -> None:
        return None
    
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string