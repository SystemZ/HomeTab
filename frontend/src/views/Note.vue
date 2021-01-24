<template>
    <v-container fluid>
        <v-snackbar top v-model="noteSaved" :timeout="1500">
            Note saved!
            <v-btn
                    :color="btnAccent"
                    text
                    @click="noteSaved = false"
            >
                Close
            </v-btn>
        </v-snackbar>
        <v-snackbar top v-model="noteSaveError" :timeout="0">
            Note save error!
            <v-btn
                    :color="btnSecondary"
                    @click="noteSaveError = false"
            >
                Close
            </v-btn>
        </v-snackbar>
        <v-dialog v-model="deleteNoteDialog" max-width="700" :fullscreen="$vuetify.breakpoint.xsOnly">
            <v-card>
                <v-card-title>
                    Delete note
                </v-card-title>
                <v-card-subtitle>
                    Are you sure that you don't need this note?
                </v-card-subtitle>
                <v-card-text>
                </v-card-text>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn
                            :color="btnPrimary"
                            text
                            @click="deleteNoteDialog = false"
                    >
                        Cancel
                    </v-btn>
                    <v-btn
                            :dark="btnDark"
                            :color="btnSecondary"
                            :disabled="noteLoading"
                            @click="deleteNote"
                    >
                        <v-icon>mdi-delete</v-icon>
                        Delete
                    </v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>

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
        <v-btn :dark="btnDark" :color="btnPrimary" class="mt-3 mb-5 mr-4" :disabled="noteLoading"
               @click="noteEditing ? noteEditing = false : noteEditing = true">
            <v-icon>mdi-pencil</v-icon>
            Edit
        </v-btn>
        <v-btn :dark="btnDark" :color="btnPrimary" class="mt-3 mb-5 mr-4" :disabled="noteLoading"
               @click.native="saveNote">
            <v-icon>mdi-content-save-outline</v-icon>
            Save
        </v-btn>
        <v-btn :dark="btnDark" color="red darken-2" class="mt-3 mb-5" :disabled="noteLoading"
               @click.native="deleteNoteDialog = true">
            <v-icon>mdi-delete</v-icon>
            Delete
        </v-btn>
        <vue-simplemde
                v-model="note.body"
                :highlight="true"
                ref="markdownEditor"
        />
    </v-container>
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
                deleteNoteDialog: false,
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
            },
            deleteNote() {
                this.noteLoading = true
                axios.delete(this.apiUrl + "/api/v1/note/" + this.id, this.authConfig())
                    .then((res) => {
                        this.$router.push({name: 'notes'})
                    })
                    .catch((err) => {
                        this.noteLoading = false
                        if (err.response.status === 401) {
                            console.log("logged out")
                            this.$root.$emit("sessionExpired")
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
