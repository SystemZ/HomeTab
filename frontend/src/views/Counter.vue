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
                            dark
                            color="red"
                            @click="toggleCounter(false)"
                            v-if="counter.running"
                    >
                        Stop
                    </v-btn>
                    <v-btn
                            class="mt-3"
                            x-large
                            dark
                            color="green"
                            @click="toggleCounter(true)"
                            v-else
                    >
                        Start
                    </v-btn>
                </div>
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
                    tags: "PC",
                    running: false,
                }
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
