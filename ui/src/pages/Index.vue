<template>
  <div
    class="q-pa-md"
    style="max-width: 500px"
  >
    <q-select
      v-model="locale"
      :options="localeOptions"
      class="q-pb-xl"
      dense
      borderless
      emit-value
      map-options
      options-dense
      style="min-width: 150px"
    />

    <form
      class="q-gutter-md"
      @submit.prevent.stop="onSubmit"
      @reset.prevent.stop="onReset"
    >
      <q-input
        v-model="login"
        filled
        :label="$t('login_label')"
        :hint="$t('login_hint')"
        name="login"
        :error-message="getValidationErrors('login')"
        :error="hasValidationErrors('login')"
      />

      <q-input
        v-model="oldPassword"
        filled
        :label="$t('oldPassword_label')"
        :hint="$t('oldPassword_hint')"
        name="oldPassword"
        :error-message="getValidationErrors('oldPassword')"
        :error="hasValidationErrors('oldPassword')"
      />

      <q-input
        v-model="newPassword"
        filled
        :label="$t('newPassword_label')"
        :hint="$t('newPassword_hint')"
        name="newPassword"
        :error-message="getValidationErrors('password')"
        :error="hasValidationErrors('password')"
      />

      <div class="q-pt-xl">
        <q-btn
          :label="$t('btn_submit')"
          type="submit"
          color="primary"
        />
        <q-btn
          :label="$t('btn_reset')"
          type="reset"
          color="primary"
          flat
          class="q-ml-sm"
        />
      </div>
    </form>
  </div>
</template>

<script>
import { ref, defineComponent } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { api } from 'boot/axios'

export default defineComponent({
  name: 'PageIndex',
  setup () {
    const $q = useQuasar()
    const { locale, t } = useI18n({ useScope: 'global' })

    locale.value = $q.lang.getLocale()

    const login = ref(null)
    const oldPassword = ref(null)
    const newPassword = ref(null)

    const errors = ref({})

    function changePasswordSend (req) {
      api.post('/api/v1', req)
        .then(
          (response) => {
            errors.value = {}

            if (response.data.success) {
              $q.notify({
                type: 'positive',
                position: 'top',
                message: response.data.message,
                caption: t('change_pass_success_caption')
              })
            }
          }
        )
        .catch(
          (error) => {
            let message = t('change_net_fail_message_default')

            if (error.response && error.response.data.errors) {
              errors.value = error.response.data.errors
              message = error.response.data.message
            } else {
              errors.value = {
                common: [message]
              }
            }

            $q.notify({
              type: 'negative',
              message: message,
              position: 'top',
              caption: t('change_pass_fail_caption')
            })
          }
        )
    }

    function getValidationErrorMessages (field) {
      for (const [key, value] of Object.entries(errors.value)) {
        if (key === field) {
          return value
        }
      }
      return []
    }

    return {
      login,
      oldPassword,

      locale,
      localeOptions: [
        { value: 'ru-RU', label: 'Russian' },
        { value: 'en-US', label: 'English' }
      ],

      getValidationErrors (field) {
        const errs = getValidationErrorMessages(field)
        if (errs.length !== 0) {
          return errs.join('; ')
        }
        return ''
      },

      hasValidationErrors (field) {
        return getValidationErrorMessages(field).length !== 0
      },

      onSubmit () {
        const req = {
          login: login.value,
          old_password: oldPassword.value,
          password: newPassword.value,
          lang: locale.value
        }
        changePasswordSend(req)
      },

      onReset () {
        login.value = null
        oldPassword.value = null
        newPassword.value = null
        errors.value = {}
      }
    }
  }
})
</script>
