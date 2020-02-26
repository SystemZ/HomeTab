<template>
    <v-container fill-height>
        <v-row class="text-center">
            <v-col cols="12">
                <span class="display-1 mr-1">{{counter.title}}</span>
                <span class="overline">{{counter.tags}}</span>
                <div v-if="!counterLoading">
                    <v-btn
                            class="mt-3"
                            x-large
                            :dark="btnDark"
                            :color="btnSecondary"
                            @click="toggleCounter(false)"
                            v-if="counter.running"
                    >
                        Stop
                    </v-btn>
                    <v-btn
                            class="mt-3"
                            x-large
                            :dark="btnDark"
                            :color="btnPrimary"
                            @click="toggleCounter(true)"
                            v-else
                    >
                        Start
                    </v-btn>
                </div>
            </v-col>
        </v-row>
        <v-row>
            <v-col>
                <v-card>
                    <v-card-title>
                        Last 100 sessions
                    </v-card-title>
                    <v-data-table
                            :headers="headers"
                            :items="sessions"
                            :loading="counterLoading"
                            :footer-props="{'disable-pagination': counterLoading}"
                            disable-sort
                            :items-per-page="5"
                    >
                        <template v-slot:item.start="{ item }">
                            {{item.start | prettyTimeDate }}
                        </template>
                        <template v-slot:item.end="{ item }">
                            {{item.end | prettyTimeDate }}
                        </template>
                        <template v-slot:progress>
                            <v-progress-linear
                                    indeterminate
                                    :height="2"
                                    color="green"
                            ></v-progress-linear>
                        </template>
                    </v-data-table>
                </v-card>
            </v-col>
        </v-row>
        <v-row>
            <v-col>
                <v-data-iterator
                        :items="records"
                        :items-per-page.sync="recordsPerPage"
                        hide-default-footer
                >
                    <template v-slot:header>
                        <v-toolbar
                                class="mb-2"
                                color="green"
                                dark
                                flat
                        >
                            <v-toolbar-title>Stats</v-toolbar-title>
                        </v-toolbar>
                    </template>

                    <template v-slot:default="props">
                        <v-row>
                            <v-col
                                    v-for="item in props.items"
                                    :key="item.name"
                                    cols="12"
                                    sm="6"
                                    md="4"
                            >
                                <v-card>
                                    <v-card-title class="subheading font-weight-bold">{{ item.name }}</v-card-title>

                                    <v-divider></v-divider>

                                    <v-list dense>
                                        <v-list-item>
                                            <v-list-item-content>Time spent:</v-list-item-content>
                                            <v-list-item-content class="align-end">{{ item.time }}</v-list-item-content>
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
        name: 'counter',
        data() {
            return {
                counterLoading: true,
                counter: {
                    id: 0,
                    title: "...",
                    tags: [],
                    running: false,
                },
                headers: [
                    {
                        text: 'Start',
                        align: 'left',
                        value: 'start',
                    }, // date & time
                    {text: 'End', value: 'end'}, // date & time
                    {text: 'Session', value: 'durationSF'}, // how long game was played/current session is taking
                    // {text: 'Created at', value: 'createdAt'},
                ],
                sessions: [
                    {
                        start: "",
                        stop: "",
                        durationSF: ""
                    },
                ],
                recordsPerPage: 3,
                records: [
                    {
                        name: "Last 7d",
                        time: "0h 0m 0s",
                    },
                    {
                        name: 'Last 30d',
                        time: "0h 0m 0s",
                    },
                    {
                        name: 'All time',
                        time: "0h 0m 0s",
                    },
                ],
            }
        },
        watch: {},
        mounted() {
            this.counter.id = this.$route.params.id
            this.getCounter()
        },
        methods: {
            authConfig() {
                return {headers: {Authorization: "Bearer " + localStorage.getItem(this.lsToken)}}
            },
            toggleCounter(start) {
                let vm = this
                let verb = "stop"
                if (start) {
                    verb = "start"
                }
                //vm.notesLoading = true
                axios.put(vm.apiUrl + "/api/v1/counter/" + vm.counter.id + "/" + verb, {}, vm.authConfig())
                    .then((res) => {
                        vm.getCounter()
                    })
                    .catch(function (err) {
                        if (err.response.status === 401) {
                            console.log("logged out")
                            vm.$root.$emit("sessionExpired")
                        } else {
                            console.log("something wrong")
                        }
                    })
            },
            getCounter() {
                let vm = this
                vm.counterLoading = true
                axios.get(vm.apiUrl + "/api/v1/counter/" + vm.counter.id + "/info", vm.authConfig())
                    .then((res) => {
                        vm.counter.title = res.data.name
                        vm.counter.running = res.data.inProgress
                        if (res.data.tags !== null) {
                            vm.counter.tags = res.data.tags.toString()
                        }
                        if (res.data.sessions !== null) {
                            vm.sessions = res.data.sessions
                        }
                        vm.records = []
                        vm.records[0] = {"name": "Last 7d", "time": res.data.stats.secondsD7F}
                        vm.records[1] = {"name": "Last 30d", "time": res.data.stats.secondsD30F}
                        vm.records[2] = {"name": "All time", "time": res.data.stats.secondsAllF}
                        vm.counterLoading = false
                    })
                    .catch(function (err) {
                        if (err.response.status === 401) {
                            console.log("logged out")
                            vm.$root.$emit("sessionExpired")
                        } else {
                            console.log("something wrong")
                        }
                    })
            },
        }
    }
</script>
