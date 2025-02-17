from .ElementoXACML import ElementoXACML

class ElementoSubjects(ElementoXACML):
    TIPO_SUBJECTS = "Subjects"
    
    def __init__(self, ht : dict):
        super().setTipo(self.TIPO_SUBJECTS)
        super().setAtributos(ht)

    def getID(self) -> str:
        return ""
    
    def isUnico(self) -> bool:
        return True

    def getAllowedChild(self) -> list:
        return ["AnySubject", "Subject"]
    
    def getAllObligatory(self) -> None:
        return None
    
    def __str__(self) -> str:
        string = f"<{self.getTipo()}"
        if self.esVacio():
            string += "/>"
        else:
            string += ">"
        return string