<template>
    <v-container fluid>

        Task done stream from last 7 days
        <div v-for="event in events" :key="event.id">
            {{event.createdAt | prettyTimeDate}}
            {{event.subject}}
            {{event.username}}
        </div>

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
