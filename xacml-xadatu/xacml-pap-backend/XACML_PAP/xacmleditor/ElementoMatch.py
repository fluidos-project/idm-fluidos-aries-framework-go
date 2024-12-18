from abc import ABC
from .ElementoXACML import ElementoXACML

class ElementoMatch(ElementoXACML, ABC):

    def getID(self) -> str:
        return super().getAtributos()['MatchId']