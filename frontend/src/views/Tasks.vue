<template>
    <v-container fluid>
        <v-dialog v-model="deleteTaskDialog" max-width="700" :fullscreen="$vuetify.breakpoint.xsOnly">
            <v-card>
                <v-card-title>
                    Delete task
                </v-card-title>
                <v-card-subtitle>
                    Are you sure that you don't need these tasks?
                </v-card-subtitle>
                <v-card-text>
                </v-card-text>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn
                            color="green darken-1"
                            text
                            @click="deleteTaskDialog = false"
                    >
                        Cancel
                    </v-btn>
                    <v-btn
                            color="red darken-1"
                            dark
                            :disabled="tasksDeleting"
                            @click="deleteTasks"
                    >
                        <v-icon>mdi-delete</v-icon>
                        Delete
                    </v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
        <v-dialog v-model="snoozeTaskDialog" max-width="700" :fullscreen="$vuetify.breakpoint.xsOnly">
            <v-card>
                <v-card-title>
                    Snooze task
                </v-card-title>
                <v-card-subtitle>
                    Select date after I should nag you again
                </v-card-subtitle>
                <v-card-text>
                    <v-row>
                        <v-col>
                            <v-date-picker
                                    v-model="taskSnoozeDateInDialog"
                                    class="mt-4"
                                    :min="taskSnoozeDateInDialogMin"
                            ></v-date-picker>
                        </v-col>
                        <v-col>
                            <v-time-picker v-model="taskSnoozeTimeInDialog" format="24hr"></v-time-picker>
                        </v-col>
                    </v-row>
                </v-card-text>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn
                            color="red darken-1"
                            text
                            @click="snoozeTaskDialog = false"
                    >
                        Cancel
                    </v-btn>
                    <v-btn
                            color="green darken-1"
                            dark
                            @click="snoozeTasks"
                    >
                        Snooze
                    </v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
        <v-dialog v-model="editTaskDialog" max-width="700" :fullscreen="$vuetify.breakpoint.xsOnly">
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
                        </v-row>
                        <v-row>
                            <v-col cols="12" sm="6">
                                <v-autocomplete
                                        :items="repeatUnits"
                                        v-model="taskRepeatUnitInDialog"
                                        label="Repeat unit"
                                ></v-autocomplete>
                            </v-col>
                            <v-col v-if="taskRepeatUnitInDialog !== ''" cols="12" sm="6">
                                <v-text-field
                                        v-model="taskRepeatEveryInDialog"
                                        label="Repeat every..."
                                ></v-text-field>
                            </v-col>
                        </v-row>
                        <v-row>
                            <v-col cols="12" sm="6">
                                <v-select
                                        v-model="taskAssignedInDialog"
                                        item-text="username"
                                        item-value="id"
                                        :items="userList"
                                        label="Assigned to"
                                ></v-select>
                            </v-col>
                            <!--
                            <v-col cols="12">
                                <v-text-field
                                        label="Additional info"
                                        v-model="taskInfoInDialog"
                                ></v-text-field>
                            </v-col>
                            <v-col cols="12" sm="6">
                                <v-autocomplete
                                        :items="['Shopping', 'Reading', 'Writing', 'Coding', 'Cleaning']"
                                        label="Tags"
                                        multiple
                                ></v-autocomplete>
                            </v-col>
                            -->
                        </v-row>
                    </v-container>
                </v-card-text>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn
                            color="red darken-1"
                            text
                            @click="editTaskDialog = false"
                    >
                        Cancel
                    </v-btn>
                    <v-btn
                            color="green darken-1"
                            dark
                            :disabled="taskSaving"
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
                <v-row>
                    <v-col cols="12" md="6" xs="12">
                        <v-btn
                                block
                                color="success"
                                @click.native="addTask"
                        >
                            <v-icon>mdi-plus</v-icon>
                            Add
                        </v-btn>
                    </v-col>
                    <v-col cols="12" md="6" xs="12">
                        <v-btn
                                block
                                color="primary"
                                @click="refreshTasks"
                        >
                            <v-icon>mdi-refresh</v-icon>
                            Refresh
                        </v-btn>
                    </v-col>
                </v-row>
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
                        <v-btn icon @click="setAsDoneTasks" :disabled="tasksDoneInProgress">
                            <v-icon>mdi-check-bold</v-icon>
                        </v-btn>
                        <v-btn icon @click="showSnoozeDialog">
                            <v-icon>mdi-alarm-snooze</v-icon>
                        </v-btn>
                        <v-btn :disabled="tasksDeleting" icon @click="confirmDeleteTasks">
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
                                <div class="mr-2">
                                    {{task.username}}
                                </div>
                                <v-icon @click.stop="showDialog(task)">mdi-pencil</v-icon>
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
                editTaskDialog: false,
                snoozeTaskDialog: false,
                deleteTaskDialog: false,
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
                taskTitleInDialog: '',
                taskInfoInDialog: '',
                taskIdInDialog: 0,
                taskAssignedInDialog: 0,
                taskSnoozeDateInDialog: '',
                taskSnoozeDateInDialogMin: '',
                taskSnoozeTimeInDialog: '',
                taskRepeatUnitInDialog: '',
                taskRepeatEveryInDialog: 0,
                tasksLoading: true,
                tasksDeleting: false,
                taskSaving: false,
                tasksDoneInProgress: false,
                userList: [],
                repeatUnits: [
                    {'text': 'No repeat', 'value': ''},
                    {'text': 'Minute', 'value': 'i'},
                    {'text': 'Hour', 'value': 'h'},
                    {'text': 'Day', 'value': 'd'},
                    {'text': 'Week', 'value': 'w'},
                    {'text': 'Month', 'value': 'm'},
                    {'text': 'Year', 'value': 'y'}
                ]
            }
        },
        mounted() {
            this.getUsers()
        },
        methods: {
            addTask() {
                if (this.newTaskTitle.length < 1) {
                    // TODO show snackbar
                    return
                }
                this.tasksLoading = true
                let url = this.apiUrl + '/api/v1/project/' + this.projectIdSelected + '/task'
                let data = {"title": this.newTaskTitle}
                axios.post(url, data, this.authConfig())
                    .then((res) => {
                        this.getTasks(this.projectIdSelected)
                    })
                    .catch(((err) => {
                        if (err.response.status === 401) {
                            console.log('logged out')
                            this.$root.$emit('sessionExpired')
                        } else if (err.response.status === 400) {
                            console.log('empty result / wrong request')
                        } else {
                            console.log('something wrong')
                        }
                    }))
                this.newTaskTitle = ''
            },
            showDialog(task) {
                this.taskIdInDialog = task.id
                this.taskTitleInDialog = task.title
                this.taskInfoInDialog = task.info
                this.taskAssignedInDialog = task.assignedTo
                this.taskRepeatUnitInDialog = task.repeatUnit
                this.taskRepeatEveryInDialog = task.repeatEvery
                this.tasks.selected = false
                this.editTaskDialog = true
            },
            showSnoozeDialog(task) {
                let now = new Date()
                let year = now.getFullYear()
                let month = now.getMonth()
                month++
                if (month < 10) {
                    month = "0" + month
                }
                let day = now.getDate()
                let today = year + "-" + month + "-" + day

                this.taskSnoozeDateInDialogMin = today
                this.taskSnoozeDateInDialog = today
                this.snoozeTaskDialog = true
            },
            confirmDeleteTasks() {
                this.deleteTaskDialog = true
            },
            deleteTasks() {
                this.tasksDeleting = true
                // get selected tasks
                let tasksForDelete = []
                this.tasks.forEach((task) => {
                    if (task.selected) {
                        tasksForDelete.push({"id": task.id, "delete": true})
                    }
                })
                // send task IDs to server
                let url = this.apiUrl + '/api/v1/project/' + this.projectIdSelected + '/task'
                axios.put(url, tasksForDelete, this.authConfig())
                    .then((res) => {
                        this.tasksDeleting = false
                        this.deleteTaskDialog = false
                        this.getTasks(this.projectIdSelected)
                    })
                    .catch(((err) => {
                        if (err.response.status === 401) {
                            console.log('logged out')
                            this.$root.$emit('sessionExpired')
                        } else if (err.response.status === 400) {
                            console.log('empty result / wrong request')
                        } else {
                            console.log('something wrong')
                        }
                    }))
                // hide dialog
                this.snoozeTaskDialog = false
                // refresh task list
                this.getTasks(this.projectIdSelected)
                // TODO add snackbar with info for user
            },
            setAsDoneTasks() {
                this.tasksDoneInProgress = true
                // get selected tasks
                let tasksForDone = []
                this.tasks.forEach((task) => {
                    if (task.selected) {
                        tasksForDone.push({"id": task.id, "done": true})
                    }
                })
                // send task IDs to server
                let url = this.apiUrl + '/api/v1/project/' + this.projectIdSelected + '/task'
                axios.put(url, tasksForDone, this.authConfig())
                    .then((res) => {
                        this.tasksDoneInProgress = false
                        // refresh task list
                        this.getTasks(this.projectIdSelected)
                        // TODO add snackbar with info for user
                    })
                    .catch(((err) => {
                        if (err.response.status === 401) {
                            console.log('logged out')
                            this.$root.$emit('sessionExpired')
                        } else if (err.response.status === 400) {
                            console.log('empty result / wrong request')
                        } else {
                            console.log('something wrong')
                        }
                    }))
            },
            snoozeTasks() {
                // get selected tasks
                let tasksForSnooze = []
                this.tasks.forEach((task) => {
                    if (task.selected) {
                        // 2020-02-23T17:47:36Z
                        let offset = new Date().getTimezoneOffset() / -60
                        if (offset < 10) {
                            offset = "0" + offset
                        }
                        // format this in Go way
                        let timeStr = this.taskSnoozeDateInDialog + "T" + this.taskSnoozeTimeInDialog + ":00+" + offset + ":00"
                        tasksForSnooze.push({"id": task.id, "snoozeTo": timeStr})
                    }
                })
                // send task IDs to server
                let url = this.apiUrl + '/api/v1/project/' + this.projectIdSelected + '/task'
                axios.put(url, tasksForSnooze, this.authConfig())
                    .then((res) => {
                        this.getTasks(this.projectIdSelected)
                    })
                    .catch(((err) => {
                        if (err.response.status === 401) {
                            console.log('logged out')
                            this.$root.$emit('sessionExpired')
                        } else if (err.response.status === 400) {
                            console.log('empty result / wrong request')
                        } else {
                            console.log('something wrong')
                        }
                    }))
                // hide dialog
                this.snoozeTaskDialog = false
                // refresh task list
                this.getTasks(this.projectIdSelected)
                // TODO add snackbar with info for user
            },
            saveTask() {
                this.taskSaving = true
                // get selected tasks
                let data = [{
                    "id": this.taskIdInDialog,
                    "title": this.taskTitleInDialog,
                    "assignTo": this.taskAssignedInDialog,
                    "repeatUnit": this.taskRepeatUnitInDialog,
                    "repeatEvery": Number(this.taskRepeatEveryInDialog),
                }]
                // send task IDs to server
                let url = this.apiUrl + '/api/v1/project/' + this.projectIdSelected + '/task'
                axios.put(url, data, this.authConfig())
                    .then((res) => {
                        this.taskSaving = false
                        // hide dialog
                        this.editTaskDialog = false
                        // refresh task list
                        this.getTasks(this.projectIdSelected)
                        //TODO add snackbar with info for user
                    })
                    .catch(((err) => {
                        if (err.response.status === 401) {
                            console.log('logged out')
                            this.$root.$emit('sessionExpired')
                        } else if (err.response.status === 400) {
                            console.log('empty result / wrong request')
                        } else {
                            console.log('something wrong')
                        }
                    }))
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
                    .catch((err) => {
                        if (err.response.status === 401) {
                            console.log('logged out')
                            this.$root.$emit('sessionExpired')
                        } else if (err.response.status === 400) {
                            console.log('empty result / wrong request')
                        } else {
                            console.log('something wrong')
                        }
                    })
            },
            refreshTasks() {
                this.getTasks(this.projectIdSelected)
            },
            getTasks(projectId) {
                this.tasksLoading = true
                let url = this.apiUrl + '/api/v1/project/' + projectId + '/task'
                axios.get(url, this.authConfig())
                    .then((res) => {
                        // save before editing
                        this.tasks = res.data
                        // assign usernames to tasks
                        for (let i = 0; i < this.tasks.length; i++) {
                            console.log(this.tasks[i])
                            for (let j = 0; j < this.userList.length; j++) {
                                if (this.userList[j].id === this.tasks[i].assignedTo) {
                                    this.tasks[i].username = this.userList[j].username
                                }
                            }
                        }
                        this.tasks.forEach((task) => {
                            let username = "-"
                            this.userList.forEach((user) => {
                                if (user.id === task.assignedTo) {
                                    username = user.username
                                }
                            })
                            task
                        })
                        this.tasksLoading = false
                    })
                    .catch((err) => {
                        if (err.response.status === 401) {
                            console.log('logged out')
                            this.$root.$emit('sessionExpired')
                        } else if (err.response.status === 400) {
                            console.log('empty result / wrong request')
                        } else {
                            console.log('something wrong')
                        }
                    })
            },
            getUsers() {
                axios.get(this.apiUrl + '/api/v1/user', this.authConfig())
                    .then((res) => {
                        this.userList = res.data
                        this.userList.unshift({"id": 0, "username": "-"})
                        this.getProjects()
                    })
                    .catch((err) => {
                        if (err.response.status === 401) {
                            console.log('logged out')
                            this.$root.$emit('sessionExpired')
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
