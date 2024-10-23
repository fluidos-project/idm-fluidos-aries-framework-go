<template>

  <div id="app" class="container-fluid">
    <div class="row">
      <div id= mainPAP class="col-lg-12">
        Policy Administration Point
      </div>
      <div id=mainPAP-central class="col-lg-12">
        Policy Administration Point
      </div>
    </div>
    <div class="row" id="MainPAP-menu" v-if="showIndex">
      <div class="col-xl-4" id="MainPAP-menu-domain">
        <div id="mainPAP-domain"> Domain </div>
        <label for="select1" style="margin-top: 5%; text-align: left;">
          <h5>Select domain</h5>
        </label>
        <select id="select1" v-model="selectedDomain" class="form-select">
          <option v-for="domain in domains" :key="domain" :value="domain">{{ domain }}</option>
        </select>
        <button id="button2" @click="showForm" style="margin-top: 5%;">
          Create new domain
        </button>
        <div v-if="showForms">
          <label for="type-domain-name" style="margin-top: 5%; text-align: left;">
            <h5>Domain name</h5>
          </label>
          <input type="text" class="form-control" id="type-domain-name" v-model="domainForm">
          <button id="button3" @click="addDomain" style="margin-top: 5%;">
            Submit
          </button>
        </div>
      </div>
      <div class="col-xl-4" id="MainPAP-menu-administration">
        <div id="mainPAP-administration"> Administration </div>
        <button id="button1" style="margin-top: 15%;" @click="indexToPolicies">
          <img src="./assets/manage_policies.png" height=25 width=25> Manage Policies
        </button>
        <button id="button2" style="margin-bottom: 25%;" @click="indexToAttributes">
          <img src="./assets/manage_attributes.png" height=25 width=25> Manage Attributes
        </button>
      </div>
      <div class="col-xl-4" id="MainPAP-menu-configuration">
        <div id="mainPAP-configuration"> Configuration </div>
        <label for="policies-storage-path" style="margin-top: 5%; text-align: left;">
          <h5>Policies storage path</h5>
        </label>
        <input type="text" class="form-control" id="policies-storage-path" value="./XACML_PAP/PAPConfigData/Policies/" readonly>
        <label for="attributes-storage-path" style="margin-top: 5%; text-align: left;">
          <h5>Attributes storage path</h5>
        </label>
        <input type="text" class="form-control" id="attributes-storage-path" value="./XACML_PAP/PAPConfigData/XACMLAtts/" readonly>
      </div>
    </div>
    <div class="row" v-if="showIndex">
      <div class="col-lg-12">
        <button id="button4" @click="indexToDomain">
          <img src="./assets/exit.png" height=25 width=25> Exit
        </button>
      </div>
    </div>
    <div class="row" id="mainManagePolicies-menu" v-if="showPolicies">
      <div id="policiesManagement" class="col-lg-12" style="font-size: 12px;border: rgb(44, 113, 148) 0.5px solid;height: 25px;">
        <p>Policies Management | Domain: {{ selectedDomain }}</p>
      </div>
    </div>
    <div class="row" v-if="showPolicies">
      <div id="A-B" class="col-lg-12">
        <div class="row">
          <div id="A" class="col-lg-2" style="background: rgb(172, 213, 246); border:rgb(44, 113, 148) 1.0px solid;">
            <div id="A.1" class="col" style="margin-bottom: 2%; margin-top: 1%; border: rgb(44, 113, 148) 1.0px solid;">
              <div style="font-size:12px; background: rgb(172, 213, 246);font-family: Arial, Helvetica, sans-serif;color:rgb(21, 27, 107); font-weight: bold;">
                Policies
              </div>
              <div>
                <div style="background: white; height: 250px;overflow: scroll;">
                  <table id="table-policies" class="table table-striped table-sm table-hover">
                    <tbody style="text-align: justify;font-size: small;">
                      <!-- 1. Politicas extraidas de data -->
                      <!-- <tr v-for="policy in policies" :key="policy.id" @click="selectedpolicy = policy" v-bind:class="{'red': policy.id == selectedpolicy.id}"> -->
                        <!-- <td>{{ policy.name }}</td> -->
                      <!-- </tr> -->
                      <!-- 2. Nueva politica agregada -->
                      <tr @click="selectedPolicy = newPolicy" v-for="newPolicy in newPolicies" :key="newPolicy.id" :value="newPolicy">
                        <td>{{ newPolicy.name }}</td>
                      </tr>
                      <!-- 3. Politicas extraidas del API con el método obtainPolicies(), cómo se llama al método @click ="obtainPolicies()"? -->
                      <tr @click="selectPolicy(Policy)" v-for="Policy in policies" :key="Policy.PolicyId.value" :value="Policy" :class="{'highlight': Policy.PolicyId.value == selectedPolicy.PolicyId.value}">
                        <!-- Cómo filtrar las politicas? -->
                        <td>{{ Policy.PolicyId.value }}</td>
                      </tr>
                      <!--Cuadro de dialogo para avisar de que hay cambios en la regla que se perderan -->
                      <div class="modal-backdrop" v-if="showModal15">
                        <div class="modal" role="dialog" id="modal15">
                          <form id="formElement-selectPolicy" style="text-align: justify;">
                            <label>
                              <h5>If you continue, the rule changes will be lost</h5>
                            </label>
                            <div class="form-group" style="margin-top: 15px">
                              <button class="btn" @click="checksChanges = false; showModal15 = false; selectPolicy(provisionalPolicy)">Continue</button>
                              <button class="btn" @click="showModal15 = false">Cancel</button>
                            </div>
                          </form>
                        </div>
                      </div>
                    </tbody>
                  </table>
                </div>
              </div>
              <div style="background-color: whitesmoke ; border: rgb(44, 113, 148) 1.0px solid; padding-top: 2px; padding-bottom: 2px;">
                <button @click="showModal1 = true" style="font-size:12px; border: black 0.5px solid; width: 27%; margin-right: 2px;" data-toggle="button" aria-pressed="false" autocomplete="off">
                  <img src="./assets/new.png" height=25 width=25> New
                </button>
                <!--Cuadro de dialogo para agregar una nueva politica -->
                <div class="modal-backdrop" v-if="showModal1">
                  <div class="modal" role="dialog" id="modal1">
                    <form id="formElement-newPolicy" style="text-align: justify;">
                      <label for="type-policy-name">
                        <h5>New Policy</h5>
                      </label>
                      <div class="form-group" style="margin-bottom: 5px;">
                        <label for="type-policy-name">Type Policy Name</label>
                      </div>
                      <div class="form-group" style="margin-bottom: 5px;">
                        <input type="text" class="form-control" id="type-policy-name" v-model="policyModal"> <!--entrada donde se toma la nueva accion -->
                      </div>
                      <div class="form-group" style="margin-top: 15px">
                        <button class="btn" @click="addPolicy">OK</button> <!--metodo para agregar un nuevo recurso -->
                        <button class="btn" @click="showModal1 = false">Close</button>
                      </div>
                    </form>
                  </div>
                </div>
                <button style="font-size:12px; border: black 0.5px solid; width: 25%; margin-right: 2px;" @click="deletePolicy">
                  <img src="./assets/delete.png" height=25 width=25> Del
                </button>
                <button style="font-size:12px; border: black 0.5px solid; width: 40%;" @click="showModal3 = true">
                  <img src="./assets/rename.png" height=25 width=25> Rename
                </button>
                <!--Cuadro de dialogo para renombrar una politica -->
                <div class="modal-backdrop" v-if="showModal3">
                  <div class="modal" role="dialog" id="modal3">
                    <form id="formElement-renamePolicy" style="text-align: justify;">
                      <label for="type-policy-name">
                        <h5>Rename Policy</h5>
                      </label>
                      <div class="form-group" style="margin-bottom: 5px;">
                        <label for="type-policy-name">Type New Policy Name</label>
                      </div>
                      <div class="form-group" style="margin-bottom: 5px;">
                        <input type="text" class="form-control" id="type-policy-name" v-model="policyModal"> <!--entrada donde se toma la nueva accion -->
                      </div>
                      <div class="form-group" style="margin-top: 15px">
                        <button class="btn" @click="renamePolicy">OK</button> <!--metodo para agregar un nuevo recurso -->
                        <button class="btn" @click="showModal3 = false">Close</button>
                      </div>
                    </form>
                  </div>
                </div>
              </div>
            </div>
            <div id="A.2" class="col" style="margin-bottom: 1%; margin-top: 1%;border: rgb(44, 113, 148) 1.0px solid;">
              <div style="font-size:12px; background: rgb(172, 213, 246);font-family: Arial, Helvetica, sans-serif;color:rgb(21, 27, 107); font-weight: bold;">
                Rules
              </div>
              <div>
                <div style="background: white;height: 150px;overflow: scroll;">
                  <table id="table-rules" class="table table-striped table-sm table-hover">
                    <tbody style="text-align: justify;font-size: small;">
                      <!--1. Reglas extraidos de data -->
                      <!-- <tr v-for="rule in rules" :key="rule.id" @click="selectedrule = rule" v-bind:class="{'red': rule.id == selectedrule.id}"> -->
                        <!-- <td>{{ rule.name }}</td> -->
                      <!-- </tr> -->
                      <!--2. Nuevo regla agregada -->
                      <tr @click="selectedRule = newRule" v-for="newRule in newRules" :key="newRule.name" :value="newRule">
                        <td>{{ newRule.name }}</td>
                      </tr>
                      <!--3.Reglas extraidas del API con el método obtainPolicies(), como se llama al método @click ="obtainPolicies()"? -->
                      <tr @click="selectRule(Rule)" v-for="Rule in selectedPolicy.Rules" :key="Rule.RuleId.value" :value="Rule" :class="{'highlight': Rule.RuleId.value == selectedRule.RuleId.value}">
                        <!--Como filtrar las reglas?-->
                        <td>{{ Rule.RuleId.value }}</td>
                      </tr>
                      <!--Cuadro de dialogo para avisar de que hay cambios en la regla que se perderan -->
                      <div class="modal-backdrop" v-if="showModal14">
                        <div class="modal" role="dialog" id="modal14">
                          <form id="formElement-selectRule" style="text-align: justify;">
                            <label>
                              <h5>If you continue, the rule changes will be lost</h5>
                            </label>
                            <div class="form-group" style="margin-top: 15px">
                              <button class="btn" @click="checksChanges = false; showModal14 = false; selectRule(provisionalRule)">Continue</button>
                              <button class="btn" @click="showModal14 = false">Cancel</button>
                            </div>
                          </form>
                        </div>
                      </div>
                    </tbody>
                  </table>
                </div>
              </div>
              <div style="background-color: whitesmoke ; border: rgb(44, 113, 148) 1.5px solid; padding-bottom: 2px;">
                <button @click="showModal2 = true" style="font-size:12px; border: black 0.5px solid; width:35%; margin-right: 2px;">
                  <img src="./assets/new.png" height=10 width=10> New
                </button>
                <!--Cuadro de dialogo para agregar una nueva regla -->
                <div class="modal-backdrop" v-if="showModal2">
                  <div class="modal" role="dialog" id="modal2">
                    <form id="formElement-newRule" style="text-align: justify;">
                      <label for="type-rule-name">
                        <h5>New Rule</h5>
                      </label>
                      <div class="form-group" style="margin-bottom: 5px;">
                        <label for="type-rule-name">Type Rule Name</label>
                      </div>
                      <div class="form-group" style="margin-bottom: 5px;">
                        <input type="text" class="form-control" id="type-rule-name" v-model="ruleModal"> <!--entrada donde se toma la nueva regla -->
                      </div>
                      <div class="form-group" style="margin-top: 15px">
                        <button class="btn" @click="addRule">OK</button>   <!--metodo para agregar una nueva regla -->
                        <button class="btn" @click="showModal2 = false">Close</button>
                      </div>
                    </form>
                  </div>
                </div>
                <button style="font-size:12px; border: black 0.5px solid; width: 35%; margin-right: 2px;" @click="deleteRule">
                  <img src="./assets/delete.png" height=10 width=10> Del
                </button>
                <button style="font-size:12px; border: black 0.5px solid; width: 20%;">
                  <img src="./assets/flecha-arriba.png" height=15 width=15>
                </button>
                <button style="font-size:12px; border: black 0.5px solid; width: 72%; margin-right: 2px;" @click="savePolicies">
                  <img src="./assets/apply.png" height=12 width=12> Apply
                </button>
                <!--Cuadro de dialogo para confirmar que las politicas han sido guardadas con exito -->
                <div class="modal-backdrop" v-if="showModal12">
                  <div class="modal" role="dialog" id="modal12">
                    <form id="formElement-savePolicies" style="text-align: justify;">
                      <label>
                        <h5>Operation carried out successfully</h5>
                      </label>
                      <div class="form-group" style="margin-top: 15px">
                        <button class="btn" @click="showModal12 = false">Close</button>
                      </div>
                    </form>
                  </div>
                </div>
                <button style="font-size:12px; border: black 0.5px solid; width: 20%;">
                  <img src="./assets/flecha-hacia-abajo.png" height=15 width=15>
                </button>
              </div>
            </div>
          </div>
          <div id="B" class="col-lg-10" style="border: rgb(44, 113, 148) 1.0px solid;background: rgb(172, 213, 246);">
            <div id="B.1" class="col" style="margin-bottom: 1%; margin-top: 1%">
              <div class="row">
                <div id="B.1.1" class="col">
                  <div style="border: rgb(44, 113, 148) 1.0px solid; font-size:12px; background: rgb(172, 213, 246);font-family: Arial, Helvetica, sans-serif;color:rgb(21, 27, 107); font-weight: bold;">
                    Resources
                  </div>
                  <div style="border: rgb(44, 113, 148) 1.0px solid; background: white; height: 310px;font-size: 12px;overflow: scroll;">
                    <div>
                      <div sytle="width 80%;align-items: center;justify-content: space-between;">
                        <div style="background: white;height: 150px;">
                          <table id="table-resources-selected" class="table table-striped table-sm">
                            <tbody style="text-align:left; font-size: small;">
                              <!--Recursos extraidos de data -->
                              <tr v-for="resource in resources" :key="resource.name.value">
                                <td>
                                  <input :id="'check_' + resource.name.value" v-model="checksResources[resource.name.value]" type="checkbox" :disabled="checkResources || disableAllChecks" @click="checksChanges = true"> {{ resource.name.value }} &lt;()&gt; {{ resource.xacml_id.value }}
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div style="min-height: 1%;"></div>
                  <div style="border: rgb(44, 113, 148) 1.0px solid; background: white;font-family: Arial, Helvetica, sans-serif;color:black ;height: 9%;">
                    <div>
                      <input id="check2" v-model="checkResources" type="checkbox" :disabled="disableAllChecks" @click="checkAllResources"> All
                    </div>
                    <!--seleccionar todos los recursos  v-model="selectAll"-->
                  </div>
                </div>
                <div id="B.1.2" class="col">
                  <div style="border: rgb(44, 113, 148) 1.0px solid; font-size:12px; background: rgb(172, 213, 246);font-family: Arial, Helvetica, sans-serif;color:rgb(21, 27, 107); font-weight: bold;">
                    Subjects
                  </div>
                  <div style="border: rgb(44, 113, 148) 1.0px solid; background: white; height: 310px;font-size: 12px;overflow: scroll;">
                    <div >
                      <div sytle="width 80%;align-items: center;justify-content: space-between;">
                        <div style="background: white;height: 150px;">
                          <table id="table-subjects-selected" class="table table-striped table-sm">
                            <tbody style="text-align: justify;font-size: small;">
                              <!--Sujetos extraidos de data -->
                              <tr v-for="subject in subjects" :key="subject.name.value">
                                <td>
                                  <input :id="'check_' + subject.name.value" v-model="checksSubjects[subject.name.value]" type="checkbox" :disabled="checkSubjects || disableAllChecks" @click="checksChanges = true"> {{ subject.name.value }} &lt;()&gt; {{ subject.xacml_id.value }}
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div style="min-height: 1%;"></div>
                  <div style="border: rgb(44, 113, 148) 1.0px solid; background: white;font-family: Arial, Helvetica, sans-serif;color:black ;height: 9%;">
                    <div>
                      <input id="check4" v-model="checkSubjects" type="checkbox" :disabled="disableAllChecks" @click="checkAllSubjects"> All
                    </div>
                  </div>
                </div>
                <div id="B.1.3" class="col">
                  <div style="border: rgb(44, 113, 148) 1.0px solid; font-size:12px; background: rgb(172, 213, 246);font-family: Arial, Helvetica, sans-serif;color:rgb(21, 27, 107); font-weight: bold;">
                    Actions
                  </div>
                  <div style="border: rgb(44, 113, 148) 1.0px solid; background: white; height: 310px;font-size: 12px;overflow: scroll;">
                    <div>
                      <div sytle="width 80%;align-items: center;justify-content: space-between;">
                        <div style="background: white;height: 150px;">
                          <table id="table-actions-selected" class="table table-striped table-sm">
                            <tbody style="text-align: justify;font-size: small;">
                              <!--Acciones extraidos de data -->
                              <tr v-for="action in actions" :key="action.name.value">
                                <td>
                                  <input :id="'check_' + action.name.value" v-model="checksActions[action.name.value]" type="checkbox" :disabled="checkActions || disableAllChecks" @click="checksChanges = true"> {{ action.name.value }} &lt;()&gt; {{ action.xacml_id.value }}
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div style="min-height: 1%;"></div>
                  <div style="border: rgb(44, 113, 148) 1.0px solid; background: white;font-family: Arial, Helvetica, sans-serif;color:black ;height: 9%;">
                    <div>
                      <input id="check6" v-model="checkActions" type="checkbox" :disabled="disableAllChecks" @click="checkAllActions"> All
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div id="B.2" class="col" style="margin-bottom: 0.5px;margin-top: 0.5px; margin-left: 0% ;border: rgb(44, 113, 148) 1px solid;background: white;">
              <div class="row">
                <div class="col" style="width: 50px;">
                  <label for="select2" style="font-size: 12px; width: 150%; height: 80%; padding-top: 10px;margin-right: 20px;">Rule CA</label>
                </div>
                <div style="width: 500px; height: 38px;" >
                  <select id="select2" class="form-select" style="font-size:13px; height: 80%" :disabled="disableAllChecks">
                    <option value="rule-CA">urn:oasis:names:tc:xacml:1.0:rule-combinig-algorithm:first-applicable</option>
                  </select>
                </div>
                <div style="width: 34%; background:rgb(172, 213, 246);"></div>
                <div style="width: 80px">
                  <label for="select3" style="font-size: 12px; width: 80%; height: 70%; padding-top: 10px;">Rule</label>
                </div>
                <div class="col" style="width: 20%; height: 80%;">
                  <select id="select3" class="form-select" style="font-size:13px; height: 100%" v-model="selectedEffect" :disabled="disableAllChecks">
                    <option value="Permit">Permit</option>
                    <option value="Deny">Deny</option>
                  </select>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div id="C" class="col-lg-12" style="border: rgb(44, 113, 148) 1px solid;">
        <button @click="policiesToIndex" style="font-size:12px; border: black 0.5px solid; width: 150px; padding-top:25px;padding-bottom:25px;margin-top: 28px;margin-bottom: 34px; margin-right: 30px;">
          <img src="./assets/back.png" height=25 width=25> Back
        </button>
        <!--Cuadro de dialogo para avisar de que hay cambios que se perderan -->
        <div class="modal-backdrop" v-if="showModal13">
          <div class="modal" role="dialog" id="modal13">
            <form id="formElement-policiesToIndex" style="text-align: justify;">
              <label>
                <h5>If you continue, the changes made will be lost</h5>
              </label>
              <div class="form-group" style="margin-top: 15px">
                <button class="btn" @click="changes = false; checksChanges = false; showModal13 = false; policiesToIndex()">Continue</button>
                <button class="btn" @click="showModal13 = false">Cancel</button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
    <div class="row" id="mainManageAttributes-menu" v-if="showAttributes">
      <div id="policiesManagement" class="col-lg-12" style="font-size: 12px;border: rgb(44, 113, 148) 0.5px solid;height: 25px;">
        <p>Attributes Management | Domain: {{ selectedDomain }}</p>
      </div>
    </div>
    <div class="row" style="min-width: 50px;" v-if="showAttributes">
      <div id="A-B" class="col-lg-12">
        <div class="row">
          <div id="B" class="col-lg-12" style="border: rgb(44, 113, 148) 1.0px solid;background: rgb(172, 213, 246); min-width: 500px;">
            <div id="B.1" class="col" style="margin-bottom: 1%; margin-top: 1% ">
              <div class="row">
                <div id="B.1.1" class="col-lg-4">
                  <div style="border: rgb(44, 113, 148) 1.0px solid; font-size:12px; background: rgb(172, 213, 246);font-family: Arial, Helvetica, sans-serif;color:rgb(21, 27, 107); font-weight: bold;">
                    Resources
                  </div>
                  <div style="border: rgb(44, 113, 148) 1.0px solid; background: white; height: 310px;font-size: 12px;overflow: scroll;">
                    <div id="getAttributes-resources">
                      <div class="container">
                        <table id="table-resources" class="table table-sm table table-hover">
                          <tbody style="text-align: justify;">
                            <!-- 1. Recursos extraidos de data -->
                            <!-- <tr v-for="resource in resources" :key="resource.id" @click="selectedresource = resource" v-bind:class="{'red': resource.id == selectedresource.id}"> -->
                              <!-- <td style="overflow: hidden; text-overflow: ellipsis;"> -->
                                <!-- {{ resource.name }} -->
                              <!-- </td> -->
                            <!-- </tr> -->
                            <!--2. Nuevo recurso agregado -->
                            <tr @click="selectedResource = newResource" v-for="newResource in newResources" :key="newResource.name.value" :value="newResource">
                              <td>
                                {{ newResource.name.value }} &lt;()&gt; {{ newResource.xacml_id.value }}
                              </td>
                            </tr>
                            <!-- 3. Recursos extraidos del API con el método obtainAttributes(), llama al método @click ="obtainAttributes()" -->
                            <tr @click="selectedResource = resource" v-for="resource in resources" :key="resource.name.value" :value="resource" :class="{'highlight': resource.name.value == selectedResource.name.value}">
                              <td> <!-- filtrado de los recursos -->
                                {{ resource.name.value }} &lt;()&gt; {{ resource.xacml_id.value }}
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    </div>
                  </div>
                  <div style="min-height: 1%;"></div>
                  <div style="background-color: white ; border: rgb(44, 113, 148) 1px solid;padding-bottom: 5px;height:9%; ">
                    <button @click="showModal4 = true" style="font-size:12px; border: black 0.5px solid;width:25%;margin-bottom: 0.5px;">
                      <img src="./assets/new.png" height=10 width=10> New Resource
                    </button>
                    <!--Cuadro de dialogo para agregar un nuevo recurso -->
                    <div class="modal-backdrop" v-if="showModal4">
                      <div class="modal" role="dialog" id="modal4">
                        <form id="formElement-newResource" style="text-align: justify;">
                          <label for="type-resource-name">
                            <h5>New Resource</h5>
                          </label>
                          <div class="form-group" style="margin-bottom: 5px;">
                            <label for="type-resource-name">Type Resource Name</label>
                          </div>
                          <div class="form-group" style="margin-bottom: 5px;">
                            <input type="text" class="form-control" id="type-resource-name" v-model="resourceModal"> <!--entrada donde se toma la nueva accion -->
                          </div>
                          <div class="form-group" style="margin-bottom: 5px;">
                            <label for="type-resource-id">
                              Type Resource ID
                            </label>
                          </div>
                          <div class="form-group" style="margin-bottom: 5px;">
                            <select type="text" class="form-control" id="type-resource-id">
                              <option>
                                urn:oasis:names:tc:xacml:1.0:resource:resource-id
                              </option>
                            </select>
                          </div>
                          <div class="form-group" style="margin-top: 15px">
                            <button class="btn" @click="addResource">OK</button> <!--metodo para agregar un nuevo recurso -->
                            <button class="btn" @click="closeModal4">Close</button>
                          </div>
                        </form>
                      </div>
                    </div>
                    <button style="font-size:12px; border: black 0.5px solid;width: 25%;" @click="deleteResource">
                      <img src="./assets/delete.png" height=10 width=10> Delete Resource
                    </button>
                    <!--Cuadro de dialogo para avisar de que no se puede eliminar un recurso porque ninguno ha sido seleccionado -->
                    <div class="modal-backdrop" v-if="showModal8">
                      <div class="modal" role="dialog" id="modal8">
                        <form id="formElement-deleteResource" style="text-align: justify;">
                          <label>
                            <h5>No resource selected</h5>
                          </label>
                          <div class="form-group" style="margin-top: 15px">
                            <button class="btn" @click="showModal8 = false">Close</button>
                          </div>
                        </form>
                      </div>
                    </div>
                  </div>
                </div>
                <div id="B.1.2" class="col-lg-4">
                  <div style="border: rgb(44, 113, 148) 1.0px solid; font-size:12px; background: rgb(172, 213, 246);font-family: Arial, Helvetica, sans-serif;color:rgb(21, 27, 107); font-weight: bold; ">
                    Actions
                  </div>
                  <div style="border: rgb(44, 113, 148) 1.0px solid; background: white; height: 310px;font-size: 12px;max-width: 100%;overflow: scroll;">
                    <div id="actions" style=" max-width: 100%;">
                      <div class="container" style=" max-width: 100%;">
                        <table id="table-actions" class="table table-sm table-hover" style=" max-width: 100%;">
                          <tbody style="text-align: justify;">
                            <!-- 1. Acciones extraidas de data -->
                            <!-- <tr v-for="action in actions" @click="selectedaction = action" v-bind:class="{'red': action.id == selectedaction.id}" v-bind:key="action.id"> -->
                              <!-- <td> -->
                                <!-- {{ action.name.value }} -->
                              <!-- </td> -->
                            <!-- </tr> -->
                            <!-- 2. Nuevo acción agregada -->
                            <tr @click="selectedAction = newAction" v-for="newAction in newActions" :key="newAction.name.value" :value="newAction">
                              <td>
                                {{ newAction.name.value }} &lt;()&gt; {{ newAction.xacml_id.value }}
                              </td>
                            </tr>
                            <!-- 3. Acciones extraidas del API con el método obtainAttributes() -->
                            <tr @click="selectedAction = action" v-for="action in actions" :key="action.name.value" :value="action" :class="{'highlight': action.name.value == selectedAction.name.value}">
                              <td> <!--Como filtrar las acciones?-->
                                {{ action.name.value }} &lt;()&gt; {{ action.xacml_id.value }}
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    </div>
                  </div>
                  <div style="min-height: 1%;"></div>
                  <div style="background-color: white ; border: rgb(44, 113, 148) 1px solid;padding-bottom: 5px; ">
                    <button @click="showModal5 = true" style="font-size:12px; border: black 0.5px solid;width:25%;">
                      <img src="./assets/new.png" height=10 width=10> New Action
                    </button>
                    <!--Cuadro de dialogo para agregar una nueva acción -->
                    <div class="modal-backdrop" v-if="showModal5">
                      <div class="modal" role="dialog" id="modal5">
                        <form id="formElement-newAction" style="text-align: justify;">
                          <label for="type-action-name">
                            <h5>New Action</h5>
                          </label>
                          <div class="form-group" style="margin-bottom: 5px;">
                            <label for="type-action-name">Type Action Name</label>
                          </div>
                          <div class="form-group" style="margin-bottom: 5px;">
                            <input type="text" class="form-control" id="type-action-name" v-model="actionModal"> <!--entrada donde se toma la nueva accion -->
                          </div>
                          <div class="form-group" style="margin-bottom: 5px;">
                            <label for="type-action-id">
                              Type Action ID
                            </label>
                          </div>
                          <div class="form-group" style="margin-bottom: 5px;">
                            <select type="text" class="form-control" id="type-action-id">
                              <option>
                                urn:oasis:names:tc:xacml:1.0:action:action-id
                              </option>
                            </select>
                          </div>
                          <div class="form-group" style="margin-top: 15px">
                            <button class="btn" @click="addAction">OK</button> <!--metodo para agregar un nuevo acción -->
                            <button class="btn" @click="closeModal5">Close</button>
                          </div>
                        </form>
                      </div>
                    </div>
                    <button style="font-size:12px; border: black 0.5px solid;width: 25%;" @click="deleteAction">
                      <img src="./assets/delete.png" height=10 width=10> Delete Action
                    </button>
                    <!--Cuadro de dialogo para avisar de que no se puede eliminar una accion porque ninguna ha sido seleccionada -->
                    <div class="modal-backdrop" v-if="showModal9">
                      <div class="modal" role="dialog" id="modal9">
                        <form id="formElement-deleteAction" style="text-align: justify;">
                          <label>
                            <h5>No action selected</h5>
                          </label>
                          <div class="form-group" style="margin-top: 15px">
                            <button class="btn" @click="showModal9 = false">Close</button>
                          </div>
                        </form>
                      </div>
                    </div>
                  </div>
                </div>
                <div id="B.1.3 " class="col-lg-4">
                  <div style="border: rgb(44, 113, 148) 1.0px solid; font-size:12px; background: rgb(172, 213, 246);font-family: Arial, Helvetica, sans-serif;color:rgb(21, 27, 107); font-weight: bold;">
                    Subjects
                  </div>
                  <div style="border: rgb(44, 113, 148) 1.0px solid; background: white; height: 310px;font-size: 12px;overflow: scroll;">
                    <div id="subjects">
                      <div class="container" style=" max-width: 100%;">
                        <table id="table-subjects" class="table table-hover table-sm" style=" max-width: 100%;">
                          <tbody style="text-align: justify;max-width: 100%;">
                            <!-- 1. Subjects extraidos de data -->
                            <!-- <tr> -->
                              <!-- <td v-for="subject in subjects" @click="selectedsubject = subject" v-bind:class="{'red': subject.id == selectedsubject.id}" v-bind:key="subject.id"> -->
                                <!-- {{ subject.name }} -->
                              <!-- </td> -->
                            <!-- </tr> -->
                            <!-- 2. Nuevo subject agregado -->
                            <tr @click="selectedSubject = newSubject" v-for="newSubject in newSubjects" :key="newSubject.name.value" :value="newSubject">
                              <td>
                                {{ newSubject.name.value }} &lt;()&gt; {{ newSubject.xacml_id.value }}
                              </td>
                            </tr>
                            <!-- 3. Subjects extraidos del API con el método obtainAttributes() -->
                            <tr @click="selectedSubject = subject" v-for="subject in subjects" :key="subject.name.value" :value="subject" :class="{'highlight': subject.name.value == selectedSubject.name.value}">
                              <td> <!--Como filtrar los sujetos?-->
                                {{ subject.name.value }} &lt;()&gt; {{ subject.xacml_id.value }}
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    </div>
                  </div>
                  <div style="min-height: 1%;"></div>
                  <div style="background-color: white ; border: rgb(44, 113, 148) 1px solid;padding-bottom: 5px;">
                    <button @click="showModal6 = true" style="font-size:12px; border: black 0.5px solid;width:25%;">
                      <img src="./assets/new.png" height=10 width=10> New Subject
                    </button>
                    <!--Cuadro de dialogo para agregar un nuevo sujeto -->
                    <div class="modal-backdrop" v-if="showModal6">
                      <div class="modal" role="dialog" id="modal6">
                        <form id="formElement-newSubject" style="text-align: justify;">
                          <label for="type-subject-name">
                            <h5>New Subject</h5>
                          </label>
                          <div class="form-group" style="margin-bottom: 5px;">
                            <label for="type-subject-name">Type Subject Name</label>
                          </div>
                          <div class="form-group" style="margin-bottom: 5px;">
                            <input type="text" class="form-control" id="type-subject-name" v-model="subjectModal"> <!--entrada donde se toma la nueva accion -->
                          </div>
                          <div class="form-group" style="margin-bottom: 5px;">
                            <label for="type-subject-id">
                              Type Subject ID
                            </label>
                          </div>
                          <div class="form-group" style="margin-bottom: 5px;">
                            <select type="text" class="form-control" id="type-subject-id" v-model="subjectIDModal">
                              <option v-for="subjectIdType in subjectIdTypes" :key="subjectIdType.subjectIdPAP" :value="subjectIdType.subjectIdPAP">{{ subjectIdType.subjectIdPAP }}</option>
                            </select>
                          </div>
                          <div class="form-group" style="margin-top: 15px">
                            <button class="btn" @click="addSubject">OK</button> <!--metodo para agregar un nuevo sujeto -->
                            <button class="btn" @click="closeModal6">Close</button>
                          </div>
                        </form>
                      </div>
                    </div>
                    <button style="font-size:12px; border: black 0.5px solid;width: 25%;" @click="deleteSubject">
                      <img src="./assets/delete.png" height=10 width=10> Delete Subject
                    </button>
                    <!--Cuadro de dialogo para avisar de que no se puede eliminar un sujeto porque ninguno ha sido seleccionado -->
                    <div class="modal-backdrop" v-if="showModal10">
                      <div class="modal" role="dialog" id="modal10">
                        <form id="formElement-deleteSubject" style="text-align: justify;">
                          <label>
                            <h5>No subject selected</h5>
                          </label>
                          <div class="form-group" style="margin-top: 15px">
                            <button class="btn" @click="showModal10 = false">Close</button>
                          </div>
                        </form>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div id="C" class="col-lg-12" style="border: rgb(44, 113, 148) 1px solid;">
        <button @click="saveAttributes" style="font-size:12px; border: black 0.5px solid; width: 150px; padding-top:25px;padding-bottom:25px;margin-top: 28px;margin-bottom: 34px; margin-right: 30px;">
          <img src="./assets/apply.png" height=25 width=25> Save All Attributes
        </button>
        <!--Cuadro de dialogo para confirmar que los atributos han sido guardados con exito -->
        <div class="modal-backdrop" v-if="showModal7">
          <div class="modal" role="dialog" id="modal7">
            <form id="formElement-saveAttributes" style="text-align: justify;">
              <label>
                <h5>Operation carried out successfully</h5>
              </label>
              <div class="form-group" style="margin-top: 15px">
                <button class="btn" @click="showModal7 = false">Close</button>
              </div>
            </form>
          </div>
        </div>
        <button @click="attributesToIndex" style="font-size:12px; border: black 0.5px solid; width: 150px; padding-top:25px;padding-bottom:25px;margin-top: 28px;margin-bottom: 34px; margin-right: 30px;">
          <img src="./assets/back.png" height=25 width=25> Back
        </button>
        <!--Cuadro de dialogo para avisar de que hay cambios que se perderan -->
        <div class="modal-backdrop" v-if="showModal11">
          <div class="modal" role="dialog" id="modal11">
            <form id="formElement-attributesToIndex" style="text-align: justify;">
              <label>
                <h5>If you continue, the changes made will be lost</h5>
              </label>
              <div class="form-group" style="margin-top: 15px">
                <button class="btn" @click="changes = false; showModal11 = false; attributesToIndex()">Continue</button>
                <button class="btn" @click="showModal11 = false">Cancel</button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>

</template>

<script>

import axios from 'axios'

export default {
  name: 'App',

  data () {
    return {
      domains : [],
      policies : [],
      attributes : [],
      resources : [],
      subjects : [],
      actions : [],
      subjectIdTypes : [],
      checksResources : null,
      checksSubjects : null,
      checksActions : null,
      selectedDomain : null,
      selectedPolicy : {
        "PolicyId": {
          "value": ""
        },
        "Rules": []
      },
      selectedRule : {
        "RuleId": {
          "value": ""
        }
      },
      selectedResource : {
        "name": {
          "value": ""
        }
      },
      selectedAction : {
        "name": {
          "value": ""
        }
      },
      selectedSubject : {
        "name": {
          "value": ""
        }
      },
      showIndex : true,
      showPolicies : false,
      showAttributes : false,
      showForms : false,
      showModal1 : false,
      showModal2 : false,
      showModal3 : false,
      showModal4 : false,
      showModal5 : false,
      showModal6 : false,
      showModal7 : false,
      showModal8 : false,
      showModal9 : false,
      showModal10 : false,
      showModal11 : false,
      showModal12 : false,
      showModal13 : false,
      showModal14 : false,
      showModal15 : false,
      domainForm : '',
      attributesForm : '',
      policiesForm : '',
      policyModal : '',
      ruleModal : '',
      resourceModal : '',
      actionModal : '',
      subjectModal : '',
      subjectIDModal : null,
      checkResources : false,
      checkSubjects : false,
      checkActions : false,
      selectedEffect : null,
      changes : false,
      checksChanges : false,
      provisionalPolicy : null,
      provisionalRule : null,
      disableAllChecks : false,
    }
  },

  methods : {
    showForm() {
      this.showForms = !this.showForms
    },

    cleanTextInputs() {
      this.domainForm = ''
      this.attributesForm = ''
      this.policiesForm = ''
      this.policyModal = ''
      this.ruleModal = ''
      this.resourceModal = ''
      this.actionModal = ''
      this.subjectModal = ''
    },

    cleanData() {
      //this.domains = []
      this.policies = []
      this.attributes = []
      this.resources = []
      this.subjects = []
      this.actions = []
      this.subjectIdTypes = []
      this.checksResources = null
      this.checksSubjects = null
      this.checksActions = null
      this.selectedPolicy = {
        "PolicyId": {
          "value": ""
        },
        "Rules": []
      }
      this.selectedRule = {
        "RuleId": {
          "value": ""
        }
      }
      this.selectedResource = {
        "name": {
          "value": ""
        }
      }
      this.selectedAction = {
        "name": {
          "value": ""
        }
      }
      this.selectedSubject = {
        "name": {
          "value": ""
        }
      }
      this.subjectIDModal = null
      this.selectedEffect = null
      this.provisionalPolicy = null
      this.provisionalRule = null
    },
    
    indexToPolicies() {
      if (this.selectedDomain) {
        this.getPolicies(this.selectedDomain)
        this.getAttributes(this.selectedDomain)
        this.showForms = false
        this.showIndex = !this.showIndex
        this.showPolicies = !this.showPolicies
        this.cleanDisableChecks()
      }
    },

    policiesToIndex() {
      if (this.changes || this.checksChanges) {
        this.showModal13 = true;
      } else {
        this.cleanData()
        //this.fetchDomains()
        this.showPolicies = !this.showPolicies
        this.showIndex = !this.showIndex
      }
    },

    policiesExit() {
      this.selectedDomain = null
      this.policiesToIndex()
    },

    indexToAttributes() {
      if (this.selectedDomain) {
        this.getAttributes(this.selectedDomain)
        this.getSubjectIdTypes()
        this.showForms = false
        this.showIndex = !this.showIndex
        this.showAttributes = !this.showAttributes
      }
    },

    attributesToIndex() {
      if (this.changes) {
        this.showModal11 = true;
      } else {
        this.cleanData()
        //this.fetchDomains()
        this.showAttributes = !this.showAttributes
        this.showIndex = !this.showIndex
      }
    },

    attributesExit() {
      this.selectedDomain = null
      this.attributesToIndex()
    },

    async fetchDomains() {
      await axios.get('/pap/api/obtaindomains')
        .then(response => {
          this.domains = response.data.domains
      })
    },

    async addDomain() {

      await axios.patch('/pap/api/saveattributes', '[]', {
        headers : {
          'domain' : this.domainForm,
          'Content-Type': 'application/json'
        }
      })

      await axios.patch('/pap/api/savepolicies', '[]', {
        headers : {
          'domain' : this.domainForm,
          'Content-Type' : 'application/json'
        }
      })

      window.location.reload()
    },

    async getAttributes(selectedDomain) {
      await axios.get('/pap/api/obtainattributes', {
        headers : {
          'domain' : selectedDomain,
        }
      })
      .then(response => {
        this.attributes = response.data
        this.processAttributes()
      })
    },

    processAttributes() {
      var textResources = '{';
      var textSubjects = '{';
      var textActions = '{';
      var aux = '';
      for (var i = 0; i < this.attributes.length; i++) {
        if (this.attributes[i].sortedValue.value == "resource") {
          aux = this.attributes[i].name.value.replace(/"/g, '\\"');
          this.attributes[i].name.value = aux;
          this.resources.push(this.attributes[i]);
          aux = aux.replace(/"/g, '\\"');
          aux = aux.replace(/"/g, '\\"');
          textResources = textResources + '"' + aux + '":false,';
        } else if (this.attributes[i].sortedValue.value == "subject") {
          this.subjects.push(this.attributes[i]);
          textSubjects = textSubjects + '"' + this.attributes[i].name.value + '":false,';
        } else if (this.attributes[i].sortedValue.value == "action") {
          this.actions.push(this.attributes[i]);
          textActions = textActions + '"' + this.attributes[i].name.value + '":false,';
        }
      }
      textResources = textResources.slice(0, textResources.length - 1) + '}';
      textSubjects = textSubjects.slice(0, textSubjects.length - 1) + '}';
      textActions = textActions.slice(0, textActions.length - 1) + '}';
      this.checksResources = JSON.parse(textResources);
      this.checksSubjects = JSON.parse(textSubjects);
      this.checksActions = JSON.parse(textActions);
    },

    selectPolicy(Policy) {
      if (!this.checksChanges) {
        this.selectedPolicy = Policy;
        this.selectedRule = {
          "RuleId": {
            "value": ""
          }
        };
        this.cleanDisableChecks();
      } else {
        this.provisionalPolicy = Policy;
        this.showModal15 = true;
      }
    },
    
    selectRule(Rule) {
      if (!this.checksChanges) {
        this.disableAllChecks = false;
        this.checkResources = false;
        this.checkSubjects = false;
        this.checkActions = false;
        var ruleHasResources = false;
        var ruleHasSubjects = false;
        var ruleHasActions = false;
        var aux = '';

        this.selectedRule = Rule;
        this.selectedEffect = this.selectedRule.Effect.value;

        for (var key1 in this.checksResources) {
          this.checksResources[key1] = false;
        }
        for (var key2 in this.checksSubjects) {
          this.checksSubjects[key2] = false;
        }
        for (var key3 in this.checksActions) {
          this.checksActions[key3] = false;
        }

        for (var key4 in this.selectedRule) {
          if (key4 == "Resources") {
            for (var i = 0; i < this.selectedRule.Resources.length; i++) {
              aux = this.selectedRule.Resources[i].AttributeValue.replace(/"/g, '\\"');
              this.checksResources[aux] = true;
              ruleHasResources = true;
            }
          }
          if (key4 == "Subjects") {
            for (var j = 0; j < this.selectedRule.Subjects.length; j++) {
              this.checksSubjects[this.selectedRule.Subjects[j].AttributeValue] = true;
              ruleHasSubjects = true;
            }
          }
          if (key4 == "Actions") {
            for (var k = 0; k < this.selectedRule.Actions.length; k++) {
              this.checksActions[this.selectedRule.Actions[k].AttributeValue] = true;
              ruleHasActions = true;
            }
          }
        }

        if (!ruleHasResources) {
          this.checkResources = true;
        }
        if (!ruleHasSubjects) {
          this.checkSubjects = true;
        }
        if (!ruleHasActions) {
          this.checkActions = true;
        }
      } else {
        this.provisionalRule = Rule;
        this.showModal14 = true;
      }
    },

    async savePolicies() {
      var newPolicies = this.policies;
      
      if (newPolicies.length > 0 && this.selectedPolicy.PolicyId.value != "" && this.selectedRule.RuleId.value != "") {
        for (var l = 0; l < newPolicies.length; l++) {
          if (newPolicies[l].PolicyId.value == this.selectedPolicy.PolicyId.value) {
            newPolicies.splice(l, 1);
          }
        }
        var newPolicy = this.selectedPolicy;
        for (var m = 0; m < newPolicy.Rules.length; m++) {
          if (newPolicy.Rules[m].RuleId.value == this.selectedRule.RuleId.value) {
            newPolicy.Rules.splice(m, 1);
          }
        }
        var newRule = this.selectedRule;
        newRule.Effect.value = this.selectedEffect;
        delete newRule.Resources;
        delete newRule.Subjects;
        delete newRule.Actions;
        var newResources = [];
        var newSubjects = [];
        var newActions = [];

        for (var key1 in this.checksResources) {
          if (this.checksResources[key1]) {
            for (var i = 0; i < this.resources.length; i++) {
              if (this.resources[i].name.value == key1) {
                var resourceText = '{"AttributeValue": "' + this.resources[i].name.value + '","AttributeId": "' + this.resources[i].xacml_id.value + '"}';
                var resourceJSON = JSON.parse(resourceText);
                newResources.push(resourceJSON);
              }
            }
          }
        }
        for (var key2 in this.checksSubjects) {
          if (this.checksSubjects[key2]) {
            for (var j = 0; j < this.subjects.length; j++) {
              if (this.subjects[j].name.value == key2) {
                var subjectText = '{"AttributeValue": "' + this.subjects[j].name.value + '","AttributeId": "' + this.subjects[j].xacml_id.value + '"}';
                var subjectJSON = JSON.parse(subjectText);
                newSubjects.push(subjectJSON);
              }
            }
          }
        }
        for (var key3 in this.checksActions) {
          if (this.checksActions[key3]) {
            for (var k = 0; k < this.actions.length; k++) {
              if (this.actions[k].name.value == key3) {
                var actionText = '{"AttributeValue": "' + this.actions[k].name.value + '","AttributeId": "' + this.actions[k].xacml_id.value + '"}';
                var actionJSON = JSON.parse(actionText);
                newActions.push(actionJSON);
              }
            }
          }
        }

        //if (newResources.length > 0) {
          newRule.Resources = newResources;
        //}
        //if (newSubjects.length > 0) {
          newRule.Subjects = newSubjects;
        //}
        //if (newActions.length > 0) {
          newRule.Actions = newActions;
        //}
        
        newPolicy.Rules.push(newRule);
        newPolicies.push(newPolicy);
      }

      await axios.patch('/pap/api/savepolicies', newPolicies, {
        headers : {
          'domain' : this.selectedDomain,
          'Content-Type' : 'application/json'
        }
      })
      this.changes = false;
      this.checksChanges = false;
      this.showModal12 = true;
      //window.location.reload();
    },

    async saveAttributes() {
      var newAttributes = [];
      var aux = '';

      for (var i = 0; i < this.resources.length; i++) {
        aux = this.resources[i].name.value.replace(/\\/g, '');
        this.resources[i].name.value = aux;
        newAttributes.push(this.resources[i]);
      }
      for (var j = 0; j < this.actions.length; j++) {
        newAttributes.push(this.actions[j]);
      }
      for (var k = 0; k < this.subjects.length; k++) {
        newAttributes.push(this.subjects[k]);
      }

      await axios.patch('/pap/api/saveattributes', newAttributes, {
        headers : {
          'domain' : this.selectedDomain,
          'Content-Type' : 'application/json'
        }
      })
      this.changes = false;
      this.showModal7 = true;
      for (var l = 0; l < this.resources.length; l++) {
        aux = this.resources[l].name.value.replace(/"/g, '\\"');
        this.resources[l].name.value = aux;
      }
      //window.location.reload();
    },

    addPolicy() {
      var policyText = '{"PolicyId": {"type": "Property", "value": "' + this.policyModal + '"}, "RuleCombiningAlgId": {"type": "Property", "value": "urn:oasis:names:tc:xacml:1.0:rule-combining-algorithm:first-applicable"}, "Rules": []}'
      var policyJSON = JSON.parse(policyText);
      this.policies.push(policyJSON);
      this.changes = true;
      this.showModal1 = false;
    },

    deletePolicy() {
      for (var i = 0; i < this.policies.length; i++) {
        if (this.policies[i].PolicyId.value == this.selectedPolicy.PolicyId.value) {
          this.policies.splice(i, 1);
          this.selectedPolicy = {
            "PolicyId": {
              "value": ""
            },
            "Rules": []
          };
          this.changes = true;
        }
      }
    },

    renamePolicy() {
      this.selectedPolicy.PolicyId.value = this.policyModal;
      this.changes = true;
      this.showModal3 = false;
    },

    addRule() {
      var ruleText = '{"RuleId": {"type": "Property", "value": "' + this.ruleModal + '"}, "Effect": {"type": "Property", "value": "Permit"}}'
      var ruleJSON = JSON.parse(ruleText);
      var tieneRules = false;
      for (var key in this.selectedPolicy) {
        if (key == "Rules") {
          tieneRules = true;
        }
      }
      if (!tieneRules) {
        this.selectedPolicy.Rules = [];
      }
      this.selectedPolicy.Rules.push(ruleJSON);
      this.changes = true;
      this.showModal2 = false;
    },

    deleteRule() {
      for (var i = 0; i < this.selectedPolicy.Rules.length; i++) {
        if (this.selectedPolicy.Rules[i].RuleId.value == this.selectedRule.RuleId.value) {
          this.selectedPolicy.Rules.splice(i, 1);
          this.selectedRule = {
            "RuleId": {
              "value": ""
            }
          };
          this.changes = true;
        }
      }
    },

    addResource() {
      var resourceText = '{"name": {"type": "Property","value": "' + this.resourceModal + '"}, "xacml_id": {"type": "Property","value": "urn:oasis:names:tc:xacml:1.0:resource:resource-id"}, "sortedValue": {"type": "Property","value": "resource"}, "xacml_DataType": {"type": "Property","value": "#string"}}'
      var resourceJSON = JSON.parse(resourceText);
      var aux = '';
      aux = resourceJSON.name.value.replace(/"/g, '\\"');
      resourceJSON.name.value = aux;
      this.resources.push(resourceJSON);
      this.resourceModal = '';
      this.changes = true;
      this.showModal4 = false;
    },

    closeModal4() {
      this.resourceModal = '';
      this.showModal4 = false;
    },

    deleteResource() {
      if (this.selectedResource.name.value != "") {
        for (var i = 0; i < this.resources.length; i++) {
          if (this.resources[i].name.value == this.selectedResource.name.value) {
            this.resources.splice(i, 1);
            this.selectedResource = {
              "name": {
                "value": ""
              }
            };
            this.changes = true;
          }
        }
      } else {
        this.showModal8 = true;
      }
    },

    addAction() {
      var actionText = '{"name": {"type": "Property","value": "' + this.actionModal + '"}, "xacml_id": {"type": "Property","value": "urn:oasis:names:tc:xacml:1.0:action:action-id"}, "sortedValue": {"type": "Property","value": "action"}, "xacml_DataType": {"type": "Property","value": "#string"}}'
      var actionJSON = JSON.parse(actionText);
      this.actions.push(actionJSON);
      this.actionModal = '';
      this.changes = true;
      this.showModal5 = false;
    },

    closeModal5() {
      this.actionModal = '';
      this.showModal5 = false;
    },

    deleteAction() {
      if (this.selectedAction.name.value != "") {
        for (var i = 0; i < this.actions.length; i++) {
          if (this.actions[i].name.value == this.selectedAction.name.value) {
            this.actions.splice(i, 1);
            this.selectedAction = {
              "name": {
                "value": ""
              }
            };
            this.changes = true;
          }
        }
      } else {
        this.showModal9 = true;
      }
    },

    addSubject() {
      var subjectText = '{"name": {"type": "Property","value": "' + this.subjectModal + '"}, "xacml_id": {"type": "Property","value": "' + this.subjectIDModal + '"}, "sortedValue": {"type": "Property","value": "subject"}, "xacml_DataType": {"type": "Property","value": "#string"}}'
      var subjectJSON = JSON.parse(subjectText);
      this.subjects.push(subjectJSON);
      this.subjectModal = '';
      this.subjectIDModal = null;
      this.changes = true;
      this.showModal6 = false;
    },

    closeModal6() {
      this.subjectModal = '';
      this.subjectIDModal = null;
      this.showModal6 = false;
    },

    deleteSubject() {
      if (this.selectedSubject.name.value != "") {
        for (var i = 0; i < this.subjects.length; i++) {
          if (this.subjects[i].name.value == this.selectedSubject.name.value) {
            this.subjects.splice(i, 1);
            this.selectedSubject = {
              "name": {
                "value": ""
              }
            };
            this.changes = true;
          }
        }
      } else {
        this.showModal10 = true;
      }
    },

    checkAllResources() {
      this.checksChanges = true;
      if (!this.checkResources) {
        for (var key in this.checksResources) {
          this.checksResources[key] = false;
        }
      } else {
        for (var key4 in this.selectedRule) {
          if (key4 == "Resources") {
            for (var i = 0; i < this.selectedRule.Resources.length; i++) {
              this.checksResources[this.selectedRule.Resources[i].AttributeValue] = true;
            }
          }
        }
      }
    },

    checkAllSubjects() {
      this.checksChanges = true;
      if (!this.checkSubjects) {
        for (var key in this.checksSubjects) {
          this.checksSubjects[key] = false;
        }
      } else {
        for (var key4 in this.selectedRule) {
          if (key4 == "Subjects") {
            for (var j = 0; j < this.selectedRule.Subjects.length; j++) {
              this.checksSubjects[this.selectedRule.Subjects[j].AttributeValue] = true;
            }
          }
        }
      }
    },

    checkAllActions() {
      this.checksChanges = true;
      if (!this.checkActions) {
        for (var key in this.checksActions) {
          this.checksActions[key] = false;
        }
      } else {
        for (var key4 in this.selectedRule) {
          if (key4 == "Actions") {
            for (var k = 0; k < this.selectedRule.Actions.length; k++) {
              this.checksActions[this.selectedRule.Actions[k].AttributeValue] = true;
            }
          }
        }
      }
    },

    cleanDisableChecks() {
      for (var key1 in this.checksResources) {
        this.checksResources[key1] = false;
      }
      for (var key2 in this.checksSubjects) {
        this.checksSubjects[key2] = false;
      }
      for (var key3 in this.checksActions) {
        this.checksActions[key3] = false;
      }
      this.checkResources = false;
      this.checkSubjects = false;
      this.checkActions = false;
      this.disableAllChecks = true;
    },

    async getPolicies(selectedDomain) {
      await axios.get('/pap/api/obtainpolicies', {
        headers : {
          'domain' : selectedDomain,
        }
      })
      .then(response => {
        this.policies = response.data
      })
    },

    async getSubjectIdTypes() {
      await axios.get('/pap/api/obtainSubjectIdTypes', {
        headers : {}
      })
      .then(response => {
        this.subjectIdTypes = response.data
      })
    },

    formatAttributesData(data) {
      let result = '';
      data.forEach((attribute) => {
        result += `Name: ${attribute.name.value}\n`;
        result += `XACML_ID: ${attribute.xacml_id.value}\n`;
        result += `SortedValue: ${attribute.sortedValue.value}\n`;
        result += `XACML_DataType: ${attribute.xacml_DataType.value}\n\n`;
      });

      return result;
    },

    formatPoliciesData(data) {
      let result = '';
      data.forEach((policy) => {
        result += `PolicyId: ${policy.PolicyId.value}\n`;
        result += `RuleCombiningAlgId: ${policy.RuleCombiningAlgId.value}\n`;

        if (policy.Rules) {
          policy.Rules.forEach((rule) => {
            result += `RuleId: ${rule.RuleId.value}\n`;
            result += `Effect: ${rule.Effect.value}\n`;

            if (rule.Subjects) {
              result += 'Subjects:\n';
              rule.Subjects.forEach((subject) => {
                result += `  ${subject.AttributeValue}\n`;
              });
            }

            if (rule.Resources) {
              result += 'Resources:\n';
              rule.Resources.forEach((resource) => {
                result += `  ${resource.AttributeValue}\n`;
              });
            }

            if (rule.Actions) {
              result += 'Actions:\n';
              rule.Actions.forEach((action) => {
                result += `  ${action.AttributeValue}\n`;
              });
            }
          });
        }

        result += '\n';
      });

      return result;
    },
  },

  async created() {
    await this.fetchDomains()
  }
}

</script>

<style>

.container-fluid {
  text-align: center;
  border: rgb(44, 113, 148) 1.5px solid;
  height: 100%;
  width: 100%;
  font-family: Arial, Helvetica, sans-serif;
}

#mainPAP {
  font-size: small;
  background: rgb(172, 213, 246);
  border: rgb(44, 113, 148) 1.5px solid;
}

#mainPAP-central {
  font-size: xx-large;
  border: rgb(44, 113, 148) 1.5px solid;
  padding: 50px 0px 50px 0px;
}

#MainPAP-menu {
  border: rgb(44, 113, 148) 1.5px solid;
  padding: 0px;
}

#MainPAP-menu-domain, #MainPAP-menu-administration, #MainPAP-menu-configuration {
  border: rgb(44, 113, 148) 1.5px solid;
}

#mainPAP-domain, #mainPAP-administration, #mainPAP-configuration {
  border: rgb(44,113,148) 1.5px solid;
  font-size: small;
  color:rgb(21, 27, 107);
  background: rgb(172, 213, 246);
  font-weight: bold;
  width: 100%;
}

#button1, #button2, #button3 {
  display: inline-block;
  font-size: 14px;
  width: 80%;
  border-radius: 4px;
  padding: auto;
  margin-block-end:20px;
  background-color: rgb(236, 236, 236);
}

#button4 {
  width: 10%;
  height: 50px;
  border-radius: 4px;
  margin-top: 4%;
  margin-bottom:3%;
  font-size:14px;
}

#mainManagePolicies {
  font-size: small;
  background: rgb(172, 213, 246);
  border: rgb(44, 113, 148) 1.5px solid;
}

#mainManagePolicies-central {
  font-size: xx-large;
  border: rgb(44, 113, 148) 1.5px solid;
  padding: 30px 0px 40px 0px;
}

#policiesManagement {
  border: rgb(44, 113, 148) 1px solid;
  background: rgb(172,213,246);
  height:28px;
  padding-top: 3px;
  margin-left: 0px;
  margin-right: 0px;
  font-size: 1.5ch;
  color:rgb(21, 27, 107);
  font-weight: bold;
  width: 100%;
}

.table>tbody>tr.highlight>td {
  background-color: rgb(172 213 246);
}

.modal-backdrop {
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: rgba(0, 0, 0, 0.3);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal {
  background: #FFFFFF;
  box-shadow: 2px 2px 20px 1px;
  overflow-x: auto;
  display: flex;
  flex-direction: column;
  top: 37.5%;
  left: 37.5%;
  width: 25%;
  height: 25%;
  border-top-width: 1px;
  border-bottom-width: 1px;
  border-right-width: 1px;
  border-left-width: 1px;
  padding-top: 16px;
  padding-bottom: 16px;
  padding-right: 16px;
  padding-left: 16px;
}

.highlight {
  background-color: rgb(172 213 246);
}

</style>