from .ElementoXACML import ElementoXACML

class ElementoPolicy(ElementoXACML):
    TIPO_POLICY = "Policy"

    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_POLICY)
        super().setAtributos(ht)

    def getID(self) -> str:
        return super().getAtributos()['PolicyId']
    
    def getAllowedChild(self) -> list:
        return [
        "Description",
        "PolicyDefaults",
        "Target",
        "Rule",
        "VariableDefinition",
        "Obligations",
        "CombinerParameters",
        "RuleCombinerParameters"
        ]
        
    def getAllObligatory(self) -> list:
        return ['Target']
    
    def getMaxNumChild(self, e: ElementoXACML) -> int:
        match e.getTipo():
            case "Target":
                return 1
            case "Description":
                return 1
            case "CombinerParameters":
                return 1
            case "RuleCombinerParameters":
                return 1
            case "PolicyDefaults":
                return 1
            case "Obligations":
                return 1
            case _:
                return super().getMaxNumChild(e)
    
    def getPosicion(self, e: ElementoXACML) -> int:
        match e.getTipo():
            case "Description":
                return 1
            case "PolicySetDefaults":
                return 2
            case "Target":
                return 3
            case "CombinerParameters" | "RuleCombinerParameters" | "VariableDefinition" | "Rule":
                return 4
            case "Obligations":
                return 5
            case _:
                return super().getPosicion(e)