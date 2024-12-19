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
const showGroup = ref(false)
const search = async () => {
  const userid = localStorage.getItem("userID")
  const res = await SearchUser(userid,username.value)
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
  await HandleAddFriendRequest(agr1,agr2,false)
}
const qiehuan = ()=>{
  showGroup.value = !showGroup.value
  // console.log(showGroup.value)
  emit('clear',{message1:showGroup.value})
  // if (showGroup.value){
  //   emit('clear',{message:false})
  // }else {
  //   emit('clear',{message:true})
  // }
}
</script>
<template>
  <div class="search">
    <input type="text" maxlength="50" placeholder="搜索" style="box-shadow: 0 0 0.1vh white;width: 150px;" @keyup.enter="search" v-model="username" @click="clear1"/>
<!--    <svg class="icon iconsearch" aria-hidden="true">-->
<!--      <use xlink:href="#icon-sanjiaoxing"></use>-->
<!--    </svg>-->

    <div class="action-icon-container">
      <slot name="action-icon"></slot>
    </div>
    <svg class="icon iconqiehuan" aria-hidden="true" @click="qiehuan">
      <use xlink:href="#icon-sanjiaoxing-copy"></use>
    </svg>
    <ul v-if="showul" class="searchList">
      <li v-for="(item, index) in searchInfo" :key="index" @click="lookInfo(item)">
        <img :src="item.Avatar" alt=""/>
        <span>{{ item.Username }}</span>
      </li>
    </ul>
    <svg class="icon iconclear" aria-hidden="true" @click="clear" v-if="clearIcon">
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
  height: 40px;
  //background-color: #f7ba2a;
  width: 200px;
  //.iconsearch {
  //  position: absolute;
  //  left: 1px;
  //  top: 16px;
  //  width: 14px;  // 固定大小
  //  height: 14px; // 固定大小
  //}

  .iconqiehuan {
    position: absolute;
    right:-1px;
    top: 13px;
    width: 16px;  // 固定大小
    height: 16px; // 固定大小
  }

  input {
    position: absolute;
    border: none;
    outline: none;
    text-indent: 3px;
    height: 20px;
    width: 100px;
    left: 5px;
    top: 10px;
    line-height: 20px;
    opacity: 0.5;
    border-radius: 5px;
    font-size: 14px; // 固定字体大小
  }

  .action-icon-container {
    position: absolute;
    right: 25px;
    top: 10px;
    cursor: pointer;
    width: 14px;  // 固定大小
    height: 14px; // 固定大小
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
        width: 40px;  // 固定宽度
        height: 40px; // 固定高度
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
        font-size: 16px; // 固定字体大小
      }
    }
  }

  .iconclear {
    position: absolute;
    right: 42px;
    top: 10px;
    cursor: pointer;
    width: 24px;  // 固定大小
    height: 24px; // 固定大小
  }

  .dialog-footer {
    color: black;

    .header {
      text-align: left;
      position: relative;
      height: 6.5vh;
      padding: 0.5vw 0.5vw 0 0.5vw;

      img {
        width: 3vw;  // 固定宽度
        height: 3vw; // 固定高度
        position: absolute;
        top: 50%;
        left: 3%;
        transform: translateY(-50%);
      }

      span {
        font-size: 1vw; // 固定字体大小
      }

      span {
        position: absolute;
        top: 50%;
        left: 20%;
        transform: translateY(-50%);
      }
    }

    .Signature {
      text-align: left;
      text-indent: 0.6vw;
      font-size: 0.8vw;
    }
  }
}

</style>