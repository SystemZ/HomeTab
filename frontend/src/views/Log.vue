<template>
    <v-container fluid>
        <v-card
                class="mx-auto"
                max-width="600"
        >
            <v-toolbar
                    :color="toolbarPrimary"
                    :dark="componentDark"
            >
                <v-toolbar-title>Last 100 sessions</v-toolbar-title>
            </v-toolbar>

            <v-list dense two-line>
                <template v-for="(counter, key) in counters">
                    <v-list-item
                            :key="counter.id"
                    >
                        <v-list-item-content>
                            <v-list-item-title>
                                <span class="subtitle-1">{{counter.title}} </span>
                                <v-chip
                                        :color="chipPrimary"
                                        :dark="componentDark"
                                        pill
                                        small>
                                    {{counter.tag}}
                                </v-chip>
                            </v-list-item-title>
                            <v-list-item-subtitle class="text--primary">
                                {{counter.username}} ended counting @ {{counter.endedAt | prettyTimeDate}}
                            </v-list-item-subtitle>
                            <v-list-item-subtitle>
                                Duration: {{counter.duration}}
                            </v-list-item-subtitle>
                        </v-list-item-content>
                    </v-list-item>
                    <v-divider
                            :key="counter.id"
                            v-if="counters.length !== key+1"
                    ></v-divider>
                </template>
            </v-list>
        </v-card>
    </v-container>
</template>

<script>
    import axios from "axios";

    export default {
        name: 'log',
        data() {
            return {
                logLoading: true,
                counters: [
                    {
                        id: 1,
                        title: "My Time at Portia",
                        tag: "PC",
                        username: "Aylin",
                        endedAt: "02:25 09/05/2020",
                        duration: "2h 12m 22s"
                    },
                    {
                        id: 2,
                        title: "The Legend of Zelda: Breath of the Wild",
                        tag: "Switch",
                        username: "System",
                        endedAt: "12:34 09/05/2020",
                        duration: "15m 08s"
                    }],
            }
        },
        mounted() {
        },
        methods: {
            authConfig() {
                return {headers: {Authorization: "Bearer " + localStorage.getItem(this.lsToken)}}
            },
            getEvents() {
                this.logLoading = true
                axios.get(this.apiUrl + "/api/v1/event", this.authConfig())
                    .then((res) => {
                        this.logLoading = false
                        this.counters = res.data
                    })
                    .catch((err) => {
                        if (err.response.status === 401) {
                            this.$root.$emit("sessionExpired")
                        } else {
                            console.log("something wrong")
                        }
                    })
            },
        }
    }
</script>
