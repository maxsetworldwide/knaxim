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

export default {
  // install,
}
