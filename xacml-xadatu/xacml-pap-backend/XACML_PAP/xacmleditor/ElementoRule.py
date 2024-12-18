from .ElementoXACML import ElementoXACML

class ElementoRule(ElementoXACML):
    TIPO_RULE = "Rule"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_RULE)
        super().setAtributos(ht)

    def getID(self) -> str:
        return super().getAtributos()['RuleId']
    
    def getAllowedChild(self) -> list:
        return [
        "Target",
        "Condition",
        "Description"
        ]
    
    def getAllObligatory(self) -> None:
        return None
    
    def isAllowedChild(self, elementoXACML : ElementoXACML) -> bool:
        match elementoXACML.getTipo():
            case "Target", "Condition":
                return True
            case _:
                return super().isAllowedChild(elementoXACML)
            
    def getPosicion(self, elementoXACML : ElementoXACML) -> int:
        match elementoXACML.getTipo():
            case "Description":
                return 1
            case "Target":
                return 2
            case "Condition":
                return 3
            case _:
                return super().getPosicion(elementoXACML)
            
    def getMaxNumChild(self, elementoXACML : ElementoXACML) -> int:
        match elementoXACML.getTipo():
            case "Description", "Target", "Condition":
                return 1
            case _:
                return super().getMaxNumChild(elementoXACML)