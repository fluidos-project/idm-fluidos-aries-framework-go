from .ElementoXACML import ElementoXACML
from .ElementoMatch import ElementoMatch

class ElementoActionMatch(ElementoMatch):
    TIPO_ACTIONMATCH = "ActionMatch"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_ACTIONMATCH)
        super().setAtributos(ht)
    
    def getAllowedChild(self) -> None:
        return None
    
    def getAllObligatory(self) -> list:
        return ["AttributeValue", 'ActionAttributeDesignator']
    
    def getPosicion(self, e : ElementoXACML) -> int:
        match e.getTipo():
            case "AttributeValue":
                return 1
            case "ActionAttributeDesignator":
                return 2
            case _:
                return super().getPosicion(e)