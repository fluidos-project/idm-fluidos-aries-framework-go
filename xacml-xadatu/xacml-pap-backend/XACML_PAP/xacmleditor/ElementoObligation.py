from .ElementoXACML import ElementoXACML

class ElementoObligation(ElementoXACML):
    TIPO_OBLIGATION = "Obligation"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_OBLIGATION)
        super().setAtributos(ht)

    def getID(self) -> str:
        return super().getAtributos()['ObligationId']
    
    def getAllowedChild(self) -> list:
        return ["AttributeAssignment"]
    
    def getAllObligatory(self) -> None:
        return None