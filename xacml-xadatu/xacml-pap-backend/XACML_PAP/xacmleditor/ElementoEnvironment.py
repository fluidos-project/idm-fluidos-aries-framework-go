from .ElementoXACML import ElementoXACML

class ElementoEnvironment(ElementoXACML):
    TIPO_ENVIRONMENT = "Environment"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_ENVIRONMENT)
        super().setAtributos(ht)

    def getID(self) -> str:
        return ""
    
    def getAllowedChild(self) -> list:
        return ["EnvironmentMatch"]
    
    def getAllObligatory(self) -> list:
        return ["EnvironmentMatch"]
    
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string