<template>
    <v-app>
        <v-navigation-drawer
                v-model="leftDrawer"
                app
        >
            <v-list dense>
                <v-list-item v-if="$store.state.loggedIn" to="/files">
                    <v-list-item-action>
                        <v-icon>mdi-file-cabinet</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Files</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
                <v-list-item v-if="$store.state.loggedIn" to="/tags">
                    <v-list-item-action>
                        <v-icon>mdi-label</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Tags</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
                <v-list-item v-if="$store.state.loggedIn" to="/scan">
                    <v-list-item-action>
                        <v-icon>mdi-file-search-outline</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Scan</v-list-item-title>
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
                color="indigo lighten-1"
                dark
        >
            <v-app-bar-nav-icon @click.stop="leftDrawer = !leftDrawer"/>
            <v-toolbar-title>GoTag</v-toolbar-title>
            <v-spacer/>
            <v-app-bar-nav-icon @click.stop="rightDrawer = !rightDrawer"/>
        </v-app-bar>

        <v-content>
            <router-view/>
        </v-content>

        <v-navigation-drawer
                app
                v-model="rightDrawer"
                right
                fixed
        >
            <tag-list/>
        </v-navigation-drawer>
    </v-app>
</template>

<script lang="ts">
    import Vue from 'vue';
    import TagList from "./components/TagList.vue";

    export default Vue.extend({
        name: 'App',
        components: {TagList},
        data: () => ({
            leftDrawer: true,
            rightDrawer: true,
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
