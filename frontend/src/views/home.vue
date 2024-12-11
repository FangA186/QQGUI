<script setup>
import Index from "../components/Index.vue";
import V3Emoji from 'vue3-emoji'
import 'vue3-emoji/dist/style.css'

import {nextTick, onMounted, ref, watch} from "vue";
import {
  ConsentFriend,
  GetApplicationRecord,
  GetConsentList,
  GetIsSpeakUserInfo,
  MessageList,
  QueryRoom,
  Receive,
    GetGroupMessageList
} from "../../wailsjs/go/api/MyApp.js"
import { useMessageStore } from '../store/messages.js'; // 根据你的文件结构调整路径
const messageStore = useMessageStore();
const consentList = ref([])
const dialogVisible = ref(false)
const ApplicationRecord = ref([])
const twoli = ref(true)
const speakList = ref([])
const oneinfo = ref([])
const show = ref(false)
const content = ref("")
// let socket
let groupSocket
const messageList = ref([])
const avatar = ref('')
const userID = ref('')
const left = ref('left')
const right = ref('right')
const groupList = ref([])
const GroupOrUser = ref(false)
const GroupRoomID = ref("")
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
  twoli.value = num === 1;
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
  console.log(speakList.value)
}
const GenerateRoomID = (uuid1, uuid2) => {
  if (uuid1 < uuid2) {
    return uuid1 + "_" + uuid2
  }
  return uuid2 + "_" + uuid1
}
const sendMessage = async () => {

  if (GroupOrUser.value===false){
    console.log(112233)
    await Receive(localStorage.getItem("userID"),
        JSON.stringify(oneinfo.value.friend_id),
        oneinfo.value.applicant_avatar,
        oneinfo.value.applicant_username,
    )
    socket.send(JSON.stringify({
      two: 2,
      content: content.value,
      send_user_id: parseInt(localStorage.getItem("userID")),
      receiver_user_id: oneinfo.value.friend_id,
      file_url: "",
      room_id: GenerateRoomID(localStorage.getItem("uuid"), oneinfo.value.friend_uuid)
    }))
  }else {
    console.log(7778899)
    groupSocket.send(JSON.stringify({
      two:2,
      content:content.value,
      send_user_id: parseInt(localStorage.getItem("userID")),
      receiver_user_id:0,
      room_id:GroupRoomID.value
    }))
  }

  content.value = ""
}
let sockets = {};  // 用来存储不同房间或对话的 WebSocket 连接

const dialogue = async (item) => {
  GroupOrUser.value = false;
  oneinfo.value = item;
  show.value = true;
  console.log(oneinfo.value);

  const userIDUUID = localStorage.getItem("uuid");
  const friendIDUUID = item.friend_uuid;
  const isGroup = 0; // 是否为群聊
  const roomid = GenerateRoomID(userIDUUID, friendIDUUID);

  // 检查该房间是否已经有 WebSocket 连接
  if (sockets[roomid] && sockets[roomid].readyState === WebSocket.OPEN) {
    // 如果该房间的 WebSocket 已经连接，则无需重新连接
    messageList.value = await GetGroupMessageList(roomid);
    console.log(messageList.value)
    console.log("该房间的 WebSocket 已经连接，无需重新连接");
    return;
  }

  // 建立新的 WebSocket 连接
  const wsUrl = `ws://127.0.0.1:8080/ws?&userIDUUID=${userIDUUID}&IsGroup=${isGroup}&roomid=${roomid}`;
  const socket = new WebSocket(wsUrl);

  // 将新连接保存在 sockets 对象中，按 roomid 管理
  sockets[roomid] = socket;

  socket.onopen = async function () {
    messageList.value = await GetGroupMessageList(roomid);
    console.log(messageList.value);
    scrollToBottom();
  };

  socket.onmessage = function (event) {
    messageList.value.push(JSON.parse(event.data));
    console.log(JSON.parse(event.data));
    scrollToBottom();
  };

  socket.onclose = function () {
    console.log("连接已关闭");
    sockets[roomid] = null;  // 关闭连接时，清除对应房间的 WebSocket 连接
  };

  socket.onerror = function (error) {
    console.error("WebSocket 错误:", error);
  };
};


const choiceFile = () => {
  /*
  * 选择文件上传*
  * */
  const fileInput = document.querySelector("#fileInput")
  if (fileInput) {
    fileInput.click(); // 手动触发文件选择器
  }
  console.log(fileInput)
  console.log(456)
}
const handleFileChange = function (event) {
  const file = event.target.files[0]; // 获取选中的文件
  if (file) {
    console.log("选择的文件:", file);
    // 在这里处理文件，例如上传或显示预览
  }
}
const scrollToBottom = () => {
  nextTick(() => {
    const chatHistoryElement = document.querySelector('.list');
    if (chatHistoryElement) {
      chatHistoryElement.scrollTop = chatHistoryElement.scrollHeight;
    }
  });
};
const queryRoom = async () => {
  const res = await QueryRoom(parseInt(localStorage.getItem("userID")))
  console.log(res["room"])
  groupList.value = res["room"]
  console.log(groupList.value)
}
watch(
    messageStore.messages,
    async (newV, oldV) => {
      if (groupList.value === null){
         await queryRoom()
      }else {
        groupList.value.push(newV[0].info.room[0])
      }
    }
);

onMounted(() => {
  consentFriend()
  getisspeakuserinfo()
  avatar.value = localStorage.getItem("avatar")
  userID.value = parseInt(localStorage.getItem("userID"))
  queryRoom(userID.value)

})

const imgG = (info) => {
  // 图片排列
  if (info.length === 2) {
    return {
      flex: "0 1"
    }
  } else if (info.length === 3) {
    return {
      flex:"0 1"
    }

  } else if (info.length === 4) {
    return {
      flex: "0 1"
    }
  } else if (info.length === 5) {
    return {
      flex: "0 1"
    }
  } else if (info.length === 6) {
    return {
      flex: "0 1"
    }
  } else if (info.length === 7) {
    return {
      flex: "0 1"
    }
  } else if (info.length === 8) {
    return {
      flex: "0 1"
    }
  } else if (info.length === 9) {
    return {
      flex: "0 1"
    }
  } else {
    // 0?
  }
}
const img3 = (info)=>{
  if (info.length===2){
    return {
      width:"1.2vw",
      height:"2.4vh"
    }
  }else if (info.length===3){
    return {
      width:"1.2vw",
      height:"2.4vh"
    }
  }else {
    return {
      width:"0.8vw",
      height:"1.6vh"
    }
  }
}
// let groupSocket = {}
const group = (item)=>{
  GroupRoomID.value = item.roomID
  GroupOrUser.value = true
  oneinfo.value = item
  oneinfo.value.friend_username = item.roomName
  show.value = true;
  const userIDUUID = localStorage.getItem("uuid");
  // console.log(oneinfo.value)
  // const userId = localStorage.getItem("userID"); // 用户 ID
  const isGroup = 1; // 是否为群聊


  // 建立新的 WebSocket 连接
  const wsUrl = `ws://127.0.0.1:8080/ws?&IsGroup=${isGroup}&roomid=${item.roomID}&userIDUUID=${userIDUUID}`;
  groupSocket = new WebSocket(wsUrl);

  groupSocket.onopen = async function () {
    console.log(groupSocket)
    const res = await GetGroupMessageList(item.roomID)
    console.log(res)
    messageList.value = await GetGroupMessageList(item.roomID)
    // console.log(messageList.value)
    scrollToBottom()
  };

  groupSocket.onmessage = function (event) {
    // messageList.value.push(JSON.parse(event.data))
    console.log(JSON.parse(event.data))
    scrollToBottom()
  };

  groupSocket.onclose = function () {
    console.log("连接已关闭");
    groupSocket = null; // 连接关闭时，清空 socket 变量
  };

  groupSocket.onerror = function (error) {
    console.error("WebSocket 错误:", error);
  };
}
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
    <template #group>
      <ul class="group">
        <li style="cursor: pointer;" v-for="(item, index) in groupList" :key="index" @click="group(item)">
          <div class="image-container">
            <div class="im">
              <div v-for="(img, imgIndex) in item.userInfo.slice(0,9)"  :key="imgIndex"  :style="imgG(item.userInfo)" class="image-item">
                <img :src="img.avatar" alt="image" :style="img3(item.userInfo)"/>
              </div>
            </div>
            <div class="s">{{ item.roomName }}</div>
          </div>
        </li>
      </ul>
    </template>
    <template #three >
      <div class="dialogueList">
        <h2>{{ oneinfo.friend_username }}</h2>
        <ul class="list">
          <li :class="[userID===item.send_user_id?'right':'left']" v-for="(item,index) in messageList" :key="item.ID">
            <div v-if="userID===item.send_user_id" class="d">
              <div style="background-color: #95ec69" class="c">{{ item.content }}
                <svg class="icon sj" aria-hidden="true" font-size="0.8vw" color="#95ec69"
                     style="transform: rotate(90deg);margin-top: 1vh;position: absolute;margin-left: 0.3vw">
                  <use xlink:href="#icon-sanjiaoxing1"></use>
                </svg>
              </div>

<!--              <img :src="userID===item.send_user_id?avatar:oneinfo.friend_avatar" alt=""-->
<!--                   style="margin-left: 1vw;margin-right: 1vw">-->
              <img :src="item.Avatar" alt=""
                   style="margin-left: 1vw;margin-right: 1vw">
            </div>
            <div v-else class="d">
<!--              <img :src="userID===item.send_user_id?avatar:oneinfo.friend_avatar" alt=""-->
<!--                   style="margin-left: 0.5vw;margin-right: 1vw">-->
              <img :src="item.Avatar" alt=""
                   style="margin-left: 1vw;margin-right: 1vw">
              <div style="background-color: #ffffff" class="c">{{ item.content }}
                <div class="x">

                </div>
              </div>
            </div>
          </li>
        </ul>
        <div class="di">
          <V3Emoji :recent="true" class="xl" @click-emoji="clickEmoji"/>
          <input
              type="file"
              id="fileInput"
              style="display: none;"
              @change="handleFileChange"
          />
          <svg class="icon wj" aria-hidden="true" style="cursor: pointer" font-size="0.8vw" @click="choiceFile">
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

.group {
  padding: 0;
  list-style: none;
  width: 100%;
  height: 93vh;
  background-color: #eae7e7;
  overflow-y: auto;
  li {
    padding: 0.3vw;
    border-bottom: 1px solid black;
    .image-container {
      width: 100%;
      display: flex;
      background-color: #eae7e7;
      height: 6.5vh;
      transition: all 0.5s;
      .im {
        display: flex;
        flex-wrap: wrap; /* 允许换行 */
        width: 3.5vw;
        padding: 0.5vh;

        .image-item {
          display: flex;
          flex-wrap: wrap;
          justify-content: space-around; // 水平均匀分布
          align-items: center; // 垂直居中
          gap: 0.5vw; // 图片之间的间距
          padding: 0.1vh;
          img {
            object-fit: fill /* 保持长宽比，并填充区域 */
          }
        }

      }

      .s {
        width: 70%;
        color: black;
        text-align: right;
        padding-right: 3vh;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }
    }
    .image-container:hover{
      background-color: #72767b;
      transition: all 0.5s;
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
    height: 4.6vh;
    line-height: 4.6vh;
    color: black;
    border-bottom: 1px solid black;
  }

  .list {
    height: 60vh;
    padding: 0;
    margin: 0;
    overflow-y: auto;

    .left {
      list-style: none;
      text-align: left;
      height: 5vh;
      display: flex;
      align-items: center;
    }

    .right {
      list-style: none;
      text-align: right;
      height: 5vh;
      display: flex;
      justify-content: flex-end;
      align-items: center;

    }

    img {
      width: 2vw;
      height: 4vh;
    }

    .d {
      margin-left: 0.5vw;
      color: black;
      display: flex;
      position: relative;

      .c {
        padding-left: 1vh;
        padding-right: 1vh;
        height: 4vh;
        line-height: 4vh;
        border-radius: 0.3vw;
        position: relative;
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