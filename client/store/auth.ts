import { defineStore } from "pinia";

interface User {
  id: string;
}

interface Token {
  accessToken: string;
  refreshToken: string;
}

interface AuthState {
  user?: User;
  token?: Token;
  isAuthenticated: boolean;
}

export const useAuthStore = defineStore("auth", {
  state: (): AuthState => ({
    user: undefined,
    isAuthenticated: false,
    token: undefined,
  }),

  actions: {
    setToken(token: Token) {
      this.token = token;
    },

    // to check if token is expired or broken
    async tokenIsValid(): Promise<boolean> {
      return true;
    },

    async signin(username: string, password: string) {
      const config = useRuntimeConfig();
      console.log(config.public.authUrl);
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

        return { error: undefined, success: true };
      } catch (error: any) {
        return { error, success: false };
      }
    },

    async signout() {
      const config = useRuntimeConfig();
      try {
        const response = await $fetch(
          `${config.public.authUrl}/api/auth/signin`,
          {
            method: "GET",
            credentials: "include",
          }
        );
        return { error: undefined, success: true };
      } catch (error: any) {
        return { error, success: false };
      }
    },
  },
});
