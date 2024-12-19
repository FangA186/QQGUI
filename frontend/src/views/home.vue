<script setup>
import Index from "../components/Index.vue";
import V3Emoji from 'vue3-emoji'
import 'vue3-emoji/dist/style.css'

import {nextTick, onMounted, ref, watch} from "vue";
import {
  ConsentFriend,
  GetApplicationRecord,
  GetConsentList,
  GetGroupMessageList,
  GetIsSpeakUserInfo,
  QueryRoom,
  Receive
} from "../../wailsjs/go/api/MyApp.js"
import {useMessageStore} from '../store/messages.js';
const messageStore = useMessageStore();
const consentList = ref([])
const dialogVisible = ref(false)
const ApplicationRecord = ref([])
const twoli = ref(true)
const speakList = ref([])
const oneinfo = ref([])
const show = ref(false)
const content = ref("")
let groupSocket
let groupSockets = {}
let sockets = {};  // 用来存储不同房间或对话的 WebSocket 连接
const messageList = ref([])
const avatar = ref('')
const userID = ref('')
const left = ref('left')
const right = ref('right')
const groupList = ref([])
const GroupOrUser = ref(false)
const GroupRoomID = ref("")
const UserRoomID = ref('')
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
const sendMessage = async (event) => {
  event.preventDefault()
  if (GroupOrUser.value === false) {
    await Receive(localStorage.getItem("userID"),
        JSON.stringify(oneinfo.value.friend_id),
        oneinfo.value.applicant_avatar,
        oneinfo.value.applicant_username,
    )
    sockets[UserRoomID.value].send(JSON.stringify({
      two: 2,
      Avatar: localStorage.getItem("avatar"),
      content: content.value,
      send_user_id: parseInt(localStorage.getItem("userID")),
      receiver_user_id: oneinfo.value.friend_id,
      file_url: "",
      room_id: GenerateRoomID(localStorage.getItem("uuid"), oneinfo.value.friend_uuid)
    }))
  } else {
    groupSockets[GroupRoomID.value].send(JSON.stringify({
      two: 2,
      Avatar: localStorage.getItem("avatar"),
      Username: oneinfo.value.applicant_username,
      content: content.value,
      send_user_id: parseInt(localStorage.getItem("userID")),
      receiver_user_id: 0,
      room_id: GroupRoomID.value
    }))
  }

  content.value = ""
}
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
      if (groupList.value === null) {
        await queryRoom()
      } else {
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
      flex: "0 1"
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
const img3 = (info) => {
  if (info.length === 2) {
    return {
      width: "24px",
      height: "24px"
    }
  } else if (info.length === 3) {
    return {
      width: "24px",
      height: "24px"
    }
  } else {
    return {
      width: "14px",
      height: "14px"
    }
  }
}

const dialogue = async (item) => {
  GroupOrUser.value = false;
  oneinfo.value = item;
  show.value = true;
  console.log(oneinfo.value);

  const userIDUUID = localStorage.getItem("uuid");
  const friendIDUUID = item.friend_uuid;
  const isGroup = 0; // 是否为群聊
  const roomid = GenerateRoomID(userIDUUID, friendIDUUID);
  UserRoomID.value = roomid
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

const group = async (item) => {
  console.log(item)
  GroupRoomID.value = item.roomID
  GroupOrUser.value = true
  oneinfo.value = item
  oneinfo.value.friend_username = item.roomName
  show.value = true;
  const userIDUUID = localStorage.getItem("uuid");
  const isGroup = 1; // 是否为群聊
  // 建立新的 WebSocket 连接
  // 检查该房间是否已经有 WebSocket 连接
  if (groupSockets[item.roomID] && groupSockets[item.roomID].readyState === WebSocket.OPEN) {
    // 如果该房间的 WebSocket 已经连接，则无需重新连接
    messageList.value = await GetGroupMessageList(item.roomID);
    console.log(messageList.value)
    console.log("该房间的 WebSocket 已经连接，无需重新连接");
    return;
  }
  const wsUrl = `ws://127.0.0.1:8080/ws?&IsGroup=${isGroup}&roomid=${item.roomID}&userIDUUID=${userIDUUID}`;
  groupSocket = new WebSocket(wsUrl);
  // 将新连接保存在 sockets 对象中，按 roomid 管理
  groupSockets[item.roomID] = groupSocket;
  groupSocket.onopen = async function () {
    messageList.value = await GetGroupMessageList(item.roomID)
    scrollToBottom()
  };

  groupSocket.onmessage = function (event) {
    messageList.value.push(JSON.parse(event.data))
    console.log(JSON.parse(event.data))
    scrollToBottom()
  };

  groupSocket.onclose = function () {
    console.log("连接已关闭");
    groupSockets[item.roomID] = null // 连接关闭时，清空 socket 变量
  };

  groupSocket.onerror = function (error) {
    console.error("WebSocket 错误:", error);
  };
}
// import  DrawerProps  from 'element-plus'
// const direction = ref<DrawerProps['direction']>('rtl')
const detailShow = ref(false)
const detail = () => {
  console.log("点击了?", detailShow.value)
  detailShow.value = !detailShow.value

}
const selectedUser = ref(null); // 当前选中的用户信息
const userInfo = (item) => {
  // 设置选中的用户信息
  selectedUser.value = item;
}
const navtoSend = ()=>{
  console.log(selectedUser.value)
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
                <div v-if="item.IsConsent===0">
                  <img :src="item.applicant_avatar" alt="">
                  <span style="text-indent: 0.5vw">{{ item.applicant_username }}
                  <button style="width: 90px" @click="consent(item)">
                    同意
                  </button>
                    <button style="position: absolute;right: 0;top: 50%;left: 50%;transform: translateY(-50%);width: 90px"
                            @click="refuse(item)">
                    拒绝
                  </button>
                  </span>

                </div>
                <div v-else-if="item.IsConsent===1">
                  <img :src="item.applicant_avatar" alt="">
                  <span style="text-indent: 0.5vw;display: block;width: 100px;">{{ item.applicant_username }}</span>
                  <button style="width: 90px" disabled>已同意</button>
                </div>
                <div v-else>
                  <img :src="item.applicant_avatar" alt="">
                  <span style="text-indent: 0.5vw;display: block;width: 100px">{{ item.applicant_username }}></span>
                  <button style="width: 90px" disabled>已拒绝</button>
                </div>
              </li>
              <li v-for="(item,index) in ApplicationRecord" style="display: flex" v-else>
                <img :src="item.friend_avatar" alt="">
                <span style="text-indent: 0.5vw;display: block;width: 100px">{{ item.friend_username }}</span>
                <div v-if="item.IsConsent===0">
                  <button style="width: 5vw;" disabled>
                    已发送验证
                  </button>
                </div>
                <div v-else-if="item.IsConsent===1">
                  <button style="width: 90px" disabled>已同意</button>
                </div>
                <div v-else>
                  <button style="width: 90px" disabled>已拒绝</button>
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
              <div v-for="(img, imgIndex) in item.userInfo.slice(0,9)" :key="imgIndex" :style="imgG(item.userInfo)"
                   class="image-item">
                <img :src="img.avatar" alt="image" :style="img3(item.userInfo)"/>
              </div>
            </div>
            <div class="s">{{ item.roomName }}</div>
          </div>
        </li>
      </ul>
    </template>
    <template #three>
      <div class="dialogueList">
        <h2>{{ oneinfo.friend_username }}
          <el-button type="info"  @click="detail" style="position: absolute;right: 15px;top: 12px">
            <svg class="icon wj" aria-hidden="true" style="cursor: pointer;" font-size="1.5vw">
              <use xlink:href="#icon-xiangzuosanjiaoxing"></use>
            </svg>
          </el-button>

        </h2>

        <ul class="list">
          <li :class="[userID===item.send_user_id?'right':'left']" v-for="(item,index) in messageList" :key="item.ID" style="padding: 0.5vw">
            <div v-if="userID===item.send_user_id" class="d">
              <div class="c">
                <span style="text-align: left;">{{ item.content }}</span>
              </div>
              <img :src="item.Avatar" alt=""
                   style="margin-left: 0.5vw">
            </div>
            <div v-else class="d">
              <img :src="item.Avatar" alt=""
                   style="margin-right: 0.5vw">
              <div class="c1">
                <span>
                  {{ item.content }}
                </span>
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
          <textarea id="story" name="story" rows="5" cols="33" v-model="content" maxlength="1000"
                    @keydown.enter="sendMessage($event)">


          </textarea>
          <button @click="sendMessage">发送(enter)</button>
        </div>
      </div>
      <el-drawer v-model="detailShow">
        <template #default>
          <div class="groupInfo">
            <div class="info" v-for="(item,index) in oneinfo.userInfo" @click="userInfo(item,$event)">
              <img :src="item.avatar" :alt="item.userName">
              <span style="color: black">
                    {{ item.userName }}
                  </span>
            </div>
            <div v-if="selectedUser" class="user-info-popup">
              <div class="showInfo">
                <img :src="selectedUser.avatar" alt="">
                <p><strong>{{ selectedUser.userName }}</strong></p>
              </div>
              <div class="other">
                <p>备注</p>
                <p>共同群聊</p>
                <p>个性签名</p>
              </div>
              <div class="func">
                <div class="sendMessage" @click="navtoSend">
                  <svg class="icon wj" aria-hidden="true" style="cursor: pointer" font-size="1.5vw" @click="choiceFile">
                    <use xlink:href="#icon-fasong2"></use>
                  </svg>
                  <p><strong>发消息</strong></p>
                </div>
                <div class="voice">
                  <svg class="icon wj" aria-hidden="true" style="cursor: pointer" font-size="1.5vw" @click="choiceFile">
                    <use xlink:href="#icon-dianhua"></use>
                  </svg>
                  <p><strong>语音聊天</strong></p>
                </div>
                <div class="video">
                  <svg class="icon wj" aria-hidden="true"  style="cursor: pointer" font-size="1.5vw" @click="choiceFile">
                    <use xlink:href="#icon-shipin1"></use>
                  </svg>
                  <p><strong>视频聊天</strong></p>
                </div>
              </div>
            </div>
          </div>
        </template>
        <template #footer>
          <div style="flex: auto">
            <el-button @click="detailShow = false">cancel</el-button>
            <el-button type="primary" @click="detailShow= false">confirm</el-button>
          </div>
        </template>
      </el-drawer>
    </template>
  </Index>
</template>

<style scoped lang="less">

.searchList {
  width: 100%;
  padding: 0;
  margin: 0;
  //height: 200px;
  li {
    //width: 120px;
    text-align: left;
    list-style: none;
    padding: 0.5vh;
    height: 40px;
    line-height:40px;
    color: black;
    display: flex;
    align-items: center; // 垂直居中
    img {
      width: 30px;   // 固定宽度
      height: 30px;  // 固定高度
      border-radius: 0.5vh;
    }
    span {
      color: black;
      text-align: left;
    }
  }
}
.searchList1{
  width: 100%;
  padding: 0;
  margin: 0;
  li {
  //width: 120px;
    text-align: left;
    list-style: none;
    padding: 0.5vh;
    color: black;
    div{
      display: flex;
      button{
        border: none;
        background-color: #72767b;
        border-radius: 5px;
        color: black;
      }
    }
    img {
      width: 30px;   // 固定宽度
      height: 30px;  // 固定高度
      border-radius: 0.5vh;
    }
    span {
      color: black;
      text-align: left;
    }
  }
}

.group {
  padding: 0;
  list-style: none;
  width: 200px;
  margin: 0;
  //height: 92vh;
  overflow-y: auto;
  li {
    //border-bottom: 1px solid black;
  //background-color: #409eff; transition: all 0.3s ease;
    height: 55px;
    //background-color: #409eff;
    .image-container {
      width: 100%;
      display: flex;
      //background-color: #eae7e7;
      //height:40px;
      transition: all 0.5s;
      .im {
        display: flex;
        flex-wrap: wrap; /* 允许换行 */
        width: 55px;
        height: 55px;
        //background-color: #409eff;
        //padding: 3px;
        .image-item {
          display: flex;
          flex-wrap: wrap;
          justify-content: space-around;
          //align-items: center;
          gap: 5px;
          padding: 1px;

          img {
            border-radius: 5px;
            object-fit: fill /* 保持长宽比，并填充区域 */
          }
        }

      }

      .s {
        width: 70%;
        color: black;
        text-align: left;
        text-indent: 5px;
        padding-right:10px;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }
    }
  }

  li:hover {
    opacity: 0.8;
    //box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
  }
}

.searchList {
  //height: 95vh;
  //background-color: #d9d8d8;

  li:hover {
    opacity: 0.7;
    //background-color: #72767b;
  }
}

.dialogueList {
  height: 100vh;
  display: flex;
  flex-direction: column;
  position: relative;
  h2 {
    text-align: left;
    text-indent: 10px;
    height: 70px;
    line-height: 70px;
    color: black;
    margin: 0;
    position: relative;
    //border-bottom: 1px solid black;
  }

  .list {
    //height: 100%;
    padding: 0;
    margin: 0;
    overflow-y: auto;
    height: 60vh;
    //background-color: darkseagreen;
    .left {
      list-style: none;
      text-align: left;
      display: flex;
      align-items: center;
    //background-color: black; margin-bottom: 0.5vw;
    }

    .right {
      list-style: none;
      text-align: right;
      display: flex;
      justify-content: flex-end;
      align-items: center;
    //background-color: #72767b; margin-bottom: 0.5vw;
    }

    img {
      width: 2vw;
      height: 4vh;
    }

    .d {
    //margin-left: 0.5vw; color: black;
      display: flex;
    //position: relative;

      .c, .c1 {
        padding: 0 1vh; /* 水平方向的内边距 */
        border-radius: 0.3vw;
        display: inline-block; /* 使得c元素可以自适应内容 */

        span {
          word-wrap: break-word;
          word-break: break-all;
          height: auto; /* 让高度自动适应内容 */
          line-height: 4vh;
          display: inline-block;
          white-space: break-spaces; /* 允许换行 */
          max-width: 35vw; /* 设置最大宽度，根据需要调整 */
          font-family: 华文新魏, serif;
        }


      }

      .c {
        background-color: #95ec69;
        color: black;
        //opacity: 0.4;
        transition: all 0.3s ease; /* 添加过渡效果 */
      }

      .c1 {
        background-color: #ffffff;
        //opacity: 0.4;
        color: black;
        transition: all 0.3s ease;
      }

      .c1:hover {
        opacity: 0.8;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.2); /* 叠加模糊阴影效果 */
      }

      .c:hover {
        opacity: 0.8;
        box-shadow: 0 0 15px rgba(0, 0, 0, 0.2); /* 叠加模糊阴影效果 */
      }

    }

  }

  .di {
    position: relative;
    height: 30vh;
    width: 100%;
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
      position: absolute;
      letter-spacing: 1px;
      //margin-top: 30px;
      resize: none;
      font-size: 1vw;
      border: none;
      color: white;
      font-weight: 700;
      font-family: Verdana, serif;
      outline: none;
      caret-color: white;
      text-indent: 0.1vw;
      width: 98%;
      background-color: rgb(64, 65, 68);
      opacity: 0.4;
      border-radius: 0.4vw;
      left: 50%;
      top: 20%;
      height: 20vh;
      transform: translateX(-50%);
      box-shadow: 0 0 4px rgba(255, 255, 255, 0.5);
    }
    button{
      position: absolute;
      bottom: 3vw;
      right:1.5vw;
      border: none;
      width: 70px;
      height: 50px;
      color: green;
      font-weight: 700;
      border-radius: 0.5vw;
      cursor: pointer;
      background-color: white;
      transition: all 0.5s ;
    }

    button:hover {
      background-color:green ;
      color: white;
      transition: all 0.5s ;
    }


  }

}

.groupInfo {
  box-shadow: 0 0 15px rgba(0, 0, 0, 0.18); /* 叠加模糊阴影效果 */
  display: grid;
  grid-template-columns: 4vw 4vw 4vw;
  border-radius: 0.5vw;
  grid-template-rows: 4vw 4vw 4vw;
  width: 20vw;
  padding: 1vw;
  margin: 0 auto;
  .info {
    display: flex;
    flex-direction: column;
    align-items: center;
    min-height: 7vh;
    position: relative;
    span {
      width: 3vw;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }

    img {
      width: 3vw;
      min-height: 4vh;
      height: 5.5vh;
      border-radius: 0.5vw;
      cursor: pointer;
      transition: all 0.3s ease-in-out; /* 添加过渡效果 */
    }

    img:hover {
      box-shadow: 0 0 10px rgba(0, 0, 0, 0.7); /* 叠加模糊阴影效果 */
      transition: all 0.3s ease; /* 添加过渡效果 */
    }
  }
  .user-info-popup {
    position: absolute;
    background: whitesmoke;
    border: 1px solid #ddd;
    box-shadow: 0 0.4vw 0.8vw rgba(0, 0, 0, 0.1);
    padding: 1vw;
    border-radius: 0.5vw;
    right: 2vw;
    width: 50%;
    margin: 0;
    top: 53%;
    left: 50%;
    transform: translate(-50%, -50%);
    display: grid;
    //grid-template-columns: 5vw 5vw 5vw;
    grid-template-rows:4vw 4vw 5vw;
    height: 27vh;
    .showInfo{
      display: flex;
      align-items: center;
      img{
        width: 2vw;
        height: 4vh;
        border-radius: 0.5vw;
        margin-right: 0.5vw;
      }
      p{
        color: black;
      }
    }
    .other{
      display: flex;
      flex-direction: column;
      p{
        display: grid;
        grid-auto-rows: 2vw 2vw 2vw;
        text-align: left;
        color: #72767b;
      }
    }
    .func{
      display: grid;
      grid-template-columns: 5vw 5vw 5vw;
      width: 100%;
      height: 10vh;
      margin-top: 3vh;
      div{
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        cursor: initial;
        p{
          color: black;
          margin-top: 0.5vh;
        }
        transition: all 0.5s;
      }
      div:hover{
        border-radius: 0.5vw;
        box-shadow: 0 0.4vw 0.8vw rgba(0, 0, 0, 0.1);
        transition: all 0.5s;
      }
    }
  }

}
</style>