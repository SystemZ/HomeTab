<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12">
        <v-img
          class="text-right pa-2"
        >
        </v-img>

        <v-row no-gutters>
          <v-col cols="12" xs="6" offset-md="4" md="1">
            <v-text-field
              v-model="itemsPerPage"
              label="Per page"
              required
              :color="inputPrimary"
            >
            </v-text-field>
          </v-col>
          <v-col cols="12" xs="6" md="1" offset-md="1">
            <v-text-field
              v-model="afterId"
              label="ID >"
              :color="inputPrimary"
            >
            </v-text-field>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="11" xs="12" md="1" offset-md="5">
            <p v-if="files.length > 0">
              ID {{ files[0].id }} - {{ files[files.length - 1].id }}
            </p>
          </v-col>
        </v-row>
        <v-row no-gutters class="mb-5">
          <v-col cols="11" xs="12">
            <v-pagination
              :disabled="filesLoading"
              :length="pages"
              v-model="page"
              total-visible="5"
              @input="getFiles"
              :color="pagePrimary"
            >
            </v-pagination>
          </v-col>
        </v-row>

        <zoom/>

        <v-item-group>
          <v-container class="pa-0">
            <v-row>
              <v-col
                v-for="(item, itemIndex) in files" :key="itemIndex"
              >
                <v-item>
                  <v-row no-gutters style="height: 220px;">
                    <v-col
                      align="center"
                      align-self="center"
                    >
                      <v-card
                        v-if="isThumbPossible(item.mime)"
                      >
                        <!--<v-icon>mdi-movie-open-outline</v-icon>-->
                        <video
                          v-if="isVideo(item.mime)"
                          @click="$root.$emit('zoomId',item.id)"
                          height="200px" width="200px" loop muted autoplay
                        >
                          <source :src="apiUrl+'/img/full/'+item.sha256" :type="item.mime">
                        </video>
                        <v-img
                          v-else-if="isThumbPossible(item.mime)"
                          :src="urlToThumb(item,200)"
                          height="200"
                          width="200"
                          @click="$root.$emit('zoomId',item.id)"
                        />
                        <span
                          v-if="item.tags"
                          class="d-inline-block text-truncate"
                          style="max-width: 200px;"
                        >
                                                    {{ item.tags }}
                                                </span>
                        <span v-else-if="fileListTagsLoading">...</span>
                        <span v-else>-</span>
                      </v-card>
                      <v-card
                        v-else
                        width="200"
                        height="200"
                      >
                        <v-container fill-height>
                          <v-row class="text-center">
                            <v-col cols="12">
                              {{ item.filename }}
                            </v-col>
                          </v-row>
                        </v-container>

                      </v-card>
                    </v-col>
                  </v-row>

                </v-item>
              </v-col>
            </v-row>
          </v-container>
        </v-item-group>

        <v-row>
          <v-col cols="11" xs="12">
            <v-pagination
              :disabled="filesLoading"
              :length="pages"
              v-model="page"
              total-visible="5"
              @input="getFiles"
              :color="pagePrimary"
            >
            </v-pagination>
          </v-col>
        </v-row>

      </v-col>
    </v-row>


  </v-container>
</template>
<script>
import axios from 'axios'
import Zoom from '../components/Files/Zoom'

export default {
  name: 'files',
  components: {Zoom},
  data () {
    return {
      // loaders
      filesLoading: true,
      fileListTagsLoading: true,
      // view specific
      tagSelected: '',
      search: '',
      files: [],
      // pagination
      page: 1,
      pages: 1,
      itemsPerPage: 56,
      afterId: 0,
      prevItemsPerPage: 0,
      prevPage: 0,
    }
  },
  mounted () {
    this.getFiles()
  },
  created () {
    this.$root.$on('tagSelected', this.selectTag)
  },
  destroyed () {
    this.$root.$off('tagSelected', this.selectTag)
  },
  methods: {
    authConfig () {
      return {headers: {Authorization: 'Bearer ' + localStorage.getItem(this.lsToken)}}
    },
    toFile (item) {
      this.$router.push({name: 'file', params: {id: item.id}})
    },
    isThumbPossible (mime) {
      if (
        mime === 'image/jpeg' ||
        mime === 'image/png' ||
        mime === 'image/gif' ||
        mime === 'video/webm' ||
        mime === 'video/mp4'
      ) {
        return true
      }
    },
    selectTag (tagName) {
      this.tagSelected = tagName
      this.getFiles()
    },
    getFiles () {
      let vm = this
      vm.notesLoading = true
      vm.fileListTagsLoading = true
      let queryType = 'next'
      let lastId = 0

      // search box handling
      if (vm.page === 1) {
        // we don't want to skip records with ID bigger than 1
        lastId = 0
      }

      // // switching items per page handling
      // if (vm.itemsPerPage !== vm.prevItemsPerPage) {
      //   vm.page = 1
      //   lastId = 0
      // }

      if (vm.files.length > 0 && vm.page > vm.prevPage) {
        // next page button
        lastId = vm.files[vm.files.length - 1].id
      } else if (vm.page < vm.prevPage) {
        // prev page button
        lastId = vm.files[0].id
        queryType = 'prev'
      }
      // save page number for next query
      vm.prevPage = vm.page

      //overwrite for fast forwarding through files
      if (this.afterId > 0) {
        lastId = this.afterId
        queryType = 'next'
        this.afterId = 0
      }

      let rawUrl = vm.apiUrl + '/api/v1/files?limit=' + vm.itemsPerPage + '&' + queryType + 'Id=' + lastId
      let data = {}
      if (vm.tagSelected.length > 0 && vm.tagSelected !== 'all') {
        data = {'q': vm.tagSelected}
      }
      axios.post(rawUrl, data, vm.authConfig())
        .then((res) => {
          vm.filesLoading = false
          // show files to user
          vm.files = res.data.files
          // meanwhile...
          // prepare list of files to get tags from
          let fileList = []
          res.data.files.forEach((file) => {
            fileList.push({'fileId': file.id})
          })
          // ask server about tags
          this.getTagsForFiles(fileList)
          // add page to pagination
          if (vm.files.length === vm.itemsPerPage) {
            vm.pages++
          }
          // scroll up so user don't need to
          vm.$vuetify.goTo(0)
        })
        .catch(function (err) {
          if (err.response.status === 401) {
            vm.$root.$emit('sessionExpired')
          } else {
            console.log('something wrong')
          }
        })
    },
    getTagsForFiles (files) {
      axios.post(this.apiUrl + '/api/v1/file/tags', files, this.authConfig())
        .then((tagsRes) => {
          this.files.forEach((file) => {
            // we need to go deeper...
            // check through all files
            tagsRes.data.forEach((tagEntry) => {
              // file get tags attached
              if (file.id === tagEntry.fileId) {
                file.tags = tagEntry.tags
              }
            })
          })

          this.fileListTagsLoading = false
        })
        .catch((err) => {
          if (err.response.status === 401) {
            this.$root.$emit('sessionExpired')
          } else {
            console.log('something wrong')
          }
        })
    }
  }
}
</script>

