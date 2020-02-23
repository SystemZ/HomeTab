<template>
    <v-container fluid>
        <v-dialog
                v-model="dialog"
                max-width="700"
                :fullscreen="$vuetify.breakpoint.xsOnly"
        >
            <v-card>
                <v-card-title class="headline">{{taskTitleInDialog}}</v-card-title>

                <v-card-text>
                    <v-container>
                        <v-row>
                            <v-col cols="12">
                                <v-text-field
                                        label="Task title"
                                        v-model="taskTitleInDialog"
                                >
                                </v-text-field>
                            </v-col>
                            <v-col cols="12">
                                <v-text-field
                                        label="Additional info"
                                        v-model="taskInfoInDialog"
                                ></v-text-field>
                            </v-col>
                            <v-col cols="12" sm="6">
                                <v-select
                                        :items="['S', 'P', 'Others']"
                                        label="Assigned to"
                                ></v-select>
                            </v-col>
                            <v-col cols="12" sm="6">
                                <v-autocomplete
                                        :items="['Shopping', 'Reading', 'Writing', 'Coding', 'Cleaning']"
                                        label="Tags"
                                        multiple
                                ></v-autocomplete>
                            </v-col>
                        </v-row>
                    </v-container>
                </v-card-text>

                <v-card-actions>
                    <v-spacer></v-spacer>

                    <v-btn
                            color="red darken-1"
                            text
                            @click="dialog = false"
                    >
                        Cancel
                    </v-btn>

                    <v-btn
                            color="green darken-1"
                            dark
                            @click="saveTask"
                    >
                        Submit
                    </v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
        <v-row>
            <v-col md="3"></v-col>
            <v-col cols="12" xs="12" md="6">
                <h1 class="headline mt-5 mb-5">New task</h1>
                <v-text-field
                        placeholder="Buy a yacht"
                        solo
                        clearable
                        v-model="newTaskTitle"
                        @keydown.enter.native="addTask"
                >
                </v-text-field>
                <v-btn
                        block
                        color="success"
                        @click.native="addTask"
                >
                    Add
                </v-btn>
            </v-col>
        </v-row>
        <v-row>
            <v-col md="3"></v-col>
            <v-col cols="12" xs="12" md="6">
                <v-card
                        class="mx-auto"
                >
                    <v-toolbar
                            color="green"
                            dark
                    >
                        <v-select
                                class="pt-8"
                                :items="projectList"
                                v-model="projectIdSelected"
                                label="Project"
                                item-text="name"
                                item-value="id"
                                outlined
                                @change="getTasks(projectIdSelected)"
                        ></v-select>

                        <v-spacer></v-spacer>
                        <v-btn icon>
                            <v-icon>mdi-check-bold</v-icon>
                        </v-btn>
                        <v-btn icon>
                            <v-icon>mdi-alarm-snooze</v-icon>
                        </v-btn>
                        <v-btn icon @click="deleteTasks">
                            <v-icon>mdi-delete</v-icon>
                        </v-btn>
                    </v-toolbar>
                    <v-progress-linear v-if="tasksLoading || projectLoading" indeterminate/>
                    <v-list
                            subheader
                    >
                        <v-subheader>Tasks</v-subheader>
                        <template v-for="(task, i) in tasks">
                            <v-list-item
                                    :key="`${i}-${task.title}`"
                            >
                                <v-list-item-action>
                                    <v-checkbox
                                            v-model="task.selected"
                                            :color="task.selected && 'green darken-2' || 'grey'"
                                    >
                                        <template v-slot:label>
                                            <div
                                                    :class="task.selected && 'grey--text' || 'black--text'"
                                                    class="ml-4"
                                                    v-text="task.title"
                                            >
                                            </div>
                                        </template>
                                    </v-checkbox>
                                </v-list-item-action>
                                <v-spacer></v-spacer>
                                <v-icon @click.stop="showDialog(task)">mdi-information-outline</v-icon>
                            </v-list-item>
                        </template>
                    </v-list>
                </v-card>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
    import axios from "axios";

    export default {
        name: 'tasks',
        data() {
            return {
                projectLoading: true,
                projectList: [{}],
                projectIdSelected: 0,
                newTaskTitle: '',
                tasks: [
                    {
                        id: 0,
                        title: '',
                        selected: false,
                        info: '',
                    }
                ],
                dialog: false,
                taskTitleInDialog: '',
                taskInfoInDialog: '',
                taskIdInDialog: 0,
                tasksLoading: true,
            }
        },
        mounted() {
            this.getProjects()
        },
        methods: {
            addTask() {
                if (this.newTaskTitle !== '') {
                    this.tasks.push({
                        'id': new Date().getUTCMilliseconds(),
                        'title': this.newTaskTitle,
                        'selected': false,
                    })
                    this.newTaskTitle = ''
                }
            },
            deleteTasks() {
                this.tasks = this.tasks.filter(function (task) {
                    return !task.selected
                })
            },
            showDialog(task) {
                this.taskIdInDialog = task.id
                this.taskTitleInDialog = task.title
                this.taskInfoInDialog = task.info
                this.tasks.selected = false
                this.dialog = true
            },
            saveTask() {
                let i
                for (i = 0; i < this.tasks.length; i++) {
                    if (this.taskIdInDialog === this.tasks[i].id) {
                        this.tasks[i].title = this.taskTitleInDialog
                        this.tasks[i].info = this.taskInfoInDialog
                    }
                }
                this.dialog = false
            },
            authConfig() {
                return {headers: {Authorization: 'Bearer ' + localStorage.getItem(this.lsToken)}}
            },
            getProjects() {
                this.projectLoading = true
                axios.get(this.apiUrl + '/api/v1/project', this.authConfig())
                    .then((res) => {
                        this.projectList = res.data
                        this.getTasks(res.data[0].id)
                        this.projectIdSelected = res.data[0].id
                        this.projectLoading = false
                    })
                    .catch(function (err) {
                        if (err.response.status === 401) {
                            console.log('logged out')
                            vm.$root.$emit('sessionExpired')
                        } else if (err.response.status === 400) {
                            console.log('empty result / wrong request')
                        } else {
                            console.log('something wrong')
                        }
                    })
            },
            getTasks(projectId) {
                this.tasksLoading = true
                let url = this.apiUrl + '/api/v1/project/' + projectId + '/task'
                axios.get(url, this.authConfig())
                    .then((res) => {
                        this.tasks = res.data
                        this.tasksLoading = false
                    })
                    .catch(function (err) {
                        if (err.response.status === 401) {
                            console.log('logged out')
                            vm.$root.$emit('sessionExpired')
                        } else if (err.response.status === 400) {
                            console.log('empty result / wrong request')
                        } else {
                            console.log('something wrong')
                        }
                    })
            }
        },
    }
</script>
