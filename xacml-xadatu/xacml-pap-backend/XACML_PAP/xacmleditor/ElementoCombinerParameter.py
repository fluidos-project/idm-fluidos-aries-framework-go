from .ElementoXACML import ElementoXACML

class ElementoCombinerParameter(ElementoXACML):
    TIPO_COMBINERPARAMETER = "CombinerParameter"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_COMBINERPARAMETER)
        super().setAtributos(ht)

    def getID(self) -> str:
        return super().getAtributos()['ParameterName']
    
    def getAllowedChild(self) -> list:
        return ["AttributeValue"]
    
    def getAllObligatory(self) -> list:
        return ["AttributeValue"]
    
    def getMaxNumChild(self, elementoXACML : ElementoXACML) -> int:
        match elementoXACML.getTipo():
            case "AttributeValue":
                return 1
            case _:
                return super().getMaxNumChild(elementoXACML)