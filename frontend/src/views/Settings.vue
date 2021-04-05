<template>
  <v-container fluid>
    <v-dialog
      v-model="showNewAccountDialog"
      width="600"
      :fullscreen="$vuetify.breakpoint.xsOnly"
    >
      <v-card>
        <v-card-title>
          Create new account
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  required
                  clearable
                  counter
                  maxlength="24"
                  label="Username"
                  prepend-icon="mdi-account"
                  type="text"
                  v-model="username"
                  :color="inputPrimary"
                >
                </v-text-field>
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  required
                  clearable
                  counter
                  maxlength="32"
                  label="Password"
                  prepend-icon="mdi-lock"
                  type="password"
                  v-model="password"
                  :color="inputPrimary"
                >
                </v-text-field>
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  required
                  clearable
                  label="Email"
                  name="email"
                  prepend-icon="mdi-email"
                  v-model="email"
                  :color="inputPrimary"
                >
                </v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-divider></v-divider>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            :color="btnSecondary"
            text
            @click="showNewAccountDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn
            :dark="componentDark"
            :color="btnPrimary"
            @click.native="addAccount"
          >
            Add
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <!-- <v-dialog
      v-model="addDevice"
      width="500"
      :fullscreen="$vuetify.breakpoint.xsOnly"
    >
      <v-card>
        <v-card-title>
          Add new device
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12" sm="6" md="4">
                <v-text-field
                  label="Device name"
                  required
                  v-model="newDevice"
                  :color="inputPrimary"
                >
                </v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-divider></v-divider>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            :color="btnSecondary"
            text
            @click="addDevice = false"
          >
            Cancel
          </v-btn>
          <v-btn
            :dark="componentDark"
            :color="btnPrimary"
            @click="saveDevice"
          >
            Add
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog
      v-model="addProject"
      width="500"
      :fullscreen="$vuetify.breakpoint.xsOnly"
    >
      <v-card>
        <v-card-title>
          Begin new project
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12" sm="6" md="4">
                <v-text-field
                  label="Project name"
                  required
                  v-model="newProject"
                  :color="inputPrimary"
                >
                </v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-divider></v-divider>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            :color="btnSecondary"
            text
            @click="addProject = false"
          >
            Cancel
          </v-btn>
          <v-btn
            :dark="componentDark"
            :color="btnPrimary"
            @click="saveProject"
          >
            Save
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog> -->

    <v-row>
      <v-col>
        <v-card
          max-width="800"
          class="mx-auto"
        >
          <v-toolbar
            :color="toolbarPrimary"
            :dark="componentDark"
          >
            <v-toolbar-title>Settings</v-toolbar-title>
          </v-toolbar>

          <v-tabs
            vertical
            :color="tabPrimary"
          >
            <v-tab>
              <v-icon left>
                mdi-account
              </v-icon>
              Accounts
            </v-tab>
            <v-tab>
              <v-icon left>
                mdi-monitor-cellphone
              </v-icon>
              Devices
            </v-tab>
            <v-tab>
              <v-icon left>
                mdi-folder-open
              </v-icon>
              Projects
            </v-tab>

            <v-tab-item>
              <v-card flat>
                <v-card-text>
                  <v-btn
                    block
                    class="mb-4"
                    :dark="componentDark"
                    :color="btnPrimary"
                    @click="showNewAccountDialog = true"
                  >
                    Create new account
                  </v-btn>
                  <v-data-table
                    :headers="accountHeaders"
                    :items="accounts"
                    :loading="tableLoading"
                  >
                    <template v-slot:progress>
                      <v-progress-linear
                        indeterminate
                        :height="2"
                        :color="progressPrimary"
                      ></v-progress-linear>
                    </template>
                  </v-data-table>
                </v-card-text>
              </v-card>
            </v-tab-item>

            <v-tab-item>
              <v-card flat>
                <v-card-text>
                  <!-- if disabled then element cannot be dark -->
                  <!-- :dark="componentDark" -->
                  <!-- @click="showDeviceDialog" -->
                  <v-btn
                    block
                    class="mb-4"
                    disabled
                    :color="btnPrimary"
                  >
                    Add device
                  </v-btn>
                  <v-data-table
                    :headers="deviceHeaders"
                    :items="devices"
                    :loading="tableLoading"
                  >
                    <template v-slot:progress>
                      <v-progress-linear
                        indeterminate
                        :height="2"
                        :color="progressPrimary"
                      ></v-progress-linear>
                    </template>
                  </v-data-table>
                </v-card-text>
              </v-card>
            </v-tab-item>

            <v-tab-item>
              <v-card flat>
                <v-card-text>
                  <!-- if disabled then element cannot be dark -->
                  <!-- :dark="componentDark" -->
                  <!-- @click="showProjectDialog" -->
                  <v-btn
                    block
                    class="mb-4"
                    disabled
                    :color="btnPrimary"
                  >
                    Add project
                  </v-btn>
                  <v-data-table
                    :headers="projectHeaders"
                    :items="projects"
                    :loading="tableLoading"
                  >
                    <template v-slot:progress>
                      <v-progress-linear
                        indeterminate
                        :height="2"
                        :color="progressPrimary"
                      ></v-progress-linear>
                    </template>
                  </v-data-table>
                </v-card-text>
              </v-card>
            </v-tab-item>

          </v-tabs>
        </v-card>
      </v-col>
    </v-row>

  </v-container>
</template>

<script>

import axios from 'axios'

export default {
  name: 'settings',
  data () {
    return {
      showNewAccountDialog: false,
      username: '',
      password: '',
      email: '',
      // addDevice: false,
      // newDevice: '',
      // addProject: false,
      // newProject: '',
      tableLoading: true,
      accounts: [],
      accountHeaders: [
        {
          text: 'ID',
          align: 'left',
          sortable: true,
          value: 'id',
        },
        {text: 'Username', value: 'username'},
        // {text: 'Created at', value: 'createdAt'},
        // {text: 'Updated at', value: 'updatedAt'},
      ],
      devices: [],
      deviceHeaders: [
        {
          text: 'ID',
          align: 'left',
          sortable: true,
          value: 'id',
        },
        {text: 'Name', value: 'name'},
        {text: 'Username', value: 'username'},
        // FIXME device token is not available by API, but is present in DB
        // {text: 'Token', value: 'token'},
      ],
      projects: [{}],
      projectHeaders: [
        {
          text: 'ID',
          align: 'left',
          sortable: true,
          value: 'id',
        },
        {text: 'Project name', value: 'name'},
      ],
    }
  },
  mounted () {
    this.getAccounts()
    this.getDevices()
    this.getProjects()
  },
  methods: {
    // showDeviceDialog () {
    //   this.addDevice = true
    // },
    // showProjectDialog () {
    //   this.addProject = true
    // },
    // saveDevice () {
    //   if (this.newDevice !== '') {
    //     this.devices.push({'name': this.newDevice})
    //     this.newDevice = ''
    //   }
    //   this.addDevice = false
    // },
    // saveProject () {
    //   if (this.newProject !== '') {
    //     this.projects.push({'title': this.newProject})
    //     this.newProject = ''
    //   }
    //   this.addProject = false
    // },
    authConfig () {
      return {headers: {Authorization: 'Bearer ' + localStorage.getItem(this.lsToken)}}
    },
    // get data
    getAccounts () {
      this.tableLoading = true
      axios.get(this.apiUrl + '/api/v1/user', this.authConfig())
        .then((res) => {
          this.tableLoading = false
          this.accounts = res.data
        })
        .catch((err) => {
          if (err.response.status === 401) {
            this.$root.$emit('sessionExpired')
          } else {
            console.log('something wrong')
          }
        })
    },
    getDevices () {
      this.tableLoading = true
      axios.get(this.apiUrl + '/api/v1/device', this.authConfig())
        .then((res) => {
          this.tableLoading = false
          this.devices = res.data
        })
        .catch((err) => {
          if (err.response.status === 401) {
            this.$root.$emit('sessionExpired')
          } else {
            console.log('something wrong')
          }
        })
    },
    getProjects () {
      this.tableLoading = true
      axios.get(this.apiUrl + '/api/v1/project', this.authConfig())
        .then((res) => {
          this.tableLoading = false
          let projectList = []
          res.data.forEach((project) => {
            if (project.id > 0) {
              projectList.push(project)
            }
          })
          this.projects = projectList
        })
        .catch((err) => {
          if (err.response.status === 401) {
            this.$root.$emit('sessionExpired')
          } else {
            console.log('something wrong')
          }
        })
    },
    // post data
    addAccount () {
      axios.post(this.apiUrl + '/api/v1/user', {
        'username': this.username,
        'password': this.password,
        'email': this.email
      }, this.authConfig())
        .then((res) => {
          //clear form
          this.username = ''
          this.password = ''
          this.email = ''
          //close dialog
          this.showNewAccountDialog = false
          //get account list
          this.getAccounts()
        })
        .catch((err) => {
          if (err.response.status === 401) {
            this.$root.$emit('sessionExpired')
          } else {
            console.log('something wrong')
          }
        })
    },
  }
}
</script>
