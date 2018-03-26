<template>
  <v-container fluid>
    <v-slide-y-transition mode="out-in">
      <v-layout column align-center>

        <ul>
          <li v-for="task in tasks">
            {{ task.title }}
          </li>
        </ul>

        <img src="@/assets/logo.png" alt="Vuetify.js" class="mb-5">
        <blockquote>
          &#8220;First, solve the problem. Then, write the code.&#8221;
          <footer>
            <small>
              <em>&mdash;John Johnson</em>
            </small>
          </footer>
        </blockquote>
      </v-layout>
    </v-slide-y-transition>
  </v-container>
</template>

<script>
  import axios from 'axios'

  export default {
    data () {
      return {
        tasks: []
      }
    },
    mounted () {
      this.getTasks()
    },
    methods: {
      getTasks () {
        axios({
          method: 'get',
          //headers: {'Authorization': 'Bearer ' + localStorage.getItem('token')},
          url: this.$config.apiUrl + 'task/1'
        }).then((res) => {
          this.tasks = res.data
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
