<template>
    <v-container fluid>
        <v-row>
            <v-col cols="2">
                <v-card
                        class="mx-auto"
                        max-width="300"
                        tile
                >
                    <v-list dense>
                        <v-subheader>All tags</v-subheader>
                        <v-list-item-group color="primary">
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
                                        {{item.tag}}
                                        <span class="grey--text">{{item.counter}}</span>
                                    </v-list-item-title>
                                </v-list-item-content>
                            </v-list-item>
                        </v-list-item-group>
                    </v-list>
                </v-card>
            </v-col>
            <v-col cols="10">
                <v-img
                        class="text-right pa-2"
                >
                </v-img>

                <v-row>
                    <v-col cols="11" xs="12">
                        <v-pagination
                                :disabled="filesLoading"
                                :length="pages"
                                v-model="page"
                                total-visible="5"
                                @input="getFiles"
                        >
                        </v-pagination>
                    </v-col>
                    <v-col cols="1" xs="12">
                        <v-text-field
                                v-model="itemsPerPage"
                                label="Per page"
                                required
                        >
                        </v-text-field>

                    </v-col>
                </v-row>

                <v-dialog v-model="bigPic" dark max-width="90%">
                    <v-card>
                        <v-card class="text-center" elevation="0">
                            <v-btn :href="apiUrl+'/img/full/'+bigPicInfo.sha256" target="_blank" class="mb-5">
                                Original
                            </v-btn>
                            <br>
                            {{bigPicInfo.filePath}}
                            <br>
                            <img
                                    @click="bigPic = false"
                                    :src="apiUrl+'/img/thumbs/700/700/'+bigPicInfo.sha256"/>
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
                                                    :disabled="tagsLoading"
                                                    v-model="bigPicInfo.tagz"
                                                    @change="addTag"
                                                    :items="tags"
                                                    chips
                                                    label="File tags"
                                                    multiple
                                                    prepend-icon="mdi-label"
                                                    :return-object="false"
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
                                                            @click:close="deleteTag(bigPicInfo,item)"
                                                    >
                                                        <strong>{{ item }}</strong>&nbsp;
                                                        <!--<span>(interest)</span>-->
                                                    </v-chip>
                                                </template>
                                            </v-combobox>
                                        </v-col>
                                    </v-row>
                                    <v-row>
                                        <v-flex cols="1" v-for="similar in similarFiles" :key="similar.sha256">
                                            <a target="_blank" :href="apiUrl+'/img/full/'+similar.sha256">
                                                {{similar.distance}}
                                            </a>
                                            <v-img @click="zoom(similar,true)"
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

                <v-item-group>
                    <v-container class="pa-0">
                        <v-row>
                            <v-col
                                    v-for="(item, itemIndex) in files" :key="itemIndex"
                            >
                                <v-item>
                                    <v-card
                                            v-if="isThumbPossible(item.mime)"
                                            width="200"
                                            height="200"
                                    >
                                        <v-img

                                                :src="apiUrl+'/img/thumbs/200/200/'+item.sha256"
                                                height="200"
                                                width="200"
                                                class="text-right pa-2"
                                                @click="zoom(item,false)"
                                        >
                                            <v-btn
                                                    icon
                                                    dark
                                            >
                                                <!--<v-icon>{{ active ? 'mdi-heart' : 'mdi-heart-outline' }}</v-icon>-->
                                            </v-btn>
                                        </v-img>

                                    </v-card>
                                    <v-card
                                            v-else
                                            width="200"
                                            height="200"
                                    >
                                        <v-container fill-height>
                                            <v-row class="text-center">
                                                <v-col cols="12">
                                                    {{item.filename}}
                                                </v-col>
                                            </v-row>
                                        </v-container>

                                    </v-card>

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
                        >
                        </v-pagination>
                    </v-col>
                    <v-col cols="1" xs="12">
                        <v-text-field
                                v-model="itemsPerPage"
                                label="Per page"
                                required
                        >
                        </v-text-field>

                    </v-col>
                </v-row>

            </v-col>
        </v-row>


    </v-container>
</template>
<script>
  import axios from 'axios'

  export default {
    name: 'files',
    data () {
      return {
        search: '',
        filesLoading: true,
        similarLoading: true,
        files: [],
        page: 1,
        pages: 1,
        itemsPerPage: 56,
        prevItemsPerPage: 0,
        prevPage: 0,
        bigPic: false,
        bigPicInfo: {},
        tagsLoading: false,
        tags: [],
        tagsRaw: [],
        tagSelected: '',
        similarFiles: [],
        // change me
        active: false,
        // tag typed by user in zoomed image
        tagTyped: ''
      }
    },
    mounted () {
      this.getFiles()
      this.getTags()
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
      zoom (item, similar) {
        if (similar && this.similarLoading) {
          // TODO notify user about loading
          // this prevents skipping tag loading for similar images
          return
        }
        this.getSimilar(item.sha256)
        this.bigPic = true
        this.bigPicInfo = item
        this.bigPicInfo.tagz = this.tagz(item.tags)
        this.bigPicInfo.tagzServer = this.tagz(item.tags)
      },
      selectTag (tagName) {
        this.tagSelected = tagName
        this.getFiles()
      },
      addTag (currentTags) {
        this.tagsLoading = true
        let diff = currentTags
          .filter(x => !this.bigPicInfo.tagzServer.includes(x))
          .concat(this.bigPicInfo.tagzServer.filter(x => !currentTags.includes(x)))

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
        let rawUrl = vm.apiUrl + '/api/v1/file/' + this.bigPicInfo.sha256 + '/tag/add'
        axios.post(rawUrl, {'tag': tagToAdd}, vm.authConfig())
          .then(() => {
            vm.bigPicInfo.tagzServer = currentTags
            vm.getTags()
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
      deleteTag (bigPicInfo, tag) {
        console.log('Deleting tag ' + tag)
        this.tagsLoading = true
        let vm = this
        let rawUrl = vm.apiUrl + '/api/v1/file/' + bigPicInfo.sha256 + '/tag/delete'
        axios.post(rawUrl, {'tag': tag}, vm.authConfig())
          .then(() => {
            for (let i = vm.bigPicInfo.tagz.length - 1; i >= 0; i--) {
              if (vm.bigPicInfo.tagz[i] === tag) {
                vm.bigPicInfo.tagz.splice(i, 1)
              }
            }
            vm.getTags()
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
      getFiles () {
        let vm = this
        vm.notesLoading = true
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

        let rawUrl = vm.apiUrl + '/api/v1/files?limit=' + vm.itemsPerPage + '&' + queryType + 'Id=' + lastId
        let data = {}
        if (vm.tagSelected.length > 0 && vm.tagSelected !== 'all') {
          data = {'q': vm.tagSelected}
        }
        axios.post(rawUrl, data, vm.authConfig())
          .then((res) => {
            vm.filesLoading = false
            vm.files = res.data.files
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
      getTags () {
        this.tagsLoading = true
        let vm = this
        let rawUrl = vm.apiUrl + '/api/v1/tags'
        axios.get(rawUrl, vm.authConfig())
          .then((res) => {
            // for tags widget
            vm.tagsRaw = res.data
            // for tags under zoomed file
            res.data.forEach((entry) => {
              vm.tags.push({'text': entry.tag + ' (' + entry.counter + ')', 'value': entry.tag})
            })
            this.tagsLoading = false
          })
          .catch(function (err) {
            if (err.response.status === 401) {
              vm.$root.$emit('sessionExpired')
            } else {
              console.log('something wrong')
            }
          })
      },
      getSimilar (sha256) {
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
            this.getTagsForFiles(fileList)
          })
          .catch((err) => {
            if (err.response.status === 401) {
              this.$root.$emit('sessionExpired')
            } else {
              console.log('something wrong')
            }
          })
      },
      getTagsForFiles (files) {
        axios.post(this.apiUrl + '/api/v1/file/tags', files, this.authConfig())
          .then((tagsRes) => {
            this.similarFiles.forEach((file) => {
              // we need to go deeper...
              // check through all files
              tagsRes.data.forEach((tagEntry) => {
                // file get tags attached
                if (file.id === tagEntry.fileId) {
                  file.tags = tagEntry.tags
                }
              })
            })
            this.similarLoading = false
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

