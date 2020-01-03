<template>
    <div>
        <h1 class="display-1">#{{id}} {{note.title}}</h1>
        <div class="caption">Tags: {{note.tags}}</div>
        <div class="caption">Created: {{note.createdAt | prettyTimeDate }}</div>
        <v-card class="mt-4">
            <v-card-text>
                {{note.body}}
            </v-card-text>
        </v-card>
    </div>
</template>

<script>
    import axios from "axios";

    export default {
        name: 'note',
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
