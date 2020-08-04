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

/**
 * knax-utils - small plugin for Knaxim to supply global utility functions.
 */
// import Vue from 'vue'
//
// const install = function (Vue) {
// }

export function humanReadableSize (size) {
  var unit = ' B'
  var amount = size
  if (amount > 700) {
    unit = ' KB'
    amount = amount / 1024
  }
  if (amount > 700) {
    unit = ' MB'
    amount = amount / 1024
  }
  if (amount > 700) {
    unit = ' GB'
    amount = amount / 1024
  }
  if (amount > 700) {
    unit = ' TB'
    amount = amount / 1024
  }
  return `${amount.toLocaleString(undefined, { maximumFractionDigits: 2 })}${unit}`
}

export function humanReadableTime (time) {
  var d = new Date(Date.parse(time))
  return `${d.toLocaleDateString()} ${d.toLocaleTimeString()}`
}

/*
 * To use the event bus:
 * import { EventBus } from '@/main'
 * EventBus.$emit('event-name', payload)
 * EventBus.$on('event-name', func)
 * EventBus.$off('event-name')
 * func (payload) {
 *   ...
 * }
 * https://alligator.io/vuejs/global-event-bus/
 *
 * Please be sure to keep track of events you emit within your components.
 */
// Deprecating in favor of new Vuex Event Bus module
// export const EventBus = new Vue()

export default {
  // install,
}
