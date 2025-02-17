from .ElementoXACML import ElementoXACML

class ElementoRuleCombinerParameters(ElementoXACML):
    TIPO_RULECOMBINERPARAMETERS = "RuleCombinerParameters"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_RULECOMBINERPARAMETERS)
        super().setAtributos(ht)

    def getID(self) -> str:
        return super().getAtributos()['RuleIdRef']
    
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