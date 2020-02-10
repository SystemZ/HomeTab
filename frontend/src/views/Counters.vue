<template>
    <div>
        <v-card class="ma-5">
            <v-card-title>
                Counters
                <v-spacer></v-spacer>
                <v-row>
                    <v-switch v-model="searchPro" label="Pro" color="green" class="mr-5"></v-switch>
                    <v-text-field
                            v-model="search"
                            append-icon="mdi-magnify"
                            label="Search"
                            single-line
                            hide-details
                            @keydown="doSearch"
                    >
                    </v-text-field>
                </v-row>
            </v-card-title>
            <v-data-table
                    :headers="headers"
                    :items="counters"
                    :search="search"
                    :loading="countersLoading"
                    :options.sync="tableOptions"
                    @pagination="getCounters"
                    :server-items-length="totalCounters"
                    :footer-props="{'disable-pagination': countersLoading}"
                    disable-sort
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
                searchPro: false,
                searchTimeout: null,
                //
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
            getCounters(pagination, resetPagination) {
                let vm = this
                vm.countersLoading = true
                let lastId = 0
                let queryType = "next"
                if (vm.counters.length > 0) {
                    if (pagination !== undefined && pagination.page > vm.prevPage) {
                        // next page button
                        lastId = vm.counters[vm.counters.length - 1].id
                    } else if (resetPagination) {
                        // search box handling
                        // we don't want to skip records with ID bigger than 1
                        lastId = 0
                        vm.tableOptions.page = 1
                    } else {
                        // prev page button
                        lastId = vm.counters[0].id
                        queryType = "prev"
                    }
                }

                let itemsPerPage = vm.prevItemsPerPage
                if (pagination !== undefined) {
                    vm.prevPage = pagination.page
                    if (pagination.itemsPerPage !== vm.prevItemsPerPage) {
                        lastId = 0
                    }
                    vm.prevItemsPerPage = pagination.itemsPerPage
                    itemsPerPage = pagination.itemsPerPage
                }

                let searchQuery = "%" + vm.search + "%"
                if (vm.searchPro) {
                    searchQuery = vm.search
                }

                let rawUrl = vm.apiUrl + "/api/v1/counter-page?limit=" + itemsPerPage + "&" + queryType + "Id=" + lastId
                let url = encodeURI(rawUrl)
                axios.post(url, {q: searchQuery}, vm.authConfig())
                    .then((res) => {
                        vm.countersLoading = false
                        vm.counters = res.data.counters
                        vm.totalCounters = res.data.pagination.allRecords
                    })
                    .catch(function (err) {
                        if (err.response.status === 401) {
                            console.log("logged out")
                            vm.$root.$emit("sessionExpired")
                        } else if (err.response.status === 400) {
                            console.log("empty result / wrong request")
                            vm.countersLoading = false
                            // make table empty as backend result is empty
                            vm.counters = []
                            vm.totalCounters = 0
                        } else {
                            console.log("something wrong")
                        }
                    })
            },
            doSearch() {
                clearTimeout(this.searchTimeout)
                this.searchTimeout = setTimeout(this.getCounters, 500, undefined, true)
            },
        }
    }
</script>
