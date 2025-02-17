from xml.sax import ErrorHandler, SAXParseException

class SAXErrorHandler(ErrorHandler):
    '''Error handler del AnalizadorSAX'''
    __valid = True
    __errores = ""

    def isValid(self) -> bool:
        '''Devuelve si es válido o no'''
        return self.__valid

    def getErrores(self) -> str:
        '''Devuelve los errores obtenidos'''
        return self.__errores

    def setErrores(self, err : str) -> None:
        '''Establece los errores obtenidos'''
        self.__errores = err

    def reset(self) -> None:
        '''Resetea valid a True'''
        self.__valid = True

    def warning(self, exception: SAXParseException) -> None:
        '''Establece un warning nuevo'''
        self.__errores += f"Warning: {exception.getMessage()}\n"
        self.__errores += f" at line {exception.getLineNumber()}, column {exception.getColumnNumber()}\n"
        self.__valid = False

    def error(self, exception: SAXParseException) -> None:
        '''Establece un error nuevo'''
        self.__errores += f"Error: {exception.getMessage()}\n"
        self.__errores += f" at line {exception.getLineNumber()}, column {exception.getColumnNumber()}\n"
        self.__valid = False

    def fatalError(self, exception: SAXParseException) -> None:
        '''Establece un error fatal nuevo'''
        self.__errores += f"Fatal Error: {exception.getMessage()}\n"
        self.__errores += f" at line {exception.getLineNumber()}, column {exception.getColumnNumber()}\n"