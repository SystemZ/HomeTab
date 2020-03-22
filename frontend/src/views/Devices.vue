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
                                        <v-list-item>
                                            <v-list-item-content>Battery</v-list-item-content>
                                            <v-list-item-content class="align-end">{{ device.battery }}%
                                            </v-list-item-content>
                                        </v-list-item>
                                        <v-list-item>
                                            <v-list-item-content>Display</v-list-item-content>
                                            <v-list-item-content class="align-end">{{ device.displayState }}
                                            </v-list-item-content>
                                        </v-list-item>
                                        <v-list-item>
                                            <v-list-item-content>Last display ON</v-list-item-content>
                                            <v-list-item-content class="align-end">
                                                {{ device.displayLastOn | prettyTimeDate }}
                                            </v-list-item-content>
                                        </v-list-item>
                                        <v-list-item>
                                            <v-list-item-content>Last display OFF</v-list-item-content>
                                            <v-list-item-content class="align-end">
                                                {{ device.displayLastOff | prettyTimeDate }}
                                            </v-list-item-content>
                                        </v-list-item>
                                        <v-list-item>
                                            <v-list-item-content>Last track</v-list-item-content>
                                            <v-list-item-content class="align-end">
                                                {{ device.musicTrack }} - {{device.musicArtist }}
                                            </v-list-item-content>
                                        </v-list-item>
                                        <v-list-item>
                                            <v-list-item-content>Last track played</v-list-item-content>
                                            <v-list-item-content class="align-end">
                                                {{ device.musicLastPlayed | prettyTimeDate }}
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
    import axios from "axios";

    export default {
        name: 'devices',
        data() {
            return {
                devicesPerPage: 2,
                devices: [],
            }
        },
        mounted() {
            this.getDevices()
        },
        methods: {
            authConfig() {
                return {headers: {Authorization: "Bearer " + localStorage.getItem(this.lsToken)}}
            },
            getDevices() {
                axios.get(this.apiUrl + "/api/v1/device", this.authConfig())
                    .then((res) => {
                        this.devices = res.data
                    })
                    .catch((err) => {
                        if (err.response.status === 401) {
                            console.log("logged out")
                            this.$root.$emit("sessionExpired")
                        } else {
                            console.log("something wrong")
                        }
                    })
            },
        }
    }
</script>
