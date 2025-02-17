from .ElementoXACML import ElementoXACML

class ElementoVariableReference(ElementoXACML):
    TIPO_VARIABLEREFERENCE = "VariableReference"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_VARIABLEREFERENCE)
        super().setAtributos(ht)

    def getID(self) -> str:
        return super().getAtributos()['VariableId']
    
    def getAllowedChild(self) -> None:
        return None
    
    def getAllObligatory(self) -> list:
        return ['Apply']
