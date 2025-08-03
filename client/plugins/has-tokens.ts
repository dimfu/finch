import { useAuthStore } from "~/store/auth";

export default defineNuxtPlugin(async (app) => {
  const accessToken = useCookie<string | undefined>("access_token");
  const refreshToken = useCookie<string | undefined>("refresh_token");
  const authStore = useAuthStore();

  // initialize user's token if available on the http cookie
  if (accessToken.value && refreshToken.value) {
    authStore.$state.hasToken = true;
  }
});
