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
            <v-col cols="12" xs="12" md="4">
                <v-text-field
                        placeholder="Buy a yacht"
                        solo
                        clearable
                        v-model="newTaskTitle"
                        @keydown.enter.native="addTask"
                >
                </v-text-field>
            </v-col>
            <v-col>
                <v-btn
                        block
                        color="success"
                        v-if="$vuetify.breakpoint.smAndDown"
                        @click.native="addTask"
                >
                    Add
                </v-btn>
                <v-btn
                        color="success"
                        large
                        height="48px"
                        v-else
                        @click.native="addTask"
                >
                    Add
                </v-btn>
            </v-col>
        </v-row>
        <v-row></v-row>
        <v-row>
            <v-col>
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
                            one-line
                    >
                        <v-subheader>Tasks</v-subheader>
                        <v-list-item-group multiple>
                            <!--  TODO fix checkbox after deleting tasks -->
                            <v-list-item
                                    v-for="(task, i) in tasks"
                                    :key="i"
                                    inactive
                            >
                                <template v-slot:default>
                                    <v-list-item-action>
                                        <v-checkbox
                                                v-model="task.selected"
                                                color="green darken-2"
                                        ></v-checkbox>
                                    </v-list-item-action>

                                    <v-list-item-content @click.stop="showDialog(task)">
                                        <v-list-item-title v-text="task.title"></v-list-item-title>
                                    </v-list-item-content>
                                </template>
                            </v-list-item>
                        </v-list-item-group>
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
                    this.tasks.push({"title": this.newTaskTitle});
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
                this.dialog = true
            },
        },
    }
</script>
