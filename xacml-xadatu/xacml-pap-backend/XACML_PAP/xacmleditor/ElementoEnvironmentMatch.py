from .ElementoXACML import ElementoXACML
from .ElementoMatch import ElementoMatch

class ElementoEnvironmentMatch(ElementoMatch):
    TIPO_ENVIRONMENTMATCH = "EnvironmentMatch"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_ENVIRONMENTMATCH)
        super().setAtributos(ht)
    
    def getAllowedChild(self) -> None:
        return None
    
    def getAllObligatory(self) -> list:
        return ["AttributeValue", "EnvironmentAttributeDesignator"]
    
    def getPosicion(self, e : ElementoXACML) -> int:
        match e.getTipo():
            case "AttributeValue":
                return 1
            case "EnvironmentAttributeDesignator":
                return 2
            case _:
                return super().getPosicion(e)