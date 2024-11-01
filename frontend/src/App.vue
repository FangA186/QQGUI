<script setup>
import a from './assets/images/lxy.webp'
import b from './assets/images/logo-universal.png'
import {ref} from "vue";
import {useRoute} from "vue-router";
import {useRouter} from "vue-router";
import V3Emoji from "vue3-emoji";
const backgroundImage = ref(a);
const route = useRoute();
const router = useRouter()
const changeBackgroundImage = (num) => {
  if (route.fullPath==='/'){
    num === 1 ? backgroundImage.value = a : backgroundImage.value = b
  }
};

</script>

<template>

  <div id="logo" :style="{ backgroundImage: `url(${backgroundImage})` }">
    <RouterView v-slot="{ Component, route }" @changeBackground="changeBackgroundImage">
      <transition name="j" mode="out-in">
        <component :is="Component" :key="route.path">
        </component>
      </transition>
    </RouterView>
  </div>
</template>

<style>
#logo {
  width: 100%;
  height: 100%;
  margin: auto;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
  display: flex;
  align-items: center;
}
</style>
