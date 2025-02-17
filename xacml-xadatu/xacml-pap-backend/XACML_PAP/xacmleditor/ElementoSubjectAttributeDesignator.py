from .ElementoAttributeDesignator import ElementoAttributeDesignator

class ElementoSubjectAttributeDesignator(ElementoAttributeDesignator):
    TIPO_SUBJECTATTRIBUTEDESIGNATOR = "SubjectAttributeDesignator"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_SUBJECTATTRIBUTEDESIGNATOR)
        super().setAtributos(ht)