from .ElementoXACML import ElementoXACML

class ElementoAttributeDesignator(ElementoXACML):
    def getID(self) -> str:
        return super().getAtributos()['AttributeId']
    
    def setDataType(self, dt : str):
        atts = super().getAtributos()
        atts['DataType'] = dt

    def isUnico(self) -> bool:
        return True
    
    def getDataType(self) -> str:
        return super().getAtributos()['DataType']
    
    def getAllowedChild(self) -> None:
        return None
    
    def getAllObligatory(self) -> None:
        return None