
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
