import Vue from 'vue'
import VueRouter from 'vue-router'

import Auth from '@/components/auth'
import FileList from '../components/file-list'
// import TextView from '../components/text-viewer'
import FileView from '../components/file-viewer.vue'
import HeaderSearchList from '../components/header-search-list'
import AppInfo from '../components/app-info'
import MemberList from '@/components/member-list'
import Profile from '@/components/profile'

Vue.use(VueRouter)

const routes = [
  {
    path: '',
    name: 'home',
    components: {
      default: FileList,
      sideview: MemberList
    }
  }, {
    path: '/login',
    name: 'login',
    component: Auth
  }, {
    path: '/register',
    name: 'register',
    component: Auth
  }, {
    path: '/request',
    name: 'request',
    component: Auth
  }, {
    path: '/profile/newpassword',
    name: 'changepass',
    component: Profile
  }, {
    path: '/search/:find',
    name: 'search',
    props: true,
    component: HeaderSearchList
  }, {
    path: '/search/:find/acronym/:acr',
    name: 'searchWithAcronym',
    props: true,
    component: HeaderSearchList
  }, {
    path: '/file/:id',
    name: 'file',
    components: {
      default: FileView,
      sideview: AppInfo
    },
    props: {
      default: true,
      sideview: true
    }
  }, {
    path: '/file/:id/acronym/:acr',
    name: 'fileWithAcronym',
    components: {
      default: FileView,
      sideview: AppInfo
    },
    props: {
      default: true,
      sideview: true
    }
  }, {
    path: '/list/:src',
    name: 'filteredFiles',
    components: {
      default: FileList,
      sideview: MemberList
    },
    props: {
      default: true
    }
  }
  // {
  //   path: '/team/:gid',
  //   name: 'teamHome',
  //   components: {
  //     default: FileList,
  //     sideview: MemberList
  //   },
  //   props: {
  //     default: true,
  //     sideview: true
  //   }
  // }, {
  //   path: '/team/:gid/:src',
  //   name: 'teamFiltered',
  //   components: {
  //     default: FileList,
  //     sideview: MemberList
  //   },
  //   props: {
  //     default: true,
  //     sideview: true
  //   }
  // }
]

const router = new VueRouter({
  mode: 'history',
  linkActiveClass: 'active',
  routes: routes
})

export default router
