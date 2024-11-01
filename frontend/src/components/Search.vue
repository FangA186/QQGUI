<script setup>
import {ref} from "vue";
import {SearchUser,HandleAddFriendRequest} from "../../wailsjs/go/api/MyApp.js"
const username = ref('')
const searchInfo = ref()
const showul = ref(false)
const clearIcon = ref(false)
const emit = defineEmits(['clear']);
const oneinfo = ref() // 存储点击的用户
const dialogVisible = ref(false)
const search = async () => {
  const userid = localStorage.getItem("userID")
  const res = await SearchUser(userid,username.value)
  // const info = await axios.get("/api/searchuser", {
  //   params: {
  //     userid,
  //     username: username.value
  //   }
  // })
  console.log(res)
  searchInfo.value = res
  res.length > 0 ? showul.value = true : showul.value = false
}
// 定义发出的事件
const clear = () => {
  username.value = ''
  clearIcon.value = false
  searchInfo.value = {}
  emit('clear', {message: true})
}

const clear1 = () => {
  clearIcon.value = true
  emit('clear', {message: false})
}
const lookInfo = (item) => {
  oneinfo.value = item
  console.log(oneinfo.value)
  dialogVisible.value = true
}
const addFriend = async ()=>{
  const agr1 = localStorage.getItem("userID")
  const agr2 = JSON.stringify(oneinfo.value.ID) // 被申请ID
  const res = await HandleAddFriendRequest(agr1,agr2,false)
  // const res = await axios.post("/api/addfriend",{
  //   user_id:localStorage.getItem("userID"),
  //   friend_id:JSON.stringify(oneinfo.value.ID),
  //   is_group:false
  // })
  console.log(res)
}
</script>
<template>
  <div class="search">
    <input type="text" maxlength="50" placeholder="搜索" style="box-shadow: 0 0 0.1vh white" @keyup.enter="search" v-model="username" @click="clear1"/>
    <svg class="icon iconsearch" aria-hidden="true" font-size="15">
      <use xlink:href="#icon-sousuo"></use>
    </svg>
    <div class="action-icon-container">
      <slot name="action-icon"></slot>
    </div>
    <ul v-if="showul" class="searchList">
      <li v-for="(item, index) in searchInfo" :key="index" @click="lookInfo(item)">
        <img :src="item.Avatar" alt=""/>
        <span>{{ item.Username }}</span>
      </li>
    </ul>
    <svg class="icon iconclear" aria-hidden="true" font-size="25" @click="clear" v-if="clearIcon">
      <use xlink:href="#icon-cuowu1"></use>
    </svg>
    <el-dialog
        v-model="dialogVisible"
        width="400"
    >
      <template #footer>
        <div class="dialog-footer">
<!--          <el-button>发消息</el-button>-->
          <div class="header">
            <img :src="oneinfo.Avatar" alt="">
            <span>{{oneinfo.Username}}</span>
          </div>
          <div class="Signature">个性签名:{{oneinfo.Signature}}</div>
        </div>
        <el-button v-if="oneinfo.is_friend!==1" style="margin-top: 1vh" @click="addFriend">添加好友</el-button>
        <div v-if="oneinfo.is_friend===1" style="margin-top: 1vh">
          <el-button>
            语音通话
          </el-button>
          <el-button type="primary" @click="dialogVisible = false">
            视频通话
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="less">
.search {
  position: relative;
  height: 5vh;
  background-color: #f7f7f7;
  .iconsearch {
    position: absolute;
    left: 1.3vw;
    top: 1.6vh;
  }

  input {
    position: absolute;
    border: none;
    outline: none;
    text-indent: 1.3vw;
    height: 3vh;
    width: 13vw;
    left: 1vw;
    top: 1vh;
    line-height: 3vh;
    background-color: #e2e2e2;
    border-radius: 0.5vh;
  }

  .action-icon-container {
    position: absolute;
    right: 2vw;
    top: 1.1vh;
    cursor: pointer;
  }

  .searchList {
    width: 100%;
    margin-top: 5vh;
    padding: 0;
    border-radius: 1vh;
    height: 95vh;

    li {
      text-align: left;
      list-style: none;
      padding: 0.5vh;
      position: relative;
      height: 5vh;
      margin-top: 0.5vh;
      background-color: #f7f7f7;
      border-radius: 1vh;

      img {
        width: 13%;
        border-radius: 1vh;
        position: absolute;
        top: 50%;
        transform: translateY(-50%);
      }

      span {
        color: black;
        position: absolute;
        top: 50%;
        left: 15%;
        transform: translateY(-50%);
      }
    }

  }

  .iconclear {
    position: absolute;
    right: 4vw;
    top: 1.1vh;
    cursor: pointer;
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
}
</style>