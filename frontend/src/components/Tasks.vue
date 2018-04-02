<template>
  <v-container fluid> <!-- fluid? -->
    <v-slide-y-transition mode="out-in">
      <v-layout row wrap>
        <v-flex xs12 lg4>
          <h2>Tools</h2>

          <v-btn v-if="refreshing" color="red" @click="getTasks()">
            <v-icon color="white">
              fa fa-refresh
            </v-icon>
          </v-btn>
          <v-btn v-else @click="getTasks()">
            <v-icon color="green">
              fa fa-refresh
            </v-icon>
          </v-btn>
          <h2 class="mt-2">Filters</h2>
          <!-- all start -->
          <v-btn v-if="showTasks == 'todo'" @click="setTasks('all')">
            <v-icon class="mr-1">
              fa fa-tasks
            </v-icon>
            TODO
          </v-btn>
          <!-- all end -->
          <!-- to do start -->
          <v-btn v-if="showTasks == 'all'" @click="setTasks('todo')">
            <v-icon class="mr-1">
              fa fa-check
            </v-icon>
            ALL
          </v-btn>
          <!-- to do end -->
          <!-- gmail start -->
          <v-btn v-if="show.gmail == true" @click="show.gmail = false">
            <v-icon class="mr-1">
              fa fa-at
            </v-icon>
            <v-icon color="green">
              fa fa-check
            </v-icon>
          </v-btn>
          <v-btn v-if="show.gmail == false" @click="show.gmail = true">
            <v-icon class="mr-1">
              fa fa-at
            </v-icon>
            <v-icon color="red">
              fa fa-times
            </v-icon>
          </v-btn>
          <!-- gmail end -->
          <!-- gitlab start -->
          <v-btn v-if="show.gitlab == true" @click="show.gitlab = false">
            <v-icon class="mr-1" color="orange darken-1">
              fa fa-gitlab
            </v-icon>
            <v-icon color="green">
              fa fa-check
            </v-icon>
          </v-btn>
          <v-btn v-if="show.gitlab == false" @click="show.gitlab = true">
            <v-icon class="mr-1" color="orange darken-1">
              fa fa-gitlab
            </v-icon>
            <v-icon color="red">
              fa fa-times
            </v-icon>
          </v-btn>
          <!-- gitlab end -->
          <!-- github start -->
          <v-btn v-if="show.github == true" @click="show.github = false">
            <v-icon class="mr-1">
              fa fa-github
            </v-icon>
            <v-icon color="green">
              fa fa-check
            </v-icon>
          </v-btn>
          <v-btn v-if="show.github == false" @click="show.github = true">
            <v-icon class="mr-1">
              fa fa-github
            </v-icon>
            <v-icon color="red">
              fa fa-times
            </v-icon>
          </v-btn>
          <!-- github end -->

          <h2>Batch edit</h2>

          <v-btn v-if="checkboxesVisible" @click="hideCheckboxes()">
            <v-icon color="green">
              fa fa-check
            </v-icon>
          </v-btn>
          <v-btn v-else @click="showCheckboxes()">
            <v-icon color="red darken-4">
              fa fa-times
            </v-icon>
          </v-btn>

          <div v-if="checkboxesVisible">

            <v-btn v-if="batchDelayOn" @click="batchDelayOn = false">
              <v-icon color="green">
                fa fa-clock-o
              </v-icon>
            </v-btn>
            <v-btn v-else @click="batchDelayOn = true">
              <v-icon color="red">
                fa fa-clock-o
              </v-icon>
            </v-btn>

            <v-flex v-if="batchDelayOn" xs12 sm4>
              <v-select
                :items="[1,3,5,7,30]"
                v-model="batchDelay"
                label="Delay by days"
              ></v-select>
            </v-flex>

            <v-btn v-if="selectedTasks.length > 0" @click="batchApply()">
              Apply
            </v-btn>

          </div>

        </v-flex>

        <v-flex xs12 lg8>
          <v-card
            class="mp-2"
            v-for="task in tasks"
            :key="task.id"
            v-if="show[task.type] == true && task.hide !== true"
            @click="redirectToOrigin(task.id)"
          >
            <v-list>
              <v-list-tile avatar>
                <v-list-tile-action class="mr-3">
                  <v-icon v-if="task.type === 'github'">fa fa-2x fa-fw fa-github</v-icon>
                  <v-icon v-if="task.type === 'gitlab'">fa fa-2x fa-fw fa-gitlab</v-icon>
                  <v-icon v-if="task.type === 'gmail'">fa fa-2x fa-fw fa-at</v-icon>
                </v-list-tile-action>
                <v-list-tile-action v-if="checkboxesVisible">
                  <v-checkbox
                    :value="task.id"
                    v-model="selectedTasks"
                  ></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-action v-else class="mr-3">
                  <v-btn v-if="task.done" color="green" fab outline small dark @click="task.done = false">
                    <v-icon>check</v-icon>
                  </v-btn>
                  <v-btn v-else color="grey" outline fab small dark @click="task.done = true">
                    <v-icon>outline</v-icon>
                  </v-btn>
                </v-list-tile-action>
                <v-list-tile-content>
                  <v-list-tile-title>
                    {{task.title}}
                    <small>{{task.projectName}} / #{{task.projectTaskId}}</small>
                  </v-list-tile-title>
                </v-list-tile-content>
                <v-list-tile-action>
                  <v-btn small flat @click="redirectToOrigin(task.id)">
                    <v-icon>
                      fa fa-fw fa-external-link
                    </v-icon>
                  </v-btn>
                </v-list-tile-action>
                <!--</v-list-tile-action>-->
                <!--<v-list-tile-avatar>-->
                <!--<img src="http://via.placeholder.com/40x40"/>-->
                <!--</v-list-tile-avatar>-->
              </v-list-tile>
            </v-list>

          </v-card>
        </v-flex>

      </v-layout>
    </v-slide-y-transition>
  </v-container>
</template>

<script>
  import axios from 'axios'

  export default {
    data () {
      return {
        //batch edit
        checkboxesVisible: false,
        selectedTasks: [],
        batchDelayOn: false,
        batchDelay: 1,
        //
        refreshing: true,
        tasks: [],
        showTasks: 'todo',
        allTasksUrl: 'tasks/all/1',
        //todoTasksUrl: 'tasks/todo/1',
        todoTasksUrl: 'tasks/focus/1',
        show: {
          'github': true,
          'gitlab': true,
          'gmail': true
        }
      }
    },
    // https://stackoverflow.com/a/34143409/1351857
    // computed: {
    //   filteredTasks: function () {
    //     var self = this;
    //     return self.tasks.filter(function (task) {
    //       return task.indexOf(self.searchQuery) !== -1;
    //     })
    //   }
    // },
    mounted () {
      this.getTasks()
    },
    methods: {
      showCheckboxes () {
        this.checkboxesVisible = true
      },
      hideCheckboxes () {
        this.checkboxesVisible = false
      },
      batchApply () {
        this.selectedTasks.forEach((t) => {
          if (this.batchDelayOn) {
            let seconds = 86400 * this.batchDelay
            this.delayBySeconds(t, seconds)
          }
          this.selectedTasks = []
        })
      },
      setTasks (filter) {
        this.showTasks = filter
        this.getTasks()
      },
      redirectToOrigin (id) {
        let win = window.open(this.$config.apiUrl + 'task/' + id + '/redirect', '_blank')
        win.focus()
      },
      markAsDone (id) {
        this.refreshing = true
        axios({
          method: 'post',
          //headers: {'Authorization': 'Bearer ' + localStorage.getItem('token')},
          url: this.$config.apiUrl + 'task/' + id + '/done'
        }).then((res) => {
          this.getTasks()
        }).catch((err) => {
          if (err.response.status === 401) {
            //this.$store.dispatch('setLoggedOut')
          } else if (err.response.status === 404) {
            //this.$router.push('/404')
          }
        })
      },
      markAsToDo (id) {
        this.refreshing = true
        axios({
          method: 'post',
          //headers: {'Authorization': 'Bearer ' + localStorage.getItem('token')},
          url: this.$config.apiUrl + 'task/' + id + '/todo'
        }).then((res) => {
          this.getTasks()
        }).catch((err) => {
          if (err.response.status === 401) {
            //this.$store.dispatch('setLoggedOut')
          } else if (err.response.status === 404) {
            //this.$router.push('/404')
          }
        })
      },
      delayBySeconds (id, seconds) {
        this.refreshing = true
        axios({
          method: 'post',
          //headers: {'Authorization': 'Bearer ' + localStorage.getItem('token')},
          url: this.$config.apiUrl + 'task/' + id + '/delay/by/' + seconds
        }).then((res) => {
          this.getTasks()
        }).catch((err) => {
          if (err.response.status === 401) {
            //this.$store.dispatch('setLoggedOut')
          } else if (err.response.status === 404) {
            //this.$router.push('/404')
          }
        })
      },
      getTasks () {
        let tasksUrl = this.todoTasksUrl
        if (this.showTasks === 'all') {
          tasksUrl = this.allTasksUrl
        }
        this.refreshing = true
        axios({
          method: 'get',
          //headers: {'Authorization': 'Bearer ' + localStorage.getItem('token')},
          url: this.$config.apiUrl + tasksUrl
        }).then((res) => {
          this.tasks = res.data
          setTimeout(() => {
            this.refreshing = false
          }, 450)
        }).catch((err) => {
          if (err.response.status === 401) {
            //this.$store.dispatch('setLoggedOut')
          } else if (err.response.status === 404) {
            //this.$router.push('/404')
          }
        })
      }
    }
  }
</script>
