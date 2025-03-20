from .ElementoXACML import ElementoXACML

class ElementoAttributeAssignment(ElementoXACML):
    TIPO_ATTRIBUTEASSIGNMENT = "AttributeAssignment"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_ATTRIBUTEASSIGNMENT)
        super().setAtributos(ht)

    def getID(self) -> str:
        return super().getAtributos()['AttributeId']
    
    def getDataType(self) -> str:
        return super().getAtributos()['DataType']
    
    def setDataType(self, dt : str):
        atts = super().getAtributos()
        atts['DataType'] = dt
    
    def getAllowedChild(self) -> None:
        return None
    
    def getAllObligatory(self) -> None:
        return None
