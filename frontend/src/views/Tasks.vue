<template>
    <v-container fluid>
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
                        v-model="task"
                >
                    {{task}}
                </v-text-field>
            </v-col>
            <v-col>
                <v-btn
                        block
                        color="success"
                        v-if="$vuetify.breakpoint.smAndDown"
                        @click="addTask"
                >
                    Add
                </v-btn>
                <v-btn
                        color="success"
                        large
                        height="48px"
                        v-else
                        @click="addTask"
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
                        <v-toolbar-title>{{project}}</v-toolbar-title>
                        <v-spacer></v-spacer>
                        <v-btn icon>
                            <v-icon>mdi-check-bold</v-icon>
                        </v-btn>
                        <v-btn icon>
                            <v-icon>mdi-alarm-snooze</v-icon>
                        </v-btn>
                        <v-btn icon>
                            <v-icon>mdi-delete</v-icon>
                        </v-btn>
                    </v-toolbar>
                    <v-list
                            subheader
                            two-line
                    >
                        <v-subheader>Tasks</v-subheader>
                        <v-list-item-group multiple>
                            <v-list-item
                                    v-for="(item, i) in items"
                                    :key="i"
                                    active-class="green--text text--darken-4"
                            >
                                <template v-slot:default="{ active, toggle }">
                                    <v-list-item-action>
                                        <v-checkbox
                                                v-model="active"
                                                color="green darken-2"
                                                @click="toggle"
                                        ></v-checkbox>
                                    </v-list-item-action>

                                    <v-list-item-content>
                                        <v-list-item-title v-text="item.title"></v-list-item-title>
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
                project: "Cool project name placeholder",
                task: "",
                items: [],
            }
        },
        mounted() {
        },
        methods: {
            addTask() {
                this.items.push({"title": this.task})
            },
        },
    }
</script>
