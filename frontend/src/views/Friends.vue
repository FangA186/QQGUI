<script setup>
import Index from "../components/Index.vue";
import {onMounted, ref} from "vue";
import {useRouter} from "vue-router";
import {FriendList, IsSpeak, Login} from "../../wailsjs/go/api/MyApp.js"
const friendsList = ref([])
const dialogVisible = ref(false)
const oneinfo = ref()
const router = useRouter()
const friends = async () => {
  const res = await FriendList(localStorage.getItem("userID"))
  friendsList.value = res.InfoList
}
const send = (item)=>{
  dialogVisible.value = true
  oneinfo.value = item
}
const sendMessage =async ()=>{
const res =await IsSpeak(localStorage.getItem("userID"),oneinfo.value.id)
  if (res.status ===200){
    await router.push("/Index")
  }
}
onMounted(()=>{
  friends()
})
</script>

<template>
  <Index>
    <template #two>
      <ul class="searchList">
        <li style="cursor: pointer;" v-for="(item,index) in friendsList" @click="send(item)">
          <img :src="item.avatar" alt="">
          <span style="text-indent: 0.5vw;color: white">{{ item.username }}</span>
        </li>
      </ul>
    </template>
    <template #three>
      <el-dialog
          v-model="dialogVisible"
          width="400"
      >
        <template #footer>
          <div class="dialog-footer">
            <div class="header">
              <img :src="oneinfo.avatar" alt="">
              <span>{{oneinfo.username}}</span>
            </div>
            <div class="Signature">个性签名:{{oneinfo.signature}}</div>
          </div>
          <div style="margin-top: 1vh">
            <el-button @click="sendMessage">发消息</el-button>
            <el-button>
              语音通话
            </el-button>
            <el-button type="primary" @click="dialogVisible = false">
              视频通话
            </el-button>
          </div>
        </template>
      </el-dialog>
    </template>
  </Index>
</template>

<style scoped lang="less">
.searchList{
  width: 100%;
  padding: 0;
  background-color: #2d2d2d;
  margin: 0;
  li {
    text-align: left;
    list-style: none;
    padding: 0.5vh;
    position: relative;
    height: 5vh;
    line-height: 5vh;
    color: black;

    img {
      width: 13%;
      //border-radius: 1vh;
      position: absolute;
      left: 6%;
      top: 50%;
      image-rendering: -moz-crisp-edges; /* Firefox */
      image-rendering: -o-crisp-edges; /* Opera */
      image-rendering: -webkit-optimize-contrast;image-rendering: crisp-edges;
    -ms-interpolation-mode: nearest-neighbor; /* IE (non-standard property) */
      height: 5vh;
      transform: translateY(-50%);
    }

    span {
      position: absolute;
      top: 35%;
      left: 19%;
      transform: translateY(-50%);
    }
  }

}
.dialog-footer{
  color: black;
  .header{
    text-align: left;
    position: relative;
    height: 6.5vh;
    padding: 0.5vw 0.5vw 0 0.5vw;
    img{
      width: 3vw;
      position: absolute;
      top: 50%;
      left: 3%;
      height: 5vh;
      transform: translateY(-50%);
    }
    span{
      font-size: 1vw;
    }
    span{
      position: absolute;
      top: 50%;
      left: 20%;
      transform: translateY(-50%);
    }
  }
  .Signature{
    text-align: left;
    text-indent: 0.6vw;
    font-size: 0.8vw;
  }
}
</style>