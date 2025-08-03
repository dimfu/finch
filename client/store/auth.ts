import { defineStore } from "pinia";

interface User {
  id: string;
}

interface AuthState {
  user?: User;
  hasToken: boolean;
  isAuthenticated: boolean;
}

export const useAuthStore = defineStore("auth", {
  state: (): AuthState => ({
    user: undefined,
    hasToken: false,
    isAuthenticated: false,
  }),

  actions: {
    setUser(userId: string) {
      this.user = {
        id: userId,
      };
      this.isAuthenticated = true;
    },

    async signin(username: string, password: string) {
      const config = useRuntimeConfig();
      try {
        const response = await $fetch(
          `${config.public.authUrl}/api/auth/signin`,
          {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: {
              username,
              password,
            },
            credentials: "include",
          }
        );

        return navigateTo("/", { external: true });
      } catch (error: any) {
        return navigateTo("/signin");
      }
    },

    async signout() {
      const config = useRuntimeConfig();
      try {
        const response = await $fetch(
          `${config.public.authUrl}/api/auth/signout`,
          {
            method: "GET",
            credentials: "include",
          }
        );
        // clear the auth store states
        this.$reset;
      } catch (error: any) {
        console.error(error);
      }
      return navigateTo("/", { external: true });
    },

    async checkAuth() {
      // do nothing if there is no tokens inside http-coochie
      if (!this.hasToken) {
        return { error: "No tokens provided", success: false };
      }

      const config = useRuntimeConfig();

      try {
        const response = await $fetch<{ userId: string }>(
          `${config.public.authUrl}/api/auth/me`,
          {
            method: "GET",
            credentials: "include",
          }
        );
        this.setUser(response.userId);
        return { error: undefined, success: true };
      } catch (error: any) {
        return { error, success: false };
      }
    },
  },
});
