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
                <v-list-item v-if="$store.state.loggedIn" to="/counters">
                    <v-list-item-action>
                        <v-icon>mdi-timer</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Counters</v-list-item-title>
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
                <v-list-item v-if="$store.state.loggedIn" to="/events">
                    <v-list-item-action>
                        <v-icon>mdi-calendar-clock</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Events</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
                <v-list-item v-if="$store.state.loggedIn" to="/devices">
                    <v-list-item-action>
                        <v-icon>mdi-cellphone</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Devices</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
                <v-list-item v-if="$store.state.loggedIn" to="/settings">
                    <v-list-item-action>
                        <v-icon>mdi-tune</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Settings</v-list-item-title>
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
            <router-view/>
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
        created() {
            this.$root.$on('sessionExpired', this.logout)
        },
        destroyed() {
            this.$root.$off('sessionExpired', this.logout)
        },
        mounted() {
            this.checkToken()
        },
        methods: {
            checkToken() {
                if (localStorage.getItem("authToken") === null) {
                    return
                }
                this.$store.dispatch('setLoggedIn')
            },
            // TODO auto logout if server responds with 401
            logout() {
                localStorage.removeItem('authToken')
                this.$store.dispatch('setLoggedOut')
                this.$router.push({name: 'login'})
            }
        }
    });
</script>
