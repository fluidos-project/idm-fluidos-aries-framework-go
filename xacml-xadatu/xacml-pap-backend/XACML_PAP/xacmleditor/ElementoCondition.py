from .ElementoXACML import ElementoXACML

class ElementoCondition(ElementoXACML):
    TIPO_CONDITION = "Condition"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_CONDITION)
        super().setAtributos(ht)

    def getID(self) -> str:
        return ""
    
    def getAllowedChild(self) -> None:
        return None
    
    def getAllObligatory(self) -> list:
        return ["Apply"]
    
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string