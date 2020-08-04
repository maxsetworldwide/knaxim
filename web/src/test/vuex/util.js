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
// helper for testing action with expected mutations
export const testAction = async (
  action,
  payload,
  expected = {},
  state = {},
  getters = {}
) => {
  expected = {
    ...{
      resolve: undefined,
      reject: undefined,
      mutations: [], // { type, payload }
      actions: [] // { type, payload }
    },
    ...expected
  }
  let mCount = 0
  let aCount = 0

  // mock commit
  const commit = (type, payload) => {
    if (mCount >= expected.mutations.length) {
      fail(`unexpected extra mutation: ${type}(${payload})`)
      return
    }
    const mutation = expected.mutations[mCount]

    try {
      expect(type).toBe(mutation.type)
      if ((payload || mutation.payload) && !mutation.ignorePayload) {
        expect(payload).toEqual(mutation.payload)
      }
    } catch (error) {
      fail(error)
    }
    mCount++
  }

  // mock dispatch
  const dispatch = (type, payload) => {
    if (aCount >= expected.actions.length) {
      fail(`unexpected extra action: ${type}(${payload})`)
      return
    }
    const otheraction = expected.actions[aCount]

    try {
      expect(type).toBe(otheraction.type)
      if ((otheraction.payload || payload) && !otheraction.ignorePayload) {
        expect(payload).toEqual(otheraction.payload)
      }
    } catch (err) {
      fail(err)
    }
    aCount++
    return Promise.resolve(true)
  }

  // call the action with mocks and arguments
  try {
    let result = await action({ commit, dispatch, state, getters }, payload)
    if (expected.resolve !== undefined) {
      expect(result).toEqual(expected.resolve)
    }
  } catch (err) {
    if (expected.reject !== undefined) {
      expect(err).toEqual(expected.reject)
    } else {
      fail(err)
    }
  }

  if (mCount !== expected.mutations.length) {
    fail('not all expected mutations called')
  }

  if (aCount !== expected.actions.length) {
    fail('not all expected actions called')
  }
}
