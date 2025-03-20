from .ElementoXACML import ElementoXACML

class ElementoFunction(ElementoXACML):
    TIPO_FUNCTION = "Function"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_FUNCTION)
        super().setAtributos(ht)

    def getID(self) -> str:
        return super().getAtributos()['FunctionId']
    
    def getAllowedChild(self) -> None:
        return None
    
    def getAllObligatory(self) -> None:
        return None