<template>
    <div>
        <v-card class="ma-5">
            <v-card-title>
                Notes
                <v-spacer></v-spacer>
                <v-text-field
                        v-model="search"
                        append-icon="mdi-magnify"
                        label="Search"
                        single-line
                        hide-details
                >
                </v-text-field>
            </v-card-title>
            <v-data-table
                    :headers="headers"
                    :items="notes"
                    :search="search"
                    :loading="notesLoading"
                    @click:row="toNote"
            >
                <template v-slot:item.createdAt="{ item }">
                    {{item.createdAt | prettyTimeDate }}
                </template>
                <template v-slot:progress>
                    <v-progress-linear
                            indeterminate
                            :height="2"
                            color="green"
                    ></v-progress-linear>
                </template>
            </v-data-table>
        </v-card>
    </div>
</template>

<script>
    import axios from "axios";

    export default {
        name: 'notes',
        data() {
            return {
                search: '',
                notesLoading: true,
                headers: [
                    {
                        text: 'ID',
                        align: 'left',
                        sortable: true,
                        value: 'id',
                    },
                    {text: 'Title', value: 'title'},
                    {text: 'Short', value: 'short'},
                    {text: 'Tags', value: 'tags'},
                    {text: 'Created at', value: 'createdAt'},
                ],
                notes: [],
            }
        },
        mounted() {
            this.getNotes()
        },
        methods: {
            authConfig() {
                return {headers: {Authorization: "Bearer " + localStorage.getItem(this.lsToken)}}
            },
            toNote(item) {
                this.$router.push({name: 'note', params: {id: item.id}})
            },
            getNotes() {
                let vm = this
                vm.notesLoading = true
                axios.get(vm.apiUrl + "/api/v1/note", vm.authConfig())
                    .then((res) => {
                        vm.notesLoading = false
                        vm.notes = res.data
                    })
                    .catch(function (err) {
                        if (err.response.status === 401) {
                            vm.$root.$emit("sessionExpired")
                        } else {
                            console.log("something wrong")
                        }
                    })
            }
        }
    }
</script>
