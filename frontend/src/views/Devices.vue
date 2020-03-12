<template>
    <v-container fill-height>
        <v-row>
            <v-col>
                <v-data-iterator
                        :items="devices"
                        :items-per-page.sync="devicesPerPage"
                        hide-default-footer
                >
                    <template v-slot:default>
                        <v-row>
                            <v-col
                                    v-for="device in devices"
                                    :key="device.id"
                                    cols="12"
                                    sm="6"
                                    md="6"
                            >
                                <v-card>
                                    <v-card-title class="subheading font-weight-bold">{{ device.name }} -
                                        {{device.username}}
                                    </v-card-title>

                                    <v-divider></v-divider>

                                    <v-list dense>
                                        <!-- FIXME in dire need of nested loop for this mess below -->
                                        <v-list-item>
                                            <v-list-item-content>Battery</v-list-item-content>
                                            <v-list-item-content class="align-end">{{ device.battery }}
                                            </v-list-item-content>
                                        </v-list-item>
                                        <v-list-item>
                                            <v-list-item-content>Display</v-list-item-content>
                                            <v-list-item-content class="align-end">{{ device.displayOn }}
                                            </v-list-item-content>
                                        </v-list-item>
                                        <v-list-item>
                                            <v-list-item-content>Last display ON</v-list-item-content>
                                            <v-list-item-content class="align-end">{{ device.lastDisplayOn }}
                                            </v-list-item-content>
                                        </v-list-item>
                                        <v-list-item>
                                            <v-list-item-content>Last display OFF</v-list-item-content>
                                            <v-list-item-content class="align-end">{{ device.lastDisplayOff }}
                                            </v-list-item-content>
                                        </v-list-item>
                                        <v-list-item>
                                            <v-list-item-content>Last track</v-list-item-content>
                                            <v-list-item-content class="align-end">{{ device.lastTrack }}
                                            </v-list-item-content>
                                        </v-list-item>
                                    </v-list>
                                </v-card>
                            </v-col>
                        </v-row>
                    </template>
                </v-data-iterator>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>

    export default {
        name: 'devices',
        data() {
            return {
                devicesPerPage: 2,
                devices: [
                    {
                        id: 1,
                        name: "Xiaomi",
                        username: "P",
                        battery: "50%",
                        displayOn: false,
                        lastDisplayOn: "15:12:28 12.03.2020",
                        lastDisplayOff: "15:13:10 12.03.2020",
                        lastTrack: "The Phantom Of The Opera - Album Version - Nightwish",
                    },
                    {
                        id: 2,
                        name: 'Samsung',
                        username: 'S',
                        battery: "50%",
                        displayOn: false,
                        lastDisplayOn: "15:12:28 12.03.2020",
                        lastDisplayOff: "15:13:10 12.03.2020",
                        lastTrack: "Centuries - Fall Out Boy",
                    },
                ],
            }
        },
        methods: {
            authConfig() {
                return {headers: {Authorization: "Bearer " + localStorage.getItem(this.lsToken)}}
            },
        }
    }
</script>
