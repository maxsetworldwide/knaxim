// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import Vue from 'vue'
import VueRouter from 'vue-router'

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
  },
  {
    path: '/reset/:passkey',
    name: 'reset',
    components: {
      default: FileList,
      sideview: MemberList
    }
  },
  {
    path: '/profile/newpassword',
    name: 'changepass',
    component: Profile
  },
  {
    path: '/search/:find',
    name: 'search',
    props: true,
    component: HeaderSearchList
  },
  {
    path: '/search/:find/acronym/:acr',
    name: 'searchWithAcronym',
    props: true,
    component: HeaderSearchList
  },
  {
    path: '/search/:find/tag/:tag',
    name: 'searchWithTag',
    props: true,
    component: HeaderSearchList
  },
  {
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
  },
  {
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
  },
  {
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
