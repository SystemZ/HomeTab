<template>
    <v-container fluid>

        <v-container fluid>
            <v-row>
                <v-col cols="12" lg="2" md="4" sm="6" xs="12">
                    <v-text-field
                            label="Title"
                            v-model="newNoteTitle"
                            @keydown.enter.native="addNote"
                    ></v-text-field>
                </v-col>
                <v-col cols="12" md="2">
                    <v-btn class="mt-3" :dark="btnDark" :color="btnPrimary" @click.native="addNote">
                        Add new
                    </v-btn>
                </v-col>
            </v-row>
        </v-container>

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
    </v-container>
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
                newNoteTitle: "",
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
                this.notesLoading = true
                axios.get(this.apiUrl + "/api/v1/note", this.authConfig())
                    .then((res) => {
                        this.notesLoading = false
                        this.notes = res.data
                    })
                    .catch((err) => {
                        if (err.response.status === 401) {
                            this.$root.$emit("sessionExpired")
                        } else {
                            console.log("something wrong")
                        }
                    })
            },
            addNote() {
                axios.post(this.apiUrl + "/api/v1/note", {"title": this.newNoteTitle}, this.authConfig())
                    .then((res) => {
                        this.$router.push("/note/" + res.data.id)
                    })
                    .catch((err) => {
                        if (err.response.status === 401) {
                            this.$root.$emit("sessionExpired")
                        } else {
                            console.log("something wrong")
                        }
                    })
            }
        }
    }
</script>
