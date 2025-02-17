from .ElementoXACML import ElementoXACML

class ElementoPolicySet(ElementoXACML):
    TIPO_POLICYSET = "PolicySet"

    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_POLICYSET)
        super().setAtributos(ht)

    def getID(self) -> str:
        return super().getAtributos()["PolicySetId"]
    
    def getAllowedChild(self) -> list:
        return [
        "Description",
        "PolicySetDefaults",
        "Target",
        "PolicySet",
        "Policy",
        "PolicySetIdReference",
        "PolicyIdReference",
        "Obligations",
        "CombinerParameters",
        "PolicySetCombinerParameters",
        "PolicyCombinerParameters"
        ]
    
    def getPosicion(self, e : ElementoXACML) -> int:
        match e.getTipo():
            case "Description":
                return 1
            case "PolicySetDefaults":
                return 2
            case "Target":
                return 3
            case ("Policy", "PolicySet", "PolicySetIdReference", 
                "PolicyIdReference", "CombinerParameters", "PolicySetCombinerParameters", 
                "PolicyCombinerParameters"):
                return 4
            case "Obligations":
                return 5
            case _:
                return super().getPosicion(e)
        
    def getAllObligatory(self) -> list:
        return ["Target"]
    
    def getMaxNumChild(self, e : ElementoXACML) -> int:
        match e.getTipo():
            case ("Description", "PolicySetDefaults", "Target", "PolicySetCombinerParameters", 
                  "PolicyCombinerParameters", "Obligations", "CombinerParameters"):
                return 1
            case _:
                return super().getMaxNumChild(e)