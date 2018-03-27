<template>
  <v-container fluid>
    <v-slide-y-transition mode="out-in">
      <v-layout row wrap>
        <v-flex xs12 lg2>
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
        </v-flex>

        <v-flex xs12 lg10>
          <v-card>
            <v-list>
              <v-list-group
                v-model="task.active"
                v-for="task in tasks"
                :key="task.id"
                prepend-icon="fa fa-circle-thin"
                no-action
                v-if="show[task.type] == true"
              >
                <!-- main task row start -->
                <v-list-tile slot="activator">
                  <v-list-tile-content>
                    <v-list-tile-title>
                      <v-icon v-if="task.done">fa fa-check</v-icon>
                      <v-icon v-if="task.type == 'github'">fa fa-github</v-icon>
                      <v-icon v-if="task.type == 'gitlab'">fa fa-gitlab</v-icon>
                      <v-icon v-if="task.type == 'gmail'">fa fa-at</v-icon>
                      {{ task.title }}
                    </v-list-tile-title>
                  </v-list-tile-content>
                </v-list-tile>
                <!-- main task row end -->

                <!-- action for task row start -->
                <v-list-tile>
                  <v-list-tile-content>
                    <v-list-tile-title>zzz</v-list-tile-title>
                  </v-list-tile-content>
                  <v-list-tile-action @click="getTasks">
                    <v-btn>
                      <v-icon>help</v-icon>
                      Refresh
                    </v-btn>
                  </v-list-tile-action>
                </v-list-tile>
                <!-- action for task row end -->

              </v-list-group>
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
        refreshing: true,
        tasks: [],
        showTasks: 'todo',
        allTasksUrl: 'tasks/all/1',
        todoTasksUrl: 'tasks/todo/1',
        show: {
          'github': true,
          'gitlab': true,
          'gmail': true
        }
      }
    },
    mounted () {
      this.getTasks()
    },
    methods: {
      setTasks (filter) {
        this.showTasks = filter
        this.getTasks()
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
