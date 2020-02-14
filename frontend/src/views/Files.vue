<template>
    <v-container fluid>

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


        <v-card v-if="bigPic" class="text-center" elevation="0">

            <v-btn :href="apiUrl+'/img/full/'+bigPicInfo.sha256" target="_blank" class="mb-5">Original</v-btn>
            <br>

            <img
                    @click="bigPic = false"
                    :src="apiUrl+'/img/thumbs/700/700/'+bigPicInfo.sha256"/>
            <br>

            <v-chip
                    v-for="(tag, tagIndex) in tagz(bigPicInfo.tags)"
                    :key="tagIndex"
                    class="ma-2"
                    color="indigo"
                    label
                    text-color="white"
                    large
                    close
                    @click:close="deleteTag(bigPicInfo,tag)"
            >
                <v-icon left>mdi-label</v-icon>
                {{tag}}
            </v-chip>
        </v-card>
        <v-item-group v-else>
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
                                        @click="bigPic = true , bigPicInfo = item"
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
        files: [],
        page: 1,
        pages: 1,
        itemsPerPage: 56,
        prevItemsPerPage: 0,
        prevPage: 0,
        bigPic: false,
        bigPicInfo: {},
        // change me
        active: false,
      }
    },
    mounted () {
      this.getFiles()
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
        if (!str.includes(',')) {
          return []
        }
        return str.split(',')
      },
      deleteTag (bigPicInfo, tag) {
        console.log(bigPicInfo, tag)
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
        axios.post(rawUrl, {}, vm.authConfig())
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
      }
    }
  }
</script>

