<template>
  <v-container fluid>
    <v-dialog
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
    </v-dialog>

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
              Account
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
                  <p>ID: {{ account.id }}</p>
                  <p>Created at: {{ account.createdAt }}</p>
                  <p>Updated at: {{ account.lastUpdateAt }}</p>
                </v-card-text>
              </v-card>
            </v-tab-item>

            <v-tab-item>
              <v-card flat>
                <v-card-text>
                  <v-btn
                    block
                    class="mb-4"
                    :dark="componentDark"
                    :color="btnPrimary"
                    @click="showDeviceDialog"
                  >
                    Add device
                  </v-btn>
                  <v-simple-table>
                    <template v-slot:default>
                      <thead>
                      <tr>
                        <th class="text-left">Name</th>
                        <th class="text-left">Token</th>
                      </tr>
                      </thead>
                      <tbody>
                      <tr v-for="device in devices" :key="device.id">
                        <td>{{ device.name }}</td>
                        <td>{{ device.token }}</td>
                      </tr>
                      </tbody>
                    </template>
                  </v-simple-table>
                </v-card-text>
              </v-card>
            </v-tab-item>

            <v-tab-item>
              <v-card flat>
                <v-card-text>
                  <v-btn
                    block
                    class="mb-4"
                    :dark="componentDark"
                    :color="btnPrimary"
                    @click="showProjectDialog"
                  >
                    Add project
                  </v-btn>
                  <v-simple-table>
                    <template v-slot:default>
                      <thead>
                      <tr>
                        <th class="text-left">Name</th>
                        <th class="text-left">Created by</th>
                      </tr>
                      </thead>
                      <tbody>
                      <tr v-for="project in projects" :key="project.id">
                        <td>{{ project.title }}</td>
                        <td>{{ project.user }}</td>
                      </tr>
                      </tbody>
                    </template>
                  </v-simple-table>
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

export default {
  name: 'settings',
  data () {
    return {
      addDevice: false,
      newDevice: '',
      addProject: false,
      newProject: '',
      account: {
        id: '2',
        createdAt: '15:30 04/04/2020',
        lastUpdateAt: '15:31 04/04/2020',
      },
      devices: [{
        id: '1',
        name: 'Galaxy',
        token: 'long-string-consisting-of-numbers-and-letters',
      },
        {
          id: '2',
          name: 'NotGalaxy',
          token: 'long-string-consisting-of-numbers-and-letters-v2',
        }],
      projects: [{
        id: '5',
        title: 'New Project',
        user: 'kitty',
      },
        {
          id: '7',
          title: 'New Project 2',
          user: 'notkitty',
        }],
    }
  },
  mounted () {
  },
  methods: {
    showDeviceDialog () {
      this.addDevice = true
    },
    showProjectDialog () {
      this.addProject = true
    },
    saveDevice () {
      if (this.newDevice !== '') {
        this.devices.push({'name': this.newDevice})
        this.newDevice = ''
      }
      this.addDevice = false
    },
    saveProject () {
      if (this.newProject !== '') {
        this.projects.push({'title': this.newProject})
        this.newProject = ''
      }
      this.addProject = false
    },
    authConfig () {
      return {headers: {Authorization: 'Bearer ' + localStorage.getItem(this.lsToken)}}
    },

  }
}
</script>
