<script setup>
import {ref} from "vue";
import {useRouter} from "vue-router";
import JSEncrypt from 'jsencrypt';
import { useUserStore } from '../store/userStore.js';
import {Login,RegisterFunc} from "../../wailsjs/go/api/MyApp.js"
const userStore = useUserStore();
const router = useRouter()
const emit = defineEmits(['changeBackground']);
const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBCgKCAQEAqjiVu5aN2oHFwUOcZoCehIqO7pprF7C/CLKipxRmhJcqyjY/STum
k7oMaHkuD3iR+4hGhHkFLoLVdEBLHTGH/LOCH+nwWx2Jc4y/vhI9XWORoNKLx0p/
85LMTR4pHhbWKWZh9KzPQf5Qu/A28d9tohhkDA+0Gk4eGMK2HOV0aA7KcdzFoM1+
6EovTl4EtYErzVdEyyeS5vqZmRzXFCjO3oaQmue/80YrXmnjC4xVJDlAZKt80mn6
9kBEPzpZ+9jBQ9l14+IeW2Enl7EaISCp5MixbT4jMt4ycq4LiY3V1u6uKNWclRR6
GfQKBE97tVxLpYEWj7gYcm4Wv3clBO7+wwIDAQAB
-----END PUBLIC KEY-----`;
const forminfo = ref({
  Username: '',
  Password: ''
})


const key = ()=>{
  const encryptor = new JSEncrypt();
  encryptor.setPublicKey(publicKey);
  const encryptedPassword = encryptor.encrypt(forminfo.value.Password);
  return {
    Username: forminfo.value.Username,
    Password: encryptedPassword, // 使用加密后的密码
  }
}

const onRegisterClick = async (num) => {
  const ke = key()
  if (forminfo.value.Username !== '' && forminfo.value.Password !== '') {
    try {
      // const res = await axios.post("/api/register", key());
      const res = await RegisterFunc(ke.Username,ke.Password)
      if (res.status === 200) {
        forminfo.value.Username = ''
        forminfo.value.Password = ''
        emit('changeBackground', num);
      } else {
        console.error(res.data.message);
      }
    } catch (error) {
      console.error("注册请求失败", error);
    }
  } else {
    emit('changeBackground', num);
  }
};

const onLoginClick = async (num) => {
  const ke = key()
  if (forminfo.value.Username !== '' && forminfo.value.Password !== '') {
    try {
      const res = await Login(ke.Username,ke.Password)
      console.log(res)
      if (res.status === 200) {
        localStorage.setItem("userID",res.userID)
        localStorage.setItem("uuid",res.UUID)
        localStorage.setItem("avatar",res.useravatar)
        userStore.setUserId(res.userID);
        emit('changeBackground', num)
        await router.push("/index")
      } else {
        console.log("账户名或密码不正确")
      }
    } catch (error) {
      console.error("登录失败", error);
    }
  } else {
    emit('changeBackground', num)
  }
}
</script>

<template>
  <div id="form">
    <form>
      <label for="Username">
        <input type="text" id="Username" placeholder="Username" v-model="forminfo.Username">
      </label>
      <label for="">
        <input type="Password" id="pwd" placeholder="Password" v-model="forminfo.Password">
      </label>
    </form>
    <div class="login_register">
      <button style="margin-right: 1vw" @click="onLoginClick(1)">login</button>
      <button style="margin-left: 1vw" @click="onRegisterClick(2)">register</button>
    </div>
    <div id="QR">
      <svg class="icon" aria-hidden="true" font-size="30">
        <use xlink:href="#icon-erweima"></use>
      </svg>
    </div>
  </div>
</template>

<style scoped lang="less">
#form {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 50%;
  width: 50%;
  background-color: #cfd9df;
  border-radius: 1vh;
  opacity: 0.84;
  margin-left: 1vw;
  position: relative;

  form {
    display: flex;
    flex-direction: column;
    width: 100%;

    input {
      border: none;
      outline: none;
      background-color: #040817;
      opacity: 0.8;
      width: 50%;
      height: 7.5vh;
      margin-top: 3vh;
      caret-color: #cfd9df;
      color: #cfd9df;
      border-radius: 1vh;
    }
  }

  .login_register {
    margin-top: 1vh;

    button {
      border: none;
      width: 9vw;
      height: 5vh;
      border-radius: 1vh;
      transition: all 1s;
    }

    button:hover {
      box-shadow: 0 0 0 1px black;
      transition: all 1s;
    }
  }

  #QR {
    position: absolute;
    bottom: 1vh;
    right: 1vw;
  }
}
</style>