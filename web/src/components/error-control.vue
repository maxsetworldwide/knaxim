<script>
/**
 * ErrorController is intended to receive global and local errors and organize
 * them to be presented to the user in various ways.  Or shoveled off into
 * another service for analysis.
 *
 * It should not be directly coupled with the creation of errors.
 *
 * Inputs
 *  vm.$refs.[ErrorController].setError(err)
 *  vm.$root.$emit('app::set::error', err)
 *  #default="{ setError }"  // Slot Prop
 *
 * Outputs:
 *  vm.$refs.[ErrorController].lastError
 *  vm.$on.error(err)
 *  #default="{ lastError }"  // Slot Prop
 */
export default {
  name: 'error-control',
  props: {
    // Listen for global error messages.
    global: {
      type: Boolean,
      defaul: false
    }
  },

  data () {
    return {
      // TODO: Track error history for debugging purposes...
      errors: [],
      lastError: ''
    }
  },
  mounted () {
    // Listen for and handle any global error messages.
    if (this.global) {
      this.$root.$on('app::set::error', (data) => {
        this.setError(data)
      })
    }
  },

  methods: {
    /* Handle all incoming errors; notify the parent.
     *
     * Use setError in the parent with $refs or anywhere in the slot using
     * the returned slot function.
    */
    setError (error) {
      this.lastError = error
      this.$emit('error', this.lastError)
    }
  },

  render () {
    return this.$scopedSlots.default({
      setError: this.setError,
      lastError: this.lastError
    })
  }
}
</script>
