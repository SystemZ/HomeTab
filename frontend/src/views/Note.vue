<template>
    <div>
        <h1 class="display-1">#{{id}} {{note.title}}</h1>
        <div class="caption">Tags: {{note.tags}}</div>
        <div class="caption">Created: {{note.createdAt | prettyTimeDate }}</div>
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
        mounted() {
            this.getNote()
            this.simplemde.codemirror.on('change', (instance, changeObj) => {
                if (!this.simplemde.isPreviewActive()) {
                    this.simplemde.togglePreview()
                }
            })
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
                    })
                    .catch(function (err) {
                        if (err.response.status === 401) {
                            console.log("logged out")
                            //vm.$root.$emit("sessionExpired")
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
