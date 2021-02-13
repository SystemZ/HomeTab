<template>
  <v-container>
    <v-card>
      <v-container>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field label="Dir to scan" v-model="dirToScan" :color="inputPrimary">
            </v-text-field>
            <v-btn
              :disabled="scanning"
              @click="doScan"
              :color="btnPrimary"
              :dark="componentDark"
            >
              Scan
            </v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-card>
    <v-snackbar v-model="showSnackbar">
      Scan queued, please wait
    </v-snackbar>
  </v-container>
</template>
<script>
import axios from 'axios'

export default {
  name: 'scan',
  data () {
    return {
      scanning: false,
      showSnackbar: false,
      dirToScan: '/mnt/smb/data/Z/filez',
    }
  },
  mounted () {
  },
  methods: {
    authConfig () {
      return {headers: {Authorization: 'Bearer ' + localStorage.getItem(this.lsToken)}}
    },
    doScan () {
      this.scanning = true
      let data = [{'path': this.dirToScan}]
      let rawUrl = this.apiUrl + '/api/v1/scan'
      axios.post(rawUrl, data, this.authConfig())
        .then(() => {
          this.showSnackbar = true
        })
        .catch((err) => {
          if (err.response.status === 401) {
            this.$root.$emit('sessionExpired')
          } else {
            console.log(err)
            console.log('something wrong')
          }
        })
    },
  }
}
</script>

