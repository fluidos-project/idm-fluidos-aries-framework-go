from .ElementoXACML import ElementoXACML

class ElementoPolicySetCombinerParameters(ElementoXACML):
    TIPO_POLICYSETCOMBINERPARAMETERS = "PolicySetCombinerParameters"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_POLICYSETCOMBINERPARAMETERS)
        super().setAtributos(ht)

    def getID(self) -> str:
        return super().getAtributos()['PolicySetIdRef']
    
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