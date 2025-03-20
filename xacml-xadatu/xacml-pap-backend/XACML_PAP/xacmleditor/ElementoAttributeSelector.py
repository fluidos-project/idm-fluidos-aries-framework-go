from .ElementoXACML import ElementoXACML

class ElementoAttributeSelector(ElementoXACML):
    TIPO_ATTRIBUTESELECTOR = "AttributeSelector"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_ATTRIBUTESELECTOR)
        super().setAtributos(ht)

    def getID(self) -> str:
        return ""
    
    def isUnico(self) -> bool:
        return True
    
    def getDataType(self) -> str:
        return super().getAtributos()['DataType']
    
    def setDataType(self, dt : str):
        atts = super().getAtributos()
        atts['DataType'] = dt
    
    def getAllowedChild(self) -> None:
        return None
    
    def getAllObligatory(self) -> None:
        return None
    
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string