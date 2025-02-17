<template>
  <div class="container mt-5">
    <h1 class="mb-4 text-background">XACML-PAP-FrontEnd-DLT-Multidomain</h1>
    <div class="row mb-3">
      <div class="col-md-6">
        <h2 class="text-background">Select a domain:</h2>
        <div class="d-flex align items-center">
          <select v-model="selectedDomain" class="form-select">
            <option v-for="domain in domains" :key="domain" :value="domain">{{ domain }}</option>
          </select>
          <ComponentButton :text="'Confirm'" @btn-click=" () => {sendDomain(); if (selectedDomain) showInfo=true; }" class="ms-3"/>
          <ComponentButton :text="'Add'" @btn-click="showForms = !showForms" class="btn btn-danger ms-3"/>
        </div>
      </div>
    </div>

    <div class="row">
      <div class="col-md-6">
        <h2 v-if="showInfo" class="text-background"> XACML Policies </h2>
        <textarea 
          v-if="showInfo" 
          v-model="policies" 
          rows="10" 
          class="form-control mb-3" 
          style="resize: none;"
          readonly>
        </textarea>
      </div>
      <div class="col-md-6">
        <h2 v-if="showInfo" class="text-background"> XACML Attributes </h2>
        <textarea 
          v-if="showInfo" 
          v-model="attributes" 
          rows="10" 
          class="form-control mb-3" 
          style="resize: none;"
          readonly>
        </textarea>
      </div>
      <div class="row my-3">
        <div class="d-grid gap-1" v-if="showInfo">
          <ComponentButton :text="'Hide'" @btn-click="() => { showInfo = !showInfo }" class="btn btn-primary ms-3"/> 
        </div>
      </div>
    </div>
    
    <div v-if="showForms">
      <div class="input-group mb-3">
        <input type="text" class="form-control" v-model="domainForm" placeholder="Domain" aria-label="Domain">
      </div>

      <div class="row">
        <div class="col-md-6">
          <div class="input-group">
            <span class="input-group-text"><b>Attributes</b></span>
            <textarea class="form-control" v-model="attributesForm" aria-label="Attributes" rows="6"></textarea>
          </div>
        </div>

        <div class="col-md-6">
          <div class="input-group">
            <span class="input-group-text"><b>Policies</b></span>
            <textarea class="form-control" v-model="policiesForm" aria-label="Policies" rows="6"></textarea>
          </div>
        </div>
      </div>
      <div class="row my-3">
        <div class="d-grid gap-1">
          <ComponentButton :text="'Submit'" @btn-click="() => { addDomain(); showForms = !showForms; }" class="btn btn-danger ms-3"/> 
        </div>
      </div>
    </div>

  </div>
</template>

<script>
import ComponentButton from './components/ComponentButton.vue'
import axios from 'axios'

export default {
  name: 'App',
  components: {
    ComponentButton,
  },

  data () {
    return {
      domains : [],
      attributes : '',
      policies : '',
      selectedDomain : null,
      showForms : false,
      showInfo : false,
      domainForm : '',
      attributesForm : '',
      policiesForm : '',
    }
  },

  methods : {
    async fetchDomains() {
      await axios.get('/api/obtaindomains')
        .then(response => {
          this.domains = response.data.domains
      })
    },

    sendDomain() {
      this.getAttributes(this.selectedDomain)
      this.getPolicies(this.selectedDomain)
    },

    async addDomain() {

      await axios.patch('/api/saveattributes', this.attributesForm , {
        headers : {
          'domain' : this.domainForm,
          'Content-Type': 'application/json'
        }
      })

      await axios.patch('/api/savepolicies', this.policiesForm , {
        headers : {
          'domain' : this.domainForm,
          'Content-Type' : 'application/json'
        }
      })

      window.location.reload()
    },

    async getAttributes(selectedDomain) {

      await axios.get('/api/obtainattributes', {
        headers : {
          'domain' : selectedDomain,
        }
      })
      .then(response => {
        this.attributes = this.formatAttributesData(response.data)
      })
    },

    async getPolicies(selectedDomain) {
      await axios.get('/api/obtainpolicies', {
        headers : {
          'domain' : selectedDomain,
        }
      })
      .then(response => {
        this.policies = this.formatPoliciesData(response.data)
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

h1 {

  font-size: 2rem; 
  padding: 10px; 
}

.text-background {
  background-color: #ffffffb0; 
  display: inline; 
}

.form-select {
  width: 100%; 
}

textarea {
  font-family: monospace;
  resize: none;
}

#app {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-height: 100vh;
  justify-content: center;
}

.component-button {
  background-color: #007bff; 
  color: #fff; 
  border: none;
  padding: 10px 20px;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.component-button:hover {
  background-color: #0056b3; 
}

.component-button:active {
  background-color: #004099; 
}

.component-button:focus {
  outline: none; 
}

body {
  background-image: url('./assets/animacion-inicio.gif');
  background-attachment: fixed;
  background-position: center center;
  background-size: cover;
}

</style>