<script setup>
import Index from "../components/Index.vue";
import V3Emoji from 'vue3-emoji'
import 'vue3-emoji/dist/style.css'
import {onMounted, ref} from "vue";
import {
  ConsentFriend,
  GetApplicationRecord,
  GetConsentList,
  GetIsSpeakUserInfo,
  Receive
} from "../../wailsjs/go/api/MyApp.js"
// import { useMessageStore } from '../store/messages.js'; // 根据你的文件结构调整路径

// // 获取 messageStore
// const messageStore = useMessageStore();
//
// // 直接使用 messageStore.messages
// const messages = messageStore.messages;
// console.log(messages)
const consentList = ref([])
const dialogVisible = ref(false)
const ApplicationRecord = ref([])
const twoli = ref(true)
const speakList = ref([])
const oneinfo = ref([])
const show = ref(false)
const content = ref("")
let socket
const clickEmoji = (val) => {
  content.value = content.value + val
}
const consentFriend = async () => {
  const userID = localStorage.getItem("userID")
  const res = await GetConsentList(userID)
  const res1 = await GetApplicationRecord(userID)
  ApplicationRecord.value = res
  consentList.value = res1
}
const dispose = (num) => {
  dialogVisible.value = !dialogVisible.value
  if (num === 1) {
    twoli.value = true
  } else {
    twoli.value = false
  }
}
const consent = async (item) => {
  const userID = localStorage.getItem("userID")
  const res = await ConsentFriend(userID, String(item.applicant_id), "1")
}
const refuse = async (item) => {
  const userID = localStorage.getItem("userID")
  const res = await ConsentFriend(userID, String(item.applicant_id), "2")
}
const getisspeakuserinfo = async () => {
  speakList.value = await GetIsSpeakUserInfo(localStorage.getItem("userID"))
}
const GenerateRoomID = (uuid1,uuid2)=>{
  if (uuid1<uuid2){
    return uuid1 + "_" + uuid2
  }
  return uuid2 + "_" + uuid1
}
const sendMessage = async () => {
  console.log(oneinfo.value)
  await Receive(localStorage.getItem("userID"),
      JSON.stringify(oneinfo.value.friend_id),
      oneinfo.value.friend_avatar,
      oneinfo.value.friend_username,
  )
  socket.send(JSON.stringify({
    content:content.value,
    send_user_id:parseInt(localStorage.getItem("userID")),
    receiver_user_id:oneinfo.value.friend_id,
    file_url: "",
    room_id:GenerateRoomID(localStorage.getItem("uuid"),oneinfo.value.friend_uuid)
  }))
  content.value = ""
}
const dialogue = (item) => {
  oneinfo.value = item
  show.value = true
  const userId = localStorage.getItem("userID"); // 用户 ID
  const userIDUUID = localStorage.getItem("uuid");
  const friendIDUUID = item.friend_uuid
  const isGroup = 0; // 是否为群聊
  const wsUrl = `ws://127.0.0.1:8080/ws?user_id=${userId}&userIDUUID=${userIDUUID}
  &friendIDUUID=${friendIDUUID}&IsGroup=${isGroup}&roomid=${GenerateRoomID(userIDUUID,friendIDUUID)}`;
  socket = new WebSocket(wsUrl);

  socket.onopen = function () {
    console.log("连接已建立");
  };

  socket.onmessage = function (event) {
    console.log(event.data)
  };

  socket.onclose = function () {
    console.log("连接已关闭");
  };

  socket.onerror = function (error) {
    console.error("WebSocket 错误:", error);
  };
}
onMounted(() => {
  consentFriend()
  getisspeakuserinfo()
})
</script>

<template>
  <Index>
    <template #two>
      <ul class="searchList">
        <li style="cursor: pointer;" v-for="(item,index) in speakList" @click="dialogue(item)">
          <img :src="item.friend_avatar" alt="">
          <span style="text-indent: 0.5vw;">{{ item.friend_username }}</span>
        </li>
        <li @click="dispose(1)" style="cursor: pointer">
          好友申请
        </li>
        <li @click="dispose(2)" style="cursor: pointer">
          申请记录
        </li>

      </ul>
      <el-dialog v-model="dialogVisible" width="400">
        <template #footer>
          <div class="dialog-footer">
            <ul class="searchList1">
              <li v-for="(item,index) in consentList" style="display: block" v-if="twoli===true">
                <img :src="item.applicant_avatar" alt="">
                <span style="text-indent: 0.5vw">{{ item.applicant_username }}
              </span>
                <div v-if="item.IsConsent===0">
                  <button style="width: 3vw" @click="consent(item)">
                    同意
                  </button>
                  <button style="position: absolute;right: 0;top: 50%;left: 50%;transform: translateY(-50%);width: 3vw"
                          @click="refuse(item)">
                    拒绝
                  </button>
                </div>
                <div v-else-if="item.IsConsent===1">
                  <button style="width: 3.5vw" disabled>已同意</button>
                </div>
                <div v-else>
                  <button style="width: 3.5vw" disabled>已拒绝</button>
                </div>
              </li>
              <li v-for="(item,index) in ApplicationRecord" style="display: block" v-else>
                <img :src="item.friend_avatar" alt="">
                <span style="text-indent: 0.5vw">{{ item.friend_username }}
              </span>
                <div v-if="item.IsConsent===0">
                  <button style="width: 5vw;" disabled>
                    已发送验证
                  </button>
                </div>
                <div v-else-if="item.IsConsent===1">
                  <button style="width: 5vw" disabled>已同意</button>
                </div>
                <div v-else>
                  <button style="width: 5vw" disabled>已拒绝</button>
                </div>
              </li>
            </ul>
          </div>
        </template>
      </el-dialog>
    </template>
    <template #three v-if="show">
      <div class="dialogueList">
        <h2>{{ oneinfo.friend_username }}</h2>
        <ul class="list">
          <li class="left">
            <img :src="oneinfo.friend_avatar" alt="">
            <div>456</div>
          </li>
          <li class="right">
            <div>45asdadadasdadasdasdasdas6</div>
            <img :src="oneinfo.friend_avatar" alt="">
          </li>
        </ul>
        <div class="di">
          <V3Emoji :recent="true" class="xl" @click-emoji="clickEmoji"/>
          <svg class="icon wj" aria-hidden="true" font-size="0.8vw">
            <use xlink:href="#icon-wenjian"></use>
          </svg>
          <textarea id="story" name="story" rows="5" cols="33" v-model="content"></textarea>
          <button @click="sendMessage">发送(enter)</button>
        </div>
      </div>
    </template>
  </Index>
</template>

<style scoped lang="less">
.searchList, .searchList1 {
  width: 100%;
  padding: 0;
  margin: 0;

  li {
    text-align: left;
    list-style: none;
    padding: 0.5vh;
    position: relative;
    height: 5vh;
    line-height: 5vh;
    color: black;

    div {
      width: 15vh;
      height: 5vh;
      position: absolute;
      right: 0
    }

    img {
      width: 13%;
      border-radius: 0.5vh;
      position: absolute;
      top: 50%;
      height: 4vh;
      image-rendering: -moz-crisp-edges; /* Firefox */
      image-rendering: -o-crisp-edges; /* Opera */
      image-rendering: -webkit-optimize-contrast;
      image-rendering: crisp-edges;
      -ms-interpolation-mode: nearest-neighbor; /* IE (non-standard property) */
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

.searchList {
  height: 95vh;
  background-color: #d9d8d8;

  li:hover {
    opacity: 0.7;
    background-color: #72767b;
  }
}

.dialogueList {
  height: 100vh;

  h2 {
    text-align: left;
    text-indent: 1vw;
    height: 8vh;
    line-height: 8vh;
  //background-color: #f7ba2a; margin: 0;
    color: black;
    border-bottom: 1px solid black;
  }

  .list {
    height: 60vh;
    padding: 0;
    margin: 0;

    .left {
      list-style: none;
      text-align: left;
      padding-left: 2vw;
      display: flex;
      align-items: center;
      margin-top: 2vh;

      img {
        width: 3vw;
        height: 4vh;
      }

      div {
        margin-left: 0.5vw;
        color: black;
      }
    }

    .right {
      list-style: none;
      text-align: right;
      padding-right: 2vw;
      height: 5vh;
      display: flex;
      justify-content: flex-end;
      align-items: center;

      img {
        width: 3vw;
        height: 4vh;
      }

      div {
        margin-right: 0.5vw;
        color: black;
      }
    }
  }

  .di {
    position: relative;

    .xl {
      width: 4vw;
      height: 0;
      position: fixed;
    }

    .wj {
      position: absolute;
      left: 4vw;
      font-size: 1.8vw;
    }

    textarea {
      letter-spacing: 1px;
      margin-top: 6vh;
      width: 98%;
      line-height: 1.5;
      resize: none;
      font-size: 1vw;
      border: none;
      color: black;
      font-weight: 700;
      font-family: Verdana, serif;
      outline: none;
      caret-color: black;
      text-indent: 0.1vw;
      height: 19vh;
      background-color: #cfd9df;
      border-top: 1px solid black;
    }

    button {
      position: absolute;
      bottom: 1vh;
      right: 2vw;
      border: none;
      width: 7vw;
      height: 5vh;
      color: green;
      font-weight: 700;
      border-radius: 0.5vw;
      cursor: pointer;
    }

    button:active {
      background-color: #cfd9df;
    }
  }

}
</style>