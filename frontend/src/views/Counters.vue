<template>
    <div>
        <v-card class="ma-5">
            <v-card-title>
                Counters
                <v-spacer></v-spacer>
                <v-text-field
                        v-model="search"
                        append-icon="mdi-magnify"
                        label="Search"
                        single-line
                        hide-details
                >
                </v-text-field>
            </v-card-title>
            <v-data-table
                    :headers="headers"
                    :items="counters"
                    :search="search"
                    :loading="countersLoading"
                    @click:row="toCounter"
            >
                <template v-slot:item.createdAt="{ item }">
                    {{item.createdAt | prettyTimeDate }}
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
    </div>
</template>

<script>
    import axios from "axios";

    export default {
        name: 'counters',
        data() {
            return {
                search: '',
                countersLoading: true,
                headers: [
                    {
                        text: 'ID',
                        align: 'left',
                        sortable: true,
                        value: 'id',
                    },
                    {text: 'Title', value: 'name'},
                    {text: 'Tags', value: 'tags'},
                    {text: 'Time', value: 'seconds'},
                    // {text: 'Created at', value: 'createdAt'},
                ],
                counters: [],
            }
        },
        mounted() {
            this.getCounters()
        },
        methods: {
            authConfig() {
                return {headers: {Authorization: "Bearer " + localStorage.getItem(this.lsToken)}}
            },
            toCounter(item) {
                this.$router.push({name: 'note', params: {id: item.id}})
            },
            getCounters() {
                let vm = this
                vm.countersLoading = true
                axios.get(vm.apiUrl + "/api/v1/counter", vm.authConfig())
                    .then((res) => {
                        vm.countersLoading = false
                        vm.counters = res.data
                    })
                    .catch(function (err) {
                        if (err.response.status === 401) {
                            console.log("logged out")
                            vm.$root.$emit("sessionExpired")
                        } else {
                            console.log("something wrong")
                        }
                    })
            }
        }
    }
</script>
