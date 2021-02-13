<template>
  <v-container
    class="fill-height"
    fluid
  >
    <v-row
      align="center"
      justify="center"
    >
      <v-col
        cols="12"
        sm="8"
        md="4"
      >
        <v-alert v-if="loginError" type="warning" class="mt-5">
          {{ loginErrorList[loginErrorReason] }}
        </v-alert>
        <v-card class="elevation-12 mt-5 mb-5" color="gray darken-3">
          <v-toolbar
            :color="toolbarPrimary"
            :dark="componentDark"
            flat
          >
            <v-toolbar-title>
              <span>Log in</span>
            </v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <v-form @keyup.native.enter="login">
              <v-text-field
                label="Username"
                name="login"
                prepend-icon="mdi-account"
                type="text"
                v-model="username"
                :color="inputPrimary"
              />

              <v-text-field
                id="password"
                label="Password"
                name="password"
                prepend-icon="mdi-lock"
                type="password"
                v-model="password"
                :color="inputPrimary"
              />
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-spacer/>
            <v-btn :dark="componentDark" :color="btnPrimary" @click.native="login">
              <span>Log in</span>
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import axios from 'axios'

export default {
  name: 'login',
  data () {
    return {
      username: '',
      password: '',
      loginError: false,
      loginErrorReason: 0,
      loginErrorList: {
        401: 'Wrong password',
        403: 'Wrong username',
        500: 'Problem with server, please try again later',
        0: 'Unexpected problem with server'
      },
    }
  },
  methods: {
    login () {
      let vm = this
      axios.post(vm.apiUrl + '/api/v1/login', {
        'username': vm.username,
        'password': vm.password
      }).then((res) => {
        vm.loginError = false
        vm.username = ''
        vm.password = ''
        localStorage.setItem('authToken', res.data.token)
        this.$store.dispatch('setLoggedIn')
        vm.$router.push({name: 'notes'})
      }).catch(function (err) {
        vm.loginError = true
        if (err.response.status === 401) {
          vm.loginErrorReason = err.response.status
          //vm.$root.$emit('sessionExpired');
        } else if (err.response.status === 403) {
          vm.loginErrorReason = err.response.status
        } else if (err.response.status === 500) {
          vm.loginErrorReason = err.response.status
        } else {
          vm.loginErrorReason = 0
        }
      })
    },
  },
}
</script>
