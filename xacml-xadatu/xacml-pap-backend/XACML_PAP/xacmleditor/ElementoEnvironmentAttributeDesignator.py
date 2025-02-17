from .ElementoAttributeDesignator import ElementoAttributeDesignator

class ElementoEnvironmentAttributeDesignator(ElementoAttributeDesignator):
    TIPO_ENVIRONMENTATTRIBUTEDESIGNATOR = "EnvironmentAttributeDesignator"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_ENVIRONMENTATTRIBUTEDESIGNATOR)
        super().setAtributos(ht)