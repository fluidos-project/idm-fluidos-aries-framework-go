from .ElementoXACML import ElementoXACML
from .ElementoMatch import ElementoMatch

class ElementoSubjectMatch(ElementoMatch):
    TIPO_SUBJECTMATCH = "SubjectMatch"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_SUBJECTMATCH)
        super().setAtributos(ht)
    
    def getAllowedChild(self) -> None:
        return None
    
    def getAllObligatory(self) -> list:
        return ["AttributeValue", "SubjectAttributeDesignator"]
    
    def getPosicion(self, e : ElementoXACML) -> int:
        match e.getTipo():
            case "AttributeValue":
                return 1
            case "SubjectAttributeDesignator":
                return 2
            case _:
                return super().getPosicion(e)