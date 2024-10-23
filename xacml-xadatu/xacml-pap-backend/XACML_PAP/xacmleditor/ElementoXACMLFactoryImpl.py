from .ElementoXACML import ElementoXACML
from .ElementoXACMLFactory import ElementoXACMLFactory
from .ElementoPolicySet import ElementoPolicySet
from .ElementoPolicySetDefaults import ElementoPolicySetDefaults
from .ElementoXPathVersion import ElementoXPathVersion
from .ElementoPolicy import ElementoPolicy
from .ElementoVariableDefinition import ElementoVariableDefinition
from .ElementoPolicyDefaults import ElementoPolicyDefaults
from .ElementoRule import ElementoRule
from .ElementoTarget import ElementoTarget
from .ElementoSubjects import ElementoSubjects
from .ElementoActions import ElementoActions
from .ElementoResources import ElementoResources
from .ElementoAnySubject import ElementoAnySubject
from .ElementoAnyAction import ElementoAnyAction
from .ElementoAnyResource import ElementoAnyResource
from .ElementoEnvironments import ElementoEnvironments
from .ElementoAttributeValue import ElementoAttributeValue
from .ElementoAction import ElementoAction
from .ElementoActionMatch import ElementoActionMatch
from .ElementoActionAttributeDesignator import ElementoActionAttributeDesignator
from .ElementoAttributeSelector import ElementoAttributeSelector
from .ElementoSubject import ElementoSubject
from .ElementoSubjectMatch import ElementoSubjectMatch
from .ElementoSubjectAttributeDesignator import ElementoSubjectAttributeDesignator
from .ElementoResource import ElementoResource
from .ElementoResourceMatch import ElementoResourceMatch
from .ElementoResourceAttributeDesignator import ElementoResourceAttributeDesignator
from .ElementoEnvironment import ElementoEnvironment
from .ElementoEnvironmentMatch import ElementoEnvironmentMatch
from .ElementoEnvironmentAttributeDesignator import ElementoEnvironmentAttributeDesignator
from .ElementoCondition import ElementoCondition
from .ElementoApply import ElementoApply
from .ElementoDescription import ElementoDescription
from .ElementoObligations import ElementoObligations
from .ElementoObligation import ElementoObligation
from .ElementoAttributeAssignment import ElementoAttributeAssignment
from .ElementoCombinerParameters import ElementoCombinerParameters
from .ElementoCombinerParameter import ElementoCombinerParameter
from .ElementoRuleCombinerParameters import ElementoRuleCombinerParameters
from .ElementoPolicyCombinerParameters import ElementoPolicyCombinerParameters
from .ElementoPolicySetCombinerParameters import ElementoPolicySetCombinerParameters
from .ElementoFunction import ElementoFunction
from .ElementoVariableReference import ElementoVariableReference
from .ElementoPolicySetIdReference import ElementoPolicySetIdReference
from .ElementoPolicyIdReference import ElementoPolicyIdReference

class ElementoXACMLFactoryImpl(ElementoXACMLFactory):

    def obtenerElementoXACML(tipo : str, atributos : dict) -> ElementoXACML:
        elem = None        
        match tipo:
            case ElementoPolicySet.TIPO_POLICYSET:
                elem = ElementoPolicySet(atributos)
            case ElementoPolicySetDefaults.TIPO_POLICYSETDEFAULTS:
                elem = ElementoPolicySetDefaults(atributos)
            case ElementoXPathVersion.TIPO_XPATHVERSION:
                elem = ElementoXPathVersion(atributos)
            case ElementoPolicy.TIPO_POLICY:
                elem = ElementoPolicy(atributos)
            case ElementoVariableDefinition.TIPO_VARIABLEDEFINITION:
                elem = ElementoVariableDefinition(atributos)
            case ElementoPolicyDefaults.TIPO_POLICYDEFAULTS:
                elem = ElementoPolicyDefaults(atributos)
            case ElementoRule.TIPO_RULE:
                elem = ElementoRule(atributos)
            case ElementoTarget.TIPO_TARGET:
                elem = ElementoTarget(atributos)
            case ElementoSubjects.TIPO_SUBJECTS:
                elem = ElementoSubjects(atributos)
            case ElementoActions.TIPO_ACTIONS:
                elem = ElementoActions(atributos)
            case ElementoResources.TIPO_RESOURCES:
                elem = ElementoResources(atributos)
            case ElementoAnySubject.TIPO_ANY_SUBJECT:
                elem = ElementoAnySubject(atributos)
            case ElementoAnyAction.TIPO_ANY_ACTION:
                elem = ElementoAnyAction(atributos)
            case ElementoAnyResource.TIPO_ANY_RESOURCE:
                elem = ElementoAnyResource(atributos)
            case ElementoEnvironments.TIPO_ENVIRONMENTS:
                elem = ElementoEnvironments(atributos)
            case ElementoAttributeValue.TIPO_ATTRIBUTEVALUE:
                elem = ElementoAttributeValue(atributos)
            case ElementoAction.TIPO_ACTION:
                elem = ElementoAction(atributos)
            case ElementoActionMatch.TIPO_ACTIONMATCH:
                elem = ElementoActionMatch(atributos)
            case ElementoActionAttributeDesignator.TIPO_ACTIONATTRIBUTEDESIGNATOR:
                elem = ElementoActionAttributeDesignator(atributos)
            case ElementoAttributeSelector.TIPO_ATTRIBUTESELECTOR:
                elem = ElementoAttributeSelector(atributos)
            case ElementoSubject.TIPO_SUBJECT:
                elem = ElementoSubject(atributos)
            case ElementoSubjectMatch.TIPO_SUBJECTMATCH:
                elem = ElementoSubjectMatch(atributos)
            case ElementoSubjectAttributeDesignator.TIPO_SUBJECTATTRIBUTEDESIGNATOR:
                elem = ElementoSubjectAttributeDesignator(atributos)
            case ElementoResource.TIPO_RESOURCE:
                elem = ElementoResource(atributos)
            case ElementoResourceMatch.TIPO_RESOURCEMATCH:
                elem = ElementoResourceMatch(atributos)
            case ElementoResourceAttributeDesignator.TIPO_RESOURCEATTRIBUTEDESIGNATOR:
                elem = ElementoResourceAttributeDesignator(atributos)
            case ElementoEnvironment.TIPO_ENVIRONMENT:
                elem = ElementoEnvironment(atributos)
            case ElementoEnvironmentMatch.TIPO_ENVIRONMENTMATCH:
                elem = ElementoEnvironmentMatch(atributos)
            case ElementoEnvironmentAttributeDesignator.TIPO_ENVIRONMENTATTRIBUTEDESIGNATOR:
                elem = ElementoEnvironmentAttributeDesignator(atributos)
            case ElementoCondition.TIPO_CONDITION:
                elem = ElementoCondition(atributos)
            case ElementoApply.TIPO_APPLY:
                elem = ElementoApply(atributos)
            case ElementoDescription.TIPO_DESCRIPTION:
                elem = ElementoDescription(atributos)
            case ElementoObligations.TIPO_OBLIGATIONS:
                elem = ElementoObligations(atributos)
            case ElementoObligation.TIPO_OBLIGATION:
                elem = ElementoObligation(atributos)
            case ElementoAttributeAssignment.TIPO_ATTRIBUTEASSIGNMENT:
                elem = ElementoAttributeAssignment(atributos)
            case ElementoCombinerParameters.TIPO_COMBINERPARAMETERS:
                elem = ElementoCombinerParameters(atributos)
            case ElementoCombinerParameter.TIPO_COMBINERPARAMETER:
                elem = ElementoCombinerParameter(atributos)
            case ElementoRuleCombinerParameters.TIPO_RULECOMBINERPARAMETERS:
                elem = ElementoRuleCombinerParameters(atributos)
            case ElementoPolicyCombinerParameters.TIPO_POLICYCOMBINERPARAMETERS:
                elem = ElementoPolicyCombinerParameters(atributos)
            case ElementoPolicySetCombinerParameters.TIPO_POLICYSETCOMBINERPARAMETERS:
                elem = ElementoPolicySetCombinerParameters(atributos)
            case ElementoFunction.TIPO_FUNCTION:
                elem = ElementoFunction(atributos)
            case ElementoVariableReference.TIPO_VARIABLEREFERENCE:
                elem = ElementoVariableReference(atributos)
            case ElementoPolicySetIdReference.TIPO_POLICYSETIDREFERENCE:
                elem = ElementoPolicySetIdReference(atributos)
            case ElementoPolicyIdReference.TIPO_POLICYIDREFERENCE:
                elem = ElementoPolicyIdReference(atributos)

        return elem
    
    def getInstancia():
        return ElementoXACMLFactoryImpl()
