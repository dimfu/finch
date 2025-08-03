<template>
  <div class="md:hidden"></div>

  <div
    class="container relative hidden h-[800px] flex-col items-center justify-center md:grid lg:max-w-none lg:grid-cols-2 lg:px-0"
  >
    <div
      class="relative hidden h-full flex-col bg-muted p-10 text-white dark:border-r lg:flex"
    >
      <div class="absolute inset-0 bg-zinc-900" />
    </div>

    <form @submit.prevent="handleSubmit" class="lg:p-8">
      <div
        class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]"
      >
        <div class="flex flex-col space-y-2 text-center">
          <h1 class="text-2xl font-semibold tracking-tight">
            Sign in to your account
          </h1>
          <p class="text-sm text-muted-foreground">
            Enter your user and password below to sign in
          </p>
        </div>
        <div class="flex flex-col space-y-2">
          <Label>Username</Label>
          <Input v-model="username" placeholder="Insert username"></Input>
        </div>
        <div class="flex flex-col space-y-2">
          <Label>Password</Label>
          <Input
            v-model="password"
            type="password"
            placeholder="*******"
          ></Input>
        </div>
        <Button type="submit">Submit</Button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ["auth"],
});
import Label from "~/components/ui/label/Label.vue";
import Input from "~/components/ui/input/Input.vue";
import Button from "~/components/ui/button/Button.vue";

import { useAuthStore } from "~/store/auth";

const username = ref("");
const password = ref("");

const auth = useAuthStore();

async function handleSubmit() {
  try {
    const result = await auth.signin(username.value, password.value);
    if (result.success) {
      console.log(auth.token);
    }
  } catch (error) {
    console.error(error);
  }
}
</script>
