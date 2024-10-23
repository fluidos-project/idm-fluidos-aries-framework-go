from .ElementoXACML import ElementoXACML

class ElementoApply(ElementoXACML):
    TIPO_APPLY = "Apply"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_APPLY)
        super().setAtributos(ht)

    def getID(self) -> str:
        return super().getAtributos()['FunctionId']
    
    def getAllowedChild(self) -> list:
        return [
        "AttributeValue",
        "Apply",
        "Function",
        "VariableReference",
        "AttributeSelector",
        "ActionAttributeDesignator",
        "ResourceAttributeDesignator",
        "SubjectAttributeDesignator",
        "EnvironmentAttributeDesignator"
        ]
    
    def getAllObligatory(self) -> None:
        return None