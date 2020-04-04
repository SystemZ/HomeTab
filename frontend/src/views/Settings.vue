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
                            :dark="btnDark"
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
                            :dark="btnDark"
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
                            color="#5C9DA0"
                            dark
                    >
                        <v-toolbar-title>Account details</v-toolbar-title>
                    </v-toolbar>

                    <v-list dense>
                        <v-list-item>
                            <v-list-item-content>ID:</v-list-item-content>
                            <v-list-item-content class="align-end">{{account.id}}</v-list-item-content>
                        </v-list-item>
                        <v-list-item>
                            <v-list-item-content>Created at:</v-list-item-content>
                            <v-list-item-content class="align-end">{{account.createdAt}}</v-list-item-content>
                        </v-list-item>
                        <v-list-item>
                            <v-list-item-content>Updated at:</v-list-item-content>
                            <v-list-item-content class="align-end">{{account.lastUpdateAt}}</v-list-item-content>
                        </v-list-item>
                    </v-list>
                </v-card>
            </v-col>
        </v-row>
        <v-row>
            <v-col>
                <v-card
                        max-width="800"
                        class="mx-auto"
                >
                    <v-toolbar
                            color="#5C9DA0"
                            dark
                    >
                        <v-toolbar-title>Devices details</v-toolbar-title>
                        <v-spacer></v-spacer>
                        <v-btn icon @click="showDeviceDialog">
                            <v-icon>mdi-plus</v-icon>
                        </v-btn>
                    </v-toolbar>

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
                </v-card>
            </v-col>
        </v-row>
        <v-row>
            <v-col>
                <v-card
                        max-width="800"
                        class="mx-auto"
                >
                    <v-toolbar
                            color="#5C9DA0"
                            dark
                    >
                        <v-toolbar-title>All projects</v-toolbar-title>
                        <v-spacer></v-spacer>
                        <v-btn icon @click="showProjectDialog">
                            <v-icon>mdi-plus</v-icon>
                        </v-btn>
                    </v-toolbar>

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
                </v-card>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
    import axios from "axios";

    export default {
        name: 'settings',
        data() {
            return {
                addDevice: false,
                newDevice: "",
                addProject: false,
                newProject: "",
                account: {
                    id: "2",
                    createdAt: "15:30 04/04/2020",
                    lastUpdateAt: "15:31 04/04/2020",
                },
                devices: [{
                    id: "1",
                    name: "Galaxy",
                    token: "long-string-consisting-of-numbers-and-letters",
                },
                    {
                        id: "2",
                        name: "NotGalaxy",
                        token: "long-string-consisting-of-numbers-and-letters-v2",
                    }],
                projects: [{
                    id: "5",
                    title: "New Project",
                    user: "kitty",
                },
                    {
                        id: "7",
                        title: "New Project 2",
                        user: "notkitty",
                    }],
            }
        },
        mounted() {
        },
        methods: {
            showDeviceDialog() {
                this.addDevice = true
            },
            showProjectDialog() {
                this.addProject = true
            },
            saveDevice() {
                if (this.newDevice !== "") {
                    this.devices.push({"name": this.newDevice});
                    this.newDevice = ""
                }
                this.addDevice = false
            },
            saveProject() {
                if (this.newProject !== "") {
                    this.projects.push({"title": this.newProject});
                    this.newProject = ""
                }
                this.addProject = false
            },
            authConfig() {
                return {headers: {Authorization: "Bearer " + localStorage.getItem(this.lsToken)}}
            },

        }
    }
</script>
