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
                ></v-text-field>
            </v-col>
            <v-col>
                <v-btn
                        block
                        color="success"
                        v-if="$vuetify.breakpoint.smAndDown"
                >
                    Add
                </v-btn>
                <v-btn
                        color="success"
                        large
                        height="48px"
                        v-else
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
                        <v-toolbar-title>Project name</v-toolbar-title>
                        <v-spacer></v-spacer>
                        <v-btn icon
                        >
                            <v-icon>mdi-check-bold</v-icon>
                        </v-btn>
                        <v-btn icon
                        >
                            <v-icon>mdi-alarm-snooze</v-icon>
                        </v-btn>
                        <v-btn icon
                        >
                            <v-icon>mdi-delete</v-icon>
                        </v-btn>
                    </v-toolbar>
                    <v-list
                            subheader
                            two-line
                            flat
                    >
                        <v-subheader>Tasks</v-subheader>
                        <!-- two items to show multiple choice -->
                        <v-list-item-group multiple>
                            <v-list-item>
                                <template v-slot:default="{ active, toggle }">
                                    <v-list-item-action>
                                        <v-checkbox
                                                v-model="active"
                                                color="success"
                                                @click="toggle"
                                        ></v-checkbox>
                                    </v-list-item-action>

                                    <v-list-item-content>
                                        <v-list-item-title>Take out the trash</v-list-item-title>
                                    </v-list-item-content>
                                </template>
                            </v-list-item>

                            <v-list-item>
                                <template v-slot:default="{ active, toggle }">
                                    <v-list-item-action>
                                        <v-checkbox
                                                v-model="active"
                                                color="success"
                                                @click="toggle"
                                        ></v-checkbox>
                                    </v-list-item-action>

                                    <v-list-item-content>
                                        <v-list-item-title>Buy batteries</v-list-item-title>
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
    }
</script>
