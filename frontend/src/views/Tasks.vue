<template>
    <v-container fluid>
        <v-dialog
                v-model="dialog"
                max-width="290"
        >
            <v-card>
                <v-card-title class="headline">{{taskTitleInDialog}}</v-card-title>

                <v-card-text>
                    Run off table persian cat jump eat fish meeeeouw but more napping, more napping all the napping is
                    exhausting.
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
                            @click="dialog = false"
                    >
                        Submit
                    </v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
        <v-row>
            <v-col>
                <v-card
                        class="mx-auto"
                >
                    <v-toolbar
                            class="mb-2"
                            color="green"
                            dark
                            flat
                    >
                        <v-toolbar-title>TODO</v-toolbar-title>
                    </v-toolbar>
                </v-card>
            </v-col>
        </v-row>
        <v-row>
            <v-col md="3"></v-col>
            <v-col cols="12" xs="12" md="6">
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
                        <v-toolbar-title>{{projectTitle}}</v-toolbar-title>
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
    export default {
        name: 'tasks',
        data() {
            return {
                projectTitle: "Cool project name placeholder",
                newTaskTitle: "",
                tasks: [],
                dialog: false,
                taskTitleInDialog: "",
            }
        },
        mounted() {
        },
        methods: {
            addTask() {
                if (this.newTaskTitle !== "") {
                    this.tasks.push({"title": this.newTaskTitle, "selected": false,});
                    this.newTaskTitle = ""
                }
            },
            deleteTasks() {
                this.tasks = this.tasks.filter(function (task) {
                    return !task.selected
                })
            },
            showDialog(task) {
                this.taskTitleInDialog = task.title
                this.tasks.selected = false
                this.dialog = true
            },
        },
    }
</script>
