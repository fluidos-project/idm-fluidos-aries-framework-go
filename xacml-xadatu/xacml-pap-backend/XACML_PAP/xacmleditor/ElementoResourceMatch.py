from .ElementoXACML import ElementoXACML
from .ElementoMatch import ElementoMatch

class ElementoResourceMatch(ElementoMatch):
    TIPO_RESOURCEMATCH = "ResourceMatch"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_RESOURCEMATCH)
        super().setAtributos(ht)
    
    def getAllowedChild(self) -> None:
        return None
    
    def getAllObligatory(self) -> list:
        return ["AttributeValue", "ResourceAttributeDesignator"]
    
    def getPosicion(self, e : ElementoXACML) -> int:
        match e.getTipo():
            case "AttributeValue":
                return 1
            case "ResourceAttributeDesignator":
                return 2
            case _:
                return super().getPosicion(e)