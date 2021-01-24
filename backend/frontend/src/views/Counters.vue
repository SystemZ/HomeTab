<template>
    <v-container fluid>

        <v-dialog
                v-model="newCounterDialogShow"
                max-width="290"
        >
            <v-card>
                <v-card-title>
                    <span class="headline">Add new counter</span>
                </v-card-title>
                <v-card-text>
                    <v-container>
                        <v-row>
                            <v-col cols="12">
                                <v-text-field :disabled="newCounterInProgress"
                                              v-model="newCounterName"
                                              label="Name*"
                                              required></v-text-field>
                            </v-col>
                        </v-row>
                        <v-row>
                            <v-col cols="12">
                                <v-text-field
                                        :disabled="newCounterInProgress"
                                        v-model="newCounterTag"
                                        label="Tag*"
                                        required></v-text-field>
                            </v-col>
                        </v-row>
                    </v-container>
                    <small>*indicates required field</small>
                </v-card-text>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn :color="btnSecondary" text @click="newCounterDialogShow = false">Close</v-btn>

                    <v-progress-circular v-if="newCounterInProgress" indeterminate color="green"
                                         class="ml-5 ml-5"></v-progress-circular>
                    <v-btn v-else :dark="btnDark" :color="btnPrimary" @click="addCounter">Add</v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>

        <v-btn class="ma-5" :dark="btnDark" :color="btnPrimary" @click="newCounterDialogShow = true">
            Add new counter
        </v-btn>

        <v-card class="ma-5">
            <v-card-title>
                All counters
                <v-spacer></v-spacer>
                <v-row>
                    <v-switch v-model="searchPro" label="Pro" :color="btnPrimary" class="mr-5"></v-switch>
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

        <v-snackbar v-model="snackbarShow">
            {{ snackbarText }}
            <v-btn color="pink" text @click="snackbarShow = false">
                Close
            </v-btn>
        </v-snackbar>

    </v-container>
</template>

<script>
  import axios from 'axios'

  export default {
    name: 'counters',
    data () {
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
          {text: 'Time', value: 'secondsF'},
          // {text: 'Created at', value: 'createdAt'},
        ],
        counters: [],
        prevItemsPerPage: 0,
        prevPage: 0,
        snackbarShow: false,
        snackbarText: '',
        newCounterInProgress: false,
        newCounterDialogShow: false,
        newCounterName: '',
        newCounterTag: 'PC',
      }
    },
    methods: {
      authConfig () {
        return {headers: {Authorization: 'Bearer ' + localStorage.getItem(this.lsToken)}}
      },
      toCounter (item) {
        this.$router.push({name: 'counter', params: {id: item.id}})
      },
      getCounters (pagination, resetPagination) {
        let vm = this
        vm.countersLoading = true
        let lastId = 0
        let queryType = 'next'
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
            queryType = 'prev'
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

        let searchQuery = '%' + vm.search + '%'
        if (vm.searchPro) {
          searchQuery = vm.search
        }

        let rawUrl = vm.apiUrl + '/api/v1/counter-page?limit=' + itemsPerPage + '&' + queryType + 'Id=' + lastId
        let url = encodeURI(rawUrl)
        axios.post(url, {q: searchQuery}, vm.authConfig())
          .then((res) => {
            vm.countersLoading = false
            vm.counters = res.data.counters
            vm.totalCounters = res.data.pagination.allRecords
          })
          .catch(function (err) {
            if (err.response.status === 401) {
              console.log('logged out')
              vm.$root.$emit('sessionExpired')
            } else if (err.response.status === 400) {
              console.log('empty result / wrong request')
              vm.countersLoading = false
              // make table empty as backend result is empty
              vm.counters = []
              vm.totalCounters = 0
            } else {
              console.log('something wrong')
            }
          })
      },
      doSearch () {
        clearTimeout(this.searchTimeout)
        this.searchTimeout = setTimeout(this.getCounters, 500, undefined, true)
      },
      addCounter () {
        // prevent user fiddling when we waiting for server's answer
        this.newCounterInProgress = true
        let data = {'name': this.newCounterName, 'tag': this.newCounterTag}
        axios.post(this.apiUrl + '/api/v1/counter', data, this.authConfig())
          .then((res) => {
            // empty form for faster adding many records and unlock it
            this.newCounterName = ''
            this.newCounterTag = 'PC'
            this.newCounterInProgress = false
            // notify user
            this.snackbarText = 'Counter added'
            this.snackbarShow = true
            // hide form
            this.newCounterDialogShow = false
            // load all counters again
            this.getCounters(undefined, true)
          })
          .catch((err) => {
            if (err.response.status === 401) {
              console.log('logged out')
              this.$root.$emit('sessionExpired')
            } else if (err.response.status === 400) {
              console.log('empty result / wrong request')
              this.snackbarText = 'Both fields should not be empty'
              this.snackbarShow = true
            } else {
              console.log('something wrong')
            }
          })
      },
    }
  }
</script>
