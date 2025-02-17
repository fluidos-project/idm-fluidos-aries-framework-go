from .ElementoXACML import ElementoXACML

class ElementoAction(ElementoXACML):
    TIPO_ACTION = "Action"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_ACTION)
        super().setAtributos(ht)

    def getID(self) -> str:
        return ""
    
    def getAllowedChild(self) -> list:
        return ["ActionMatch"]
    
    def getAllObligatory(self) -> list:
        return ["ActionMatch"]
    
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string