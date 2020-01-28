<template>
    <div>
        <v-snackbar top v-model="noteSaved" :timeout="1500">
            Note saved!
            <v-btn
                    color="green"
                    text
                    @click="noteSaved = false"
            >
                Close
            </v-btn>
        </v-snackbar>
        <v-snackbar top v-model="noteSaveError" :timeout="0">
            Note save error!
            <v-btn
                    color="red"
                    @click="noteSaveError = false"
            >
                Close
            </v-btn>
        </v-snackbar>

        <h1 class="display-1">
            #{{id}}
            <v-text-field v-if="noteEditing" v-model="note.title"></v-text-field>
            <span v-else>{{note.title}}</span>
        </h1>
        <div class="caption">
            Tags:
            <v-text-field v-if="noteEditing" v-model="note.tags"></v-text-field>
            <span v-else>{{note.tags}}</span>
        </div>
        <div class="caption">Created: {{note.createdAt | prettyTimeDate }}</div>
        <v-btn dark color="primary" class="mt-3 mb-5 mr-4" :disabled="noteLoading"
               @click="noteEditing ? noteEditing = false : noteEditing = true">
            <v-icon>mdi-pencil</v-icon>
            Edit
        </v-btn>
        <v-btn dark color="green" class="mt-3 mb-5" :disabled="noteLoading" @click.native="saveNote">
            <v-icon>mdi-content-save-outline</v-icon>
            Save
        </v-btn>
        <vue-simplemde
                v-model="note.body"
                :highlight="true"
                ref="markdownEditor"
        />
    </div>
</template>

<script>
    import axios from "axios";
    import VueSimplemde from 'vue-simplemde'
    import hljs from 'highlight.js';

    window.hljs = hljs;

    export default {
        name: 'note',
        components: {
            VueSimplemde
        },
        data() {
            return {
                noteLoading: true,
                noteSaved: false,
                noteSaveError: false,
                noteEditing: false,
                note: {
                    id: 0,
                    title: "",
                    body: "",
                    tags: "",
                    createdAt: "",
                }
            }
        },
        computed: {
            id() {
                return this.$route.params.id
            },
            simplemde() {
                return this.$refs.markdownEditor.simplemde
            }
        },
        watch: {
            noteEditing(noteEditingNew, noteEditingOld) {
                if (noteEditingNew && this.simplemde.isPreviewActive()) {
                    this.simplemde.togglePreview()
                } else if (!noteEditingNew && this.simplemde.isPreviewActive()) {
                    this.simplemde.togglePreview()
                } else if (!noteEditingNew && !this.simplemde.isPreviewActive()) {
                    this.simplemde.togglePreview()
                }
            }
        },
        mounted() {
            this.getNote()
        },
        methods: {
            authConfig() {
                return {headers: {Authorization: "Bearer " + localStorage.getItem(this.lsToken)}}
            },
            getNote() {
                let vm = this
                vm.notesLoading = true
                axios.get(vm.apiUrl + "/api/v1/note/" + vm.id, vm.authConfig())
                    .then((res) => {
                        vm.noteLoading = false
                        vm.note = res.data
                        setTimeout(() => {
                            vm.simplemde.togglePreview()
                        }, 0)
                    })
                    .catch(function (err) {
                        if (err.response.status === 401) {
                            console.log("logged out")
                            vm.$root.$emit("sessionExpired")
                        } else {
                            console.log("something wrong")
                        }
                    })
            },
            saveNote() {
                let vm = this
                vm.notesLoading = true
                let dataToSend = {
                    'body': vm.note.body,
                    'title': vm.note.title,
                    'tags': vm.note.tags,
                }
                axios.put(vm.apiUrl + "/api/v1/note/" + vm.id, dataToSend, vm.authConfig())
                    .then((res) => {
                        vm.noteLoading = false
                        vm.noteSaveError = false
                        vm.noteSaved = true
                    })
                    .catch(function (err) {
                        vm.noteSaveError = true
                        if (err.response.status === 401) {
                            console.log("logged out")
                            vm.$root.$emit("sessionExpired")
                        } else {
                            console.log("something wrong")
                        }
                    })
            }
        }
    }
</script>
<style>
    @import '~simplemde/dist/simplemde.min.css';
    @import '~github-markdown-css';
    @import '~highlight.js/styles/atom-one-dark.css';
    /* Highlight theme list: https://github.com/isagalaev/highlight.js/tree/master/src/styles */
</style>
