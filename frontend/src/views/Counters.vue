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
                    :options.sync="tableOptions"
                    @pagination="getCounters"
                    :server-items-length="totalCounters"
            >
                <!--@click:row="toCounter"-->
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
                tableOptions: {},
                totalCounters: 0,
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
                prevItemsPerPage: 0,
                prevPage: 0,
            }
        },
        methods: {
            authConfig() {
                return {headers: {Authorization: "Bearer " + localStorage.getItem(this.lsToken)}}
            },
            toCounter(item) {
                this.$router.push({name: 'note', params: {id: item.id}})
            },
            getCounters(pagination) {
                let vm = this
                vm.countersLoading = true
                let lastId = 0
                let queryType = "next"
                if (vm.counters.length > 0) {
                    if (pagination.page > vm.prevPage) {
                        lastId = vm.counters[vm.counters.length - 1].id
                    } else {
                        lastId = vm.counters[0].id
                        queryType = "prev"
                    }
                }
                vm.prevPage = pagination.page
                if (pagination.itemsPerPage !== vm.prevItemsPerPage) {
                    lastId = 0
                }
                vm.prevItemsPerPage = pagination.itemsPerPage
                axios.get(vm.apiUrl + "/api/v1/counter-page?limit=" + pagination.itemsPerPage + "&" + queryType + "Id=" + lastId, vm.authConfig())
                    .then((res) => {
                        vm.countersLoading = false
                        vm.counters = res.data.counters
                        vm.totalCounters = res.data.pagination.allRecords
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
