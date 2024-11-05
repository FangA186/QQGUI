<script setup>
import LeftNav from "./LeftNav.vue";
import Search from "./Search.vue";
import {onMounted, ref} from "vue";
import {GetIsSpeakUserInfo} from "../../wailsjs/go/api/MyApp.js"
import V3Emoji from "vue3-emoji";
const showmessages = ref(true)
// 处理接收到的事件
function handleClear(data) {
  showmessages.value = data.message
}
const avatar = ref('')

onMounted(()=>{
  avatar.value = localStorage.getItem("avatar")
})
</script>

<template>
  <div id="index">
    <div class="navigation_bar">
      <div class="img_div">
        <img :src="avatar" alt="">
      </div>
      <LeftNav></LeftNav>
    </div>
    <div class="chat_search">
      <Search @clear="handleClear">
        <template #action-icon>
          <svg class="icon iconadd" aria-hidden="true" font-size="25">
            <use xlink:href="#icon-jiahao"></use>
          </svg>
        </template>
      </Search>
      <div class="messages" v-if="showmessages">
        <slot name="two"></slot>
      </div>
    </div>
    <div class="chat_frame">
      <slot name="three"></slot>
    </div>
  </div>
</template>

<style scoped lang="less">
#index {
  width: 100%;
  height: 100%;
  background-color: #cfd9df;
  display: flex;
  .navigation_bar {
    width: 3%;
    height: 100%;
    background-color: #2e2e2e;

    .img_div {
      width: 100%;

      img {
        width: 80%;
        height: 2.3vw;
        margin-top: 1vh;
        image-rendering: -moz-crisp-edges; /* Firefox */
        image-rendering: -o-crisp-edges; /* Opera */
        image-rendering: -webkit-optimize-contrast;image-rendering: crisp-edges;
        -ms-interpolation-mode: nearest-neighbor; /* IE (non-standard property) */
      }
    }
  }

  .chat_search {
    display: flex;
    flex-direction: column;
    width: 18%;
    background-color: #f7f7f7;
    .messages {
      height: 95vh;
    }
  }

  .chat_frame {
    display: flex;
    flex-direction: column;
    width: 79%;
  }
}
</style>