#
#Copyright Odin Solutions S.L. All Rights Reserved.
#
#SPDX-License-Identifier: Apache-2.0
#

# **
# * What does this class do?
# * 
# * This class depends from what class or component?
# * What classes or components depend from this class?
# * 
# * 
# * This class Generates a 
# * 
# * @author dgarcia@odins.es
# * 
# * Modified by:
# * 
# * @author jasanchez@odins.es
# *
class CapabilityTokenRequest:
    subjectType = None
    subject = None
    device = None
    action = None
    resource = None
    # *
    # 	 * 
    # 	 * This constructor generates the object with that contains the strings passed as parameters. 
    # 	 * We can call it directly, or use the getCapabilityTokenRequest function that given a JSON
    # 	 *
    # 	 * @param subjectType
    # 	 * @param subject
    # 	 * @param action
    # 	 * @param device
    # 	 * @param resource
    def __init__(self, subjectType,  subject,  action,  device,  resource) :
        self.subjectType = str(subjectType)
        self.subject = str(subject)
        self.action = str(action)
        self.device = str(device)
        self.resource = str(resource)
    # *
    # 	 * 
    # 	 * Given the function parameter, a String containing a JSON document that contains the subject, device, action and resource,
    # 	 * we return a CapabilityTokenRequest object with said information. 
    # 	 * 
    # 	 * 
    # 	 * Idea of the JSON document
    # 	 * 
    # 	 * {
    # 	 *  "st" : "Type of the Subject -- (The format of this field is subject to change)",
    # 	 * 	"su" : "Identity of the Subject -- (The format of this field is subject to change)",
    # 	 * 	"de" : "Generally, the URL of the device -- e.g, coaps://CAFE:DCAF:8080",
    # 	 *  "ac" : "Generally, a CRUD command, POST, PUT, GET, DELETE, etc.",
    # 	 *  "re" : "Generally, the LOCATION path or URI path or the resource"
    # 	 * }
    # 	 * 
    # 	 * Example of JSON document passed through  capabilitytokenrequest
    # 	 * 
    # 	 * {
    # 	 *  "st" : "urn:ietf:params:scim:schemas:core:2.0:id",
    # 	 * 	"su" : "(aBVL}}7X7@Gv?Z*",
    # 	 * 	"de" : "coaps://CAFE:DCAF:8080",
    # 	 *  "ac" : "POST",
    # 	 *  "re" : "temperature"
    # 	 * }
    # 	 * 
    # 	 * This would translate to a 
    # 	 * 
    # 	 * 		POST coaps://CAFE:DCAF:8080/temperature 
    # 	 * 
    # 	 * that is authorized by the capability token.
    # 	 * 
    # 	 * @param String capabilitytokenrequest
    # 	 * @return CapabilityTokenRequest
    # // ???DEPRECATED???
    # public static CapabilityTokenRequest getCapabilityTokenRequestFromJSON (String capabilitytokenrequest)
    # {
    # 	JsonParser parser = new JsonParser();
    # 	JsonElement jelement = parser.parse(capabilitytokenrequest);
    # 	JsonObject token_json = jelement.getAsJsonObject();
    #
    # 	JsonElement ts_json = token_json.get("ts");
    # 	String ts = ts_json.getAsString();
    #
    # 	JsonElement id_json = token_json.get("su");
    # 	String su = id_json.getAsString();
    #
    # 	JsonElement ii_json = token_json.get("de");
    # 	String de = ii_json.getAsString();
    #
    # 	JsonElement is_json = token_json.get("ac");
    # 	String ac = is_json.getAsString();
    #
    # 	JsonElement su_json = token_json.get("re");
    # 	String re = su_json.getAsString();
    #
    # 	return new CapabilityTokenRequest(ts, su, ac, de, re);
    # }

    # Getters
    def  getSubjectType(self) :
        return str(self.subjectType)

    def  getSubject(self) :
        return str(self.subject)

    def  getDevice(self) :
        return str(self.device)

    def  getAction(self) :
        return str(self.action)

    def  getResource(self) :
        return str(self.resource)

    # Setters
    def setResource(self, resource) :
        self.resource = str(resource)

    def setDevice(self, device) :
        self.device = str(device)

    def setAction(self, action) :
        self.action = str(action)

    def setSubject(self, subject) :
        self.subject = str(subject)

    def setSubjectType(self, subjectType) :
        self.subject = str(subjectType)

    # ToString
    # ???DEPRECATED???
    def  toString(self) :
        try:
            request = "{\"su\": \"" + self.getSubject() + "\",\"de\":\"" + self.getDevice() + "\",\"ac\":\"" + self.getAction() + "\",\"re\":\"" + self.getResource() + "\"}"
            return request
        except Exception as e:
            print("CapabilityTokenRequest.toString(ERROR): Could not obtain toString value")
            print(e)
            return None