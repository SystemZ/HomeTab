<template>
    <v-dialog v-model="dialog" dark max-width="90%">
        <v-card>
            <v-card class="text-center" elevation="0">
                <div v-if="mainFileLoading" class="pa-5 ma-5">
                    <v-skeleton-loader
                            class="mx-auto"
                            width="600"
                            height="300"
                            type="image"
                    />
                </div>
                <div v-else>
                     <span class="mr-3">
                                 ID <kbd>{{file.id}}</kbd>
                            </span>
                    <v-btn :href="apiUrl+'/img/full/'+file.sha256" target="_blank" class="mb-2 mt-2 mr-5" color="grey darken-3">
                        Original
                    </v-btn>
                    <v-btn @click="(showFilePath) ? showFilePath = false : showFilePath = true" color="grey darken-3">
                        Path
                    </v-btn>
                    <div v-if="showFilePath">
                        <kbd class="mt-1 mb-2">{{file.filePath}}</kbd>
                    </div>
                    <div v-if="isVideo(file.mime) && dialog">
                        <video height="600" controls loop autoplay>
                            <source :src="apiUrl+'/img/full/'+file.sha256" :type="file.mime">
                            Your browser does not support the video tag.
                        </video>
                    </div>
                    <div v-else>
                        <img @click="dialog = false" :src="urlToThumb(file,700)"/>
                    </div>
                </div>
                <br>
                <v-form>
                    <v-container>
                        <v-row>
                            <v-col
                                    cols="12"
                                    lg="6"
                                    offset-lg="3"
                            >
                                <v-combobox
                                        :search-input.sync="tagTyped"
                                        v-model="file.tagz"
                                        @change="addTag"
                                        :items="tags"
                                        chips
                                        label="File tags"
                                        multiple
                                        prepend-icon="mdi-label"
                                        :return-object="false"
                                        :loading="allTagsLoading || mainFileTagsLoading"
                                        :disabled="allTagsLoading || mainFileTagsLoading"
                                        :disable-lookup="allTagsLoading || mainFileTagsLoading"
                                >

                                    <!-- clearable -->
                                    <template v-slot:selection="{ attrs, item, select, selected }">
                                        <v-chip
                                                v-bind="attrs"
                                                color="indigo"
                                                dark
                                                class="ma-2"
                                                large
                                                label
                                                :input-value="selected"
                                                close
                                                @click="select"
                                                @click:close="deleteTag(file,item)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                            <!--<span>(interest)</span>-->
                                        </v-chip>
                                    </template>
                                </v-combobox>
                            </v-col>
                        </v-row>
                        <v-row v-if="similarLoading">
                            <v-progress-linear
                                    indeterminate
                                    color="primary"
                            />
                        </v-row>
                        <v-row v-else>
                            <v-flex cols="1" v-for="similar in similarFiles" :key="similar.sha256">
                                <a target="_blank" :href="apiUrl+'/img/full/'+similar.sha256">
                                    {{similar.distance}}
                                </a>
                                <v-img @click="zoom(similar.id)"
                                       width="150"
                                       height="150"
                                       :src="apiUrl+'/img/thumbs/150/150/'+similar.sha256">
                                </v-img>
                                <span v-if="similar.tags"
                                      class="d-inline-block text-truncate"
                                      style="max-width: 150px;">{{similar.tags}}</span>
                                <span v-else-if="similarLoading">...</span>
                                <span v-else>-</span>
                            </v-flex>
                        </v-row>
                    </v-container>
                </v-form>

            </v-card>

        </v-card>

    </v-dialog>
</template>

<script>
  import axios from 'axios'

  export default {
    name: 'Zoom',

    data: () => ({
      // main dialog visibility
      dialog: false,
      // loaders
      mainFileLoading: true,
      mainFileTagsLoading: true,
      allTagsLoading: true,
      similarLoading: true,
      // view specific
      tags: [],
      file: {},
      showFilePath: false,
      similarFiles: [],
      // tag typed by user
      tagTyped: '',
    }),
    mounted () {
    },
    created () {
      this.$root.$on('zoomId', this.zoom)
    },
    destroyed () {
      this.$root.$off('zoomId', this.zoom)
    },
    methods: {
      authConfig () {
        return {headers: {Authorization: 'Bearer ' + localStorage.getItem(this.lsToken)}}
      },
      zoom (fileId) {
        // reset all loaders
        this.mainFileLoading = true
        this.allTagsLoading = true
        this.mainFileTagsLoading = true
        this.similarLoading = true
        // show content to user
        this.dialog = true
        // async data gathering
        this.getAllTags()
        this.getFile(fileId)
      },
      getTagsForFiles (files, similar) {
        if (!similar) this.mainFileTagsLoading = true
        axios.post(this.apiUrl + '/api/v1/file/tags', files, this.authConfig())
          .then((tagsRes) => {
            if (similar) {
              // for thumbs under main file
              this.similarFiles.forEach((file) => {
                // we need to go deeper...
                // check through all files
                tagsRes.data.forEach((tagEntry) => {
                  // file get tags attached
                  if (file.id === tagEntry.fileId) {
                    file.tags = tagEntry.tags
                  }
                })
                this.similarLoading = false
              })
            } else {
              // for main file
              tagsRes.data.forEach((tagEntry) => {
                if (this.file.id === tagEntry.fileId) {
                  this.file.tags = tagEntry.tags
                }
              })
              // apply tags
              this.file.tagz = this.tagz(this.file.tags)
              this.file.tagzServer = this.tagz(this.file.tags)
              this.mainFileTagsLoading = false
            }
          })
          .catch((err) => {
            if (err.response.status === 401) {
              this.$root.$emit('sessionExpired')
            } else {
              console.log('something wrong')
            }
          })
      },
      tagz (str) {
        // no tags
        if (!str.includes(',') && str.length < 1) {
          return []
        }
        // two or more
        if (str.includes(',') && str.length > 0) {
          return str.split(',')
        }
        // one tag
        return [str]
      },
      getFile (id) {
        this.mainFileLoading = true
        axios.get(this.apiUrl + '/api/v1/file/' + id, this.authConfig())
          .then((res) => {
            this.file = res.data
            this.mainFileLoading = false
            // ask server about tags
            this.getTagsForFiles([{'fileId': res.data.id}], false)
            // get similar thumbs under main file
            this.getSimilar(this.file.sha256)
          })
          .catch((err) => {
            if (err.response.status === 401) {
              this.$root.$emit('sessionExpired')
            } else {
              console.log('something wrong')
            }
          })
      },
      getAllTags () {
        this.allTagsLoading = true
        let rawUrl = this.apiUrl + '/api/v1/tags'
        axios.get(rawUrl, this.authConfig())
          .then((res) => {
            // for tags under zoomed file
            res.data.forEach((entry) => {
              this.tags.push({'text': entry.tag + ' (' + entry.counter + ')', 'value': entry.tag})
            })
            this.allTagsLoading = false
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
      addTag (currentTags) {
        //this.tagsLoading = true
        let diff = currentTags
          .filter(x => !this.file.tagzServer.includes(x))
          .concat(this.file.tagzServer.filter(x => !currentTags.includes(x)))

        let tagToAdd = ''

        if (diff.length < 1) {
          console.log('detected to try an empty tag, cancelling...')
          return
        }

        // this requires combobox :return-object="false"
        // it makes tag selecting work properly without objects visible for user
        // https://github.com/vuetifyjs/vuetify/issues/5358#issuecomment-431312918
        tagToAdd = diff[0]
        console.log('Adding tag ' + tagToAdd)

        // this removes manually typed tag if user clicks on tag suggestion
        this.tagTyped = ''

        let vm = this
        let rawUrl = vm.apiUrl + '/api/v1/file/' + this.file.sha256 + '/tag/add'
        axios.post(rawUrl, {'tag': tagToAdd}, vm.authConfig())
          .then(() => {
            vm.file.tagzServer = currentTags
            // get global tag list
            vm.getAllTags()
            // ask server about tags for this one file on main list
            this.getTagsForFiles([{'fileId': this.file.id}], false)
          })
          .catch(function (err) {
            if (err.response.status === 401) {
              vm.$root.$emit('sessionExpired')
            } else {
              console.log(err)
              console.log('something wrong')
            }
          })

      },
      deleteTag (file, tag) {
        console.log('Deleting tag ' + tag)
        // this.tagsLoading = true
        let vm = this
        let rawUrl = vm.apiUrl + '/api/v1/file/' + file.sha256 + '/tag/delete'
        axios.post(rawUrl, {'tag': tag}, vm.authConfig())
          .then(() => {
            for (let i = vm.file.tagz.length - 1; i >= 0; i--) {
              if (vm.file.tagz[i] === tag) {
                vm.file.tagz.splice(i, 1)
              }
            }
            // get global tag list
            vm.getAllTags()
            // ask server about tags for this one file on main list
            this.getTagsForFiles([{'fileId': this.file.id}], false)
          })
          .catch(function (err) {
            console.log(err)
            if (err.response.status === 401) {
              vm.$root.$emit('sessionExpired')
            } else {
              console.log('something wrong')
            }
          })
      },
      getSimilar (sha256) {
        this.similarLoading = true
        axios.get(this.apiUrl + '/api/v1/file/' + sha256 + '/similar', this.authConfig())
          .then((res) => {
            // prevent strange user clicks
            this.similarLoading = true
            // show user suggested files
            this.similarFiles = res.data
            // meanwhile...
            // prepare list of files to get tags from
            let fileList = []
            res.data.forEach((file) => {
              fileList.push({'fileId': file.id})
            })
            // ask server about tags
            this.getTagsForFiles(fileList, true)
            this.similarLoading = false
          })
          .catch((err) => {
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
