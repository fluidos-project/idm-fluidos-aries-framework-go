from .ElementoXACML import ElementoXACML

class ElementoSubject(ElementoXACML):
    TIPO_SUBJECT = "Subject"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_SUBJECT)
        super().setAtributos(ht)

    def getID(self) -> str:
        return ""
    
    def getAllowedChild(self) -> list:
        return ["SubjectMatch"]
    
    def getAllObligatory(self) -> list:
        return ["SubjectMatch"]
    
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string