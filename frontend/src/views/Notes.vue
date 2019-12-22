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
            >
                <template v-slot:item.createdAt="{ item }">
                    {{item.createdAt | prettyTimeDate}}
                </template>
            </v-data-table>
        </v-card>
    </div>
</template>

<script>
    export default {
        name: 'notes',
        data() {
            return {
                search: '',
                headers: [
                    {
                        text: 'ID',
                        align: 'left',
                        sortable: false,
                        value: 'id',
                    },
                    {text: 'Title', value: 'title'},
                    {text: 'Short', value: 'short'},
                    {text: 'Tags', value: 'tags'},
                    {text: 'Created at', value: 'createdAt'},
                ],
                notes: [
                    {
                        id: 1,
                        title: "Smaller plan",
                        short: "No",
                        tags: ["lvlup", "ticket"],
                        createdAt: "2018-10-24T17:53:56+01:00",
                    },
                ],
            }
        },
    }
</script>
