<template>
    <v-app>
        <v-navigation-drawer
                v-model="drawer"
                app
        >
            <v-list dense>
                <v-list-item v-if="$store.state.loggedIn" to="/tasks">
                    <v-list-item-action>
                        <v-icon>mdi-clipboard-text</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Tasks</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
                <v-list-item v-if="$store.state.loggedIn" to="/notes">
                    <v-list-item-action>
                        <v-icon>mdi-note</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Notes</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>

                <v-list-item v-if="!$store.state.loggedIn" to="/login">
                    <v-list-item-action>
                        <v-icon>mdi-login</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Login</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
                <v-list-item v-else @click.native="logout" link>
                    <v-list-item-action>
                        <v-icon>mdi-exit-run</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Logout</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>

            </v-list>
        </v-navigation-drawer>

        <v-app-bar
                app
                color="green darken-1"
                dark
        >
            <v-app-bar-nav-icon @click.stop="drawer = !drawer"/>
            <v-toolbar-title>TaskTab</v-toolbar-title>
        </v-app-bar>

        <v-content>
            <v-container
                    fluid
            >
                <router-view/>
            </v-container>
        </v-content>
    </v-app>
</template>

<script lang="ts">
    import Vue from 'vue';

    export default Vue.extend({
        name: 'App',

        data: () => ({
            drawer: null,
        }),
        methods: {
            logout() {
                localStorage.removeItem('authToken')
                this.$store.dispatch('setLoggedOut')
                this.$router.push({name: 'login'})
            }
        }
    });
</script>
