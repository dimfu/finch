import { useAuthStore } from "~/store/auth";

export default defineNuxtRouteMiddleware(async (to, from) => {
  // ignore if its coming from the server
  if (import.meta.server) return;
  const authStore = useAuthStore();

  // check whether the client is already signed in or not
  const { error, success } = await authStore.checkAuth();
  if (!success) {
    console.warn(error);
    return;
  }

  if (from.path == "/signin" && authStore.isAuthenticated) {
    // TODO: redirect the client to the main route for the protected route
    return navigateTo("/", { external: true });
  }
});
