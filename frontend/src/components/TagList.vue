<template>
  <v-card
    class="mx-auto"
    max-width="300"
    tile
  >
    <v-list dense>
      <v-subheader>All tags</v-subheader>
      <v-list-item-group :color="tagPrimary">
        <v-list-item
          @click="selectTag('all')"
        >
          <v-list-item-icon>
            <v-icon>mdi-earth</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>
              All
              <span class="grey--text">?</span>
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item
          @click="selectTag('none')"
        >
          <v-list-item-icon>
            <v-icon>mdi-label-off</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>
              Without
              <span class="grey--text">?</span>
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item
          v-for="(item, i) in tagsRaw" :key="i"
          @click="selectTag(item.tag)"
        >
          <v-list-item-icon>
            <v-icon>mdi-label-outline</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>
              {{ item.tag }}
              <span class="grey--text">{{ item.counter }}</span>
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list-item-group>
    </v-list>
  </v-card>
</template>

<script>
import axios from 'axios'

export default {
  name: 'TagList',

  data: () => ({
    tagsLoading: true,
    tagsRaw: [],
    tags: [],
  }),
  mounted () {
    this.getTags()
  },
  methods: {
    authConfig () {
      return {headers: {Authorization: 'Bearer ' + localStorage.getItem(this.lsToken)}}
    },
    selectTag (tagName) {
      this.tagSelected = tagName
      this.$root.$emit('tagSelected', tagName)
    },
    getTags () {
      this.tagsLoading = true
      let rawUrl = this.apiUrl + '/api/v1/tags'
      axios.get(rawUrl, this.authConfig())
        .then((res) => {
          // for tags widget
          this.tagsRaw = res.data
          // for tags under zoomed file
          res.data.forEach((entry) => {
            this.tags.push({'text': entry.tag + ' (' + entry.counter + ')', 'value': entry.tag})
          })
          this.tagsLoading = false
        })
        .catch((err) => {
          console.log(err)
          if (err.response.status === 401) {
            this.$root.$emit('sessionExpired')
          } else {
            console.log('something wrong')
          }
        })
    },
  }
}
</script>
