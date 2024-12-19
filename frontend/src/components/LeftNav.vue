<script setup>
import {useRouter} from "vue-router";
const router = useRouter()
</script>

<template>
  <ul class="navigation_list">
    <li>
     <slot name="img1"></slot>
    </li>
    <li @click="router.push('/Index')">
      <svg class="icon" aria-hidden="true" font-size="30">
        <use xlink:href="#icon-xinxi"></use>
      </svg>
    </li>
    <li @click="router.push('/Friends')">
      <svg class="icon" aria-hidden="true" font-size="30">
        <use xlink:href="#icon-haoyou1"></use>
      </svg>
    </li>
    <li @click="router.push('/Public')">
      <svg class="icon" aria-hidden="true" font-size="30">
        <use xlink:href="#icon-pengyouquan2"></use>
      </svg>
    </li>
  </ul>
</template>

<style scoped lang="less">
.navigation_list {
  padding: 0;
  //width: 100%;
  //background-color: #67c23a;
  li {
    list-style: none;
    text-align: center;
    margin-top: 3vh;
  }

  li:hover {
    cursor: pointer;
  }
}
</style>