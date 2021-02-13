<template>
    <v-container fluid>
        <v-card
                max-width="600"
                class="mx-auto"
        >
            <v-toolbar
                    :color="toolbarPrimary"
                    :dark="componentDark"
            >
                <v-toolbar-title>Last 7 days</v-toolbar-title>
            </v-toolbar>

            <v-list two-line dense>
                <template v-for="(event, key) in events">
                    <v-list-item
                            :key="event.id"
                    >
                        <v-list-item-content>
                            <v-list-item-title v-text="event.subject"></v-list-item-title>
                            <v-list-item-subtitle>{{event.username}} @ {{event.createdAt | prettyTimeDate}}</v-list-item-subtitle>
                        </v-list-item-content>
                    </v-list-item>
                    <v-divider
                            :key="event.id"
                            v-if="events.length !== key+1"
                    ></v-divider>
                </template>
            </v-list>
        </v-card>
    </v-container>
</template>

<script>
    import axios from "axios";

    export default {
        name: 'events',
        data() {
            return {
                eventsLoading: true,
                events: [{}],
            }
        },
        mounted() {
            this.getEvents()
        },
        methods: {
            authConfig() {
                return {headers: {Authorization: "Bearer " + localStorage.getItem(this.lsToken)}}
            },
            getEvents() {
                this.eventsLoading = true
                axios.get(this.apiUrl + "/api/v1/event", this.authConfig())
                    .then((res) => {
                        this.eventsLoading = false
                        this.events = res.data
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
