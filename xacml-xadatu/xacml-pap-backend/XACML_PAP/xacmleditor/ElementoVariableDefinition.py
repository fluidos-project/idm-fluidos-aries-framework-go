from .ElementoXACML import ElementoXACML

class ElementoVariableDefinition(ElementoXACML):
    TIPO_VARIABLEDEFINITION = "VariableDefinition"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_VARIABLEDEFINITION)
        super().setAtributos(ht)

    def getID(self) -> str:
        return super().getAtributos()['VariableId']
    
    def getAllowedChild(self) -> None:
        return None
    
    def getAllObligatory(self) -> list:
        return ["Apply"]