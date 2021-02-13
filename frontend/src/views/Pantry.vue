<template>
    <v-container
            class="fill-height"
            fluid
    >
        <v-row
                align="center"
                justify="center"
        >
            <v-col
                    cols="12"
                    md="4"
                    sm="8"
            >
                <v-card>
                    <v-card-title>
                        <span class="headline">Add to Pantry</span>
                    </v-card-title>
                    <v-card-text>
                        <v-container>
                            <v-row>
                                <v-col cols="12">
                                    <v-text-field
                                            clearable
                                            label="Name"
                                            required
                                            v-model="item.name"
                                            :color="inputPrimary"
                                    ></v-text-field>
                                </v-col>
                                <v-col cols="12">
                                    <v-autocomplete
                                            :items="['Basilur', 'Tesco', 'Prymat']"
                                            label="Brand"
                                            required
                                            v-model="item.brand"
                                            :color="inputPrimary"
                                    >
                                    </v-autocomplete>
                                </v-col>
                                <v-col cols="12">
                                    <v-text-field
                                            label="Bar code"
                                            required
                                            v-model="item.code"
                                            :color="inputPrimary"
                                    ></v-text-field>
                                </v-col>
                                <v-col cols="12" sm="6">
                                    <v-select
                                            :items="['Fridge', 'Refrigerator', 'Under oven', 'Somewhere else']"
                                            label="Localization"
                                            required
                                            v-model="item.place"
                                            :color="inputPrimary"
                                    ></v-select>
                                </v-col>
                                <v-col cols="12" sm="6">
                                    <v-autocomplete
                                            :items="['Drinks', 'Spices', 'Cereals', 'Veggies & Fruit']"
                                            label="Type of food"
                                            v-model="item.type"
                                            :color="inputPrimary"
                                    ></v-autocomplete>
                                </v-col>
                            </v-row>
                        </v-container>
                    </v-card-text>
                    <v-card-actions>
                        <v-spacer></v-spacer>
                        <v-btn :color="btnPrimary" @click="saveItem" text>Save</v-btn>
                    </v-card-actions>
                </v-card>
            </v-col>
        </v-row>

        <v-row>
            <v-col>
                <v-data-table
                        :headers="headers"
                        :items="items"
                        :items-per-page="5"
                        class="elevation-1"
                >
                </v-data-table>
            </v-col>
        </v-row>
    </v-container>

</template>

<script>
    import axios from "axios";

    export default {
        name: 'pantry',
        data() {
            return {
                item: {
                    id: 2,
                    name: "",
                    brand: "",
                    code: 0,
                    place: "",
                    type: "",
                },
                items: [
                    {
                        id: 1,
                        name: "Green tea",
                        brand: "Basilur",
                        code: 1234567890,
                        place: "Somewhere else",
                        type: "Drinks",
                    },
                ],
                headers: [
                    {
                        text: 'ID',
                        align: 'left',
                        sortable: true,
                        value: 'id',
                    },
                    {text: 'Name', sortable: true, value: 'name'},
                    {text: 'Brand', sortable: false, value: 'brand'},
                    {text: 'Place', sortable: false, value: 'place'},
                    {text: 'Type', sortable: false, value: 'type'},
                ],
            }
        },
        mounted() {
        },
        methods: {
            saveItem() {
                this.items.push({
                    "id": this.item.id++,
                    "name": this.item.name,
                    "brand": this.item.brand,
                    "code": this.item.code,
                    "place": this.item.place,
                    "type": this.item.type,
                });
                this.item.name = ""
                this.item.brand = ""
                this.item.code = 0
                this.item.place = ""
                this.item.type = ""
            },
            authConfig() {
                return {headers: {Authorization: "Bearer " + localStorage.getItem(this.lsToken)}}
            }
        },
    }
</script>
