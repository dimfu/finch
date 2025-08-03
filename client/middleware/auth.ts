import { useAuthStore } from "~/store/auth";

export default defineNuxtRouteMiddleware(async (to, from) => {
  const authStore = useAuthStore();
  if (from.path == "/signin" && authStore.isAuthenticated) {
    const tokenIsValid = await authStore.tokenIsValid();
    if (tokenIsValid) {
      return navigateTo("/");
    }
  }
});
