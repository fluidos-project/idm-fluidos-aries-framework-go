from .ElementoAttributeDesignator import ElementoAttributeDesignator

class ElementoActionAttributeDesignator(ElementoAttributeDesignator):
    TIPO_ACTIONATTRIBUTEDESIGNATOR = "ActionAttributeDesignator"
    
    def __init__(self, ht):
        super().setTipo(self.TIPO_ACTIONATTRIBUTEDESIGNATOR)
        super().setAtributos(ht)