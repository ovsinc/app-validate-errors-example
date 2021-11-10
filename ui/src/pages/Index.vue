<template>
  <div class="q-pa-md" style="max-width: 300px">
    <form @submit.prevent.stop="onSubmit" @reset.prevent.stop="onReset" class="q-gutter-md">
      <q-input
        filled
        v-model="login"
        label="Your name *"
        hint="Name and surname"
        name="login"
        :error-message="getValidationErrors('login')"
        :error="hasValidationErrors('login')"
      />

      <q-input
        filled
        v-model="oldPassword"
        label="Your name *"
        hint="Name and surname"
        name="oldPassword"
        :error-message="getValidationErrors('oldPassword')"
        :error="hasValidationErrors('oldPassword')"
      />

      <q-input
        filled
        v-model="newPassword"
        label="Your name *"
        hint="Name and surname"
        name="newPassword"
        :error-message="getValidationErrors('newPassword')"
        :error="hasValidationErrors('newPassword')"
      />

      <q-input
        filled
        v-model="newPasswordAgain"
        label="Your name *"
        hint="Name and surname"
        name="newPasswordAgain"
        :error-message="getValidationErrors('newPasswordAgain')"
        :error="hasValidationErrors('newPasswordAgain')"
      />

      <div>
        <q-btn label="Submit" type="submit" color="primary" />
        <q-btn label="Reset" type="reset" color="primary" flat class="q-ml-sm" />
      </div>
    </form>
  </div>
</template>

<script>
import { ref, defineComponent } from 'vue'
import { useQuasar } from 'quasar'

export default defineComponent({
  name: 'PageIndex',
  setup () {
    const $q = useQuasar()

    const login = ref(null)
    const oldPassword = ref(null)
    const newPassword = ref(null)
    const newPasswordAgain = ref(null)

    // const request = ref({})
    const response = ref({})

    function getValidationErrorMessages (field) {
      for (const [key, value] of Object.entries(response.value)) {
        if (key === field) {
          return value
        }
      }
      return []
    }

    function getValidationErrors (field) {
      const errors = getValidationErrorMessages(field)
      if (errors.length !== 0) {
        return errors.join('\r\n')
      }
      return ''
    }

    function hasValidationErrors (field) {
      if (getValidationErrorMessages(field).length !== 0) {
        showValidationError()
        return true
      }
      return false
    }

    function setValidationErrors (payload) {
      response.value = payload
    }

    function showValidationError () {
      $q.notify({
        type: 'negative',
        message: 'Validation failure',
        caption: 'please check the inputs'
      })
    }

    return {
      getValidationErrors,
      hasValidationErrors,

      login,
      oldPassword,
      newPassword,
      newPasswordAgain,

      computed: {
        isValid: function () {}
      },

      onSubmit () {
        const data = {
          login: ['hello world', 'is world']
        }
        setValidationErrors(data)
      },

      onReset () {
        login.value = null
        oldPassword.value = null
        newPassword.value = null
        newPasswordAgain.value = null
        setValidationErrors({})
      }
    }
  }
})
</script>
