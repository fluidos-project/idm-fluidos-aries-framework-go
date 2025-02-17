from .ElementoXACML import ElementoXACML

class ElementoEnvironments(ElementoXACML):
    TIPO_ENVIRONMENTS = "Environments"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_ENVIRONMENTS)
        super().setAtributos(ht)

    def getID(self) -> str:
        return ""
    
    def isUnico(self) -> bool:
        return True

    def getAllowedChild(self) -> list:
        return ["Environment"]
    
    def getAllObligatory(self) -> list:
        return ['Environment']
    
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string