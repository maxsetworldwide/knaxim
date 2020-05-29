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
