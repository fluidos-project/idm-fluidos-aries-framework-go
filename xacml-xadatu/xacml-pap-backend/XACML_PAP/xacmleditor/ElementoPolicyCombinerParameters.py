from .ElementoXACML import ElementoXACML

class ElementoPolicyCombinerParameters(ElementoXACML):
    TIPO_POLICYCOMBINERPARAMETERS = "PolicyCombinerParameters"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_POLICYCOMBINERPARAMETERS)
        super().setAtributos(ht)

    def getID(self) -> str:
        return super().getAtributos()['PolicyIdRef']
    
    def getAllowedChild(self) -> list:
        return ["CombinerParameter"]
    
    def getAllObligatory(self) -> None:
        return None
    
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string