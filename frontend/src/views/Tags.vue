<template>
  <v-container fluid>
    <v-card
      class="mx-auto"
      max-width="300"
      tile
    >
      <v-list dense>
        <v-subheader>All tags</v-subheader>
        <v-list-item-group :color="tagPrimary">
          <v-list-item
            v-for="(item, i) in tags"
            :key="i"
          >
            <v-list-item-icon>
              <v-icon>mdi-label</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>
                {{ item.tag }}
                <span class="grey--text">{{ item.counter }}</span></v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list-item-group>
      </v-list>
    </v-card>
  </v-container>
</template>
<script>
import axios from 'axios'

export default {
  name: 'files',
  data () {
    return {
      tagsLoading: true,
      tags: [],
    }
  },
  mounted () {
    this.getTags()
  },
  methods: {
    authConfig () {
      return {headers: {Authorization: 'Bearer ' + localStorage.getItem(this.lsToken)}}
    },
    getTags () {
      this.tagsLoading = true
      let vm = this
      let rawUrl = vm.apiUrl + '/api/v1/tags'
      axios.get(rawUrl, vm.authConfig())
        .then((res) => {
          vm.tags = res.data
          this.tagsLoading = false
        })
        .catch(function (err) {
          if (err.response.status === 401) {
            vm.$root.$emit('sessionExpired')
          } else {
            console.log('something wrong')
          }
        })
    }
  }
}
</script>

