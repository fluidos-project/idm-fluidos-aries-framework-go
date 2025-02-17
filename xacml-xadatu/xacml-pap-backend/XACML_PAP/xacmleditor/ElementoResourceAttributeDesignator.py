from .ElementoAttributeDesignator import ElementoAttributeDesignator

class ElementoResourceAttributeDesignator(ElementoAttributeDesignator):
    TIPO_RESOURCEATTRIBUTEDESIGNATOR = "ResourceAttributeDesignator"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_RESOURCEATTRIBUTEDESIGNATOR)
        super().setAtributos(ht)