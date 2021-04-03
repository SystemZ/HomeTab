<template>
  <v-app>
    <v-navigation-drawer
      v-model="leftDrawer"
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
        <v-list-item v-if="$store.state.loggedIn" to="/log">
          <v-list-item-action>
            <v-icon>mdi-clipboard-list-outline</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Log</v-list-item-title>
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
        <v-list-item v-if="$store.state.loggedIn" to="/files">
          <v-list-item-action>
            <v-icon>mdi-image</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Images</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item v-if="$store.state.loggedIn" to="/tags">
          <v-list-item-action>
            <v-icon>mdi-tag-multiple</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Image tags</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item v-if="$store.state.loggedIn" to="/scan">
          <v-list-item-action>
            <v-icon>mdi-image-search</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Scan images</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item v-if="$store.state.loggedIn" to="/pantry">
          <v-list-item-action>
            <v-icon>mdi-silverware-fork-knife</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Pantry</v-list-item-title>
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
      :color="navbarPrimary"
      :dark="componentDark"
    >
      <v-app-bar-nav-icon @click.stop="leftDrawer = !leftDrawer"/>
      <v-toolbar-title>HomeTab</v-toolbar-title>
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
    rightDrawer: false,
  }),
  created() {
    this.$root.$on('sessionExpired', this.logout)
    this.$root.$on('openRightDrawer', this.openRightDrawer)
    this.$root.$on('closeRightDrawer', this.closeRightDrawer)
  },
  destroyed() {
    this.$root.$off('sessionExpired', this.logout)
    this.$root.$off('openRightDrawer', this.openRightDrawer)
    this.$root.$off('closeRightDrawer', this.closeRightDrawer)
  },
  mounted() {
    this.checkToken()
  },
  methods: {
    openRightDrawer() {
      this.rightDrawer = true
    },
    closeRightDrawer() {
      this.rightDrawer = false
    },
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
