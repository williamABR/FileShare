<template>
  <h1>
    {{
      checkCookie("userId")
        ? "Bienvenido, Ya puedes compartir archivos"
        : "Usuario no registrado"
    }}
  </h1>
  <!--<h1>Hola, {{ getCookie("userId") }}</h1>-->
  <div class="container" id="app">
    <div class="subContainer">
      <form
        enctype="multipart/form-data"
        @submit.prevent="prueba($event)"
        name="formSend"
        id="formSend"
      >
        <h1>Compartir Archivo</h1>
        <div class="form-group">
          <label for="exampleFormControlFile1">Seleccionar un archivo:</label>
          <br />
          <input
            type="file"
            class="form-control-file"
            id="exampleFormControlFile1"
            name="myFile"
          />
          <div class="channels input-group">
            <div class="input-group-prepend mx-2">
              <label class="input-group-text" for="inputGroupSelect01"
                >Canal:
              </label>
            </div>
            <select
              class="form-select"
              aria-label="Default select example"
              name="channel"
              id="selectChannel"
            >
              <option value="channel1" selected>Channel 1</option>
              <option value="channel2">Channel 2</option>
              <option value="channel3">Channel 3</option>
            </select>
          </div>
          <input type="submit" value="Compartir Archivo" class="btnUpload" />
        </div>
      </form>
      <!--Subscription-->
      <form @submit.prevent="subscribir($event)" name="formSub" id="formSub">
        <h1>Subscribir a Canal</h1>
        <div class="form-group">
          <div class="channels input-group">
            <div class="input-group-prepend mx-2">
              <label class="input-group-text" for="inputGroupSelect01"
                >Canal:
              </label>
            </div>
            <select
              class="form-select"
              aria-label="Default select example"
              name="channel1"
              id="selectChannel1"
            >
              <option value="channel1" selected>Channel 1</option>
              <option value="channel2">Channel 2</option>
              <option value="channel3">Channel 3</option>
            </select>
          </div>
          <input type="submit" value="Subscribirse a Canal" class="btnUpload" />
        </div>
      </form>
    </div>
    <div class="subContainer">
      <div class="card cardFiles">
        <div class="card-body">
          <ul class="list-group list-group-flush" :key="actualizar">
            <li v-for="item in files" :key="item.id">
              <a :href="'http://localhost:8080/file/' + item.fileName">
                {{ item.fileName }}
              </a>
              <div>{{ item.canal }}</div>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
  <!--<button @click="prueba">Get posts</button>-->
</template>

<script>
import { uuid } from "vue-uuid"; // Import uuid
import { crono } from "vue-crono";
import axios from "axios";

//const form = document.querySelector("form");

export default {
  mixins: [crono],
  name: "HelloWorld",
  data() {
    return {
      posts: null,
      files: [],
      subs: [],
      actualizar: 0,
    };
  },
  props: {
    msg: String,
  },
  methods: {
    async subscribir() {
      const selectAux = document.getElementById("selectChannel1");
      const headers = { "Content-Type": "application/x-www-form-urlencoded" };
      var obj = {"idUser":this.getCookie("userId"),"channel":selectAux.value}
      this.posts = await axios
        .post("http://localhost:8080/subscription", obj, {headers})
        .then((response) => {
          for (var i = 0; i < response.length; i++){
          this.files.push = response[i]
          }
        })
        .catch((error) => {
          this.errorMessage = error.message;
          console.error("There was an error!", error);
        });
    },
    async prueba() {
      const formData = new FormData(document.getElementById("formSend"));
      formData.append("idUser", this.getCookie("userId"));
      const headers = { "Content-Type": "multipart/form-data" };
      this.posts = await axios
        .post("http://localhost:8080/upload", formData, { headers })
        .then((response) => {
          var e = document.getElementById("selectChannel");
          var value = e.value;

          var objectFile = {
            canal: value,
            fileName: response.data,
          };
          this.files.push(objectFile);
          this.actualizar += 1;
          console.log(this.files);
        })
        .catch((error) => {
          this.errorMessage = error.message;
          console.error("There was an error!", error);
        });
    },
    setCookie(cname, cvalue, exdays) {
      const d = new Date();
      d.setTime(d.getTime() + exdays * 24 * 60 * 60 * 1000);
      let expires = "expires=" + d.toUTCString();
      document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
    },
    getCookie(cname) {
      let name = cname + "=";
      let decodedCookie = decodeURIComponent(document.cookie);
      let ca = decodedCookie.split(";");
      for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) == " ") {
          c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
          return c.substring(name.length, c.length);
        }
      }
      return "";
    },

    checkCookie(cname) {
      let user = this.getCookie(cname);
      if (user != "") {
        return true;
      } else {
        //return false;
        this.setCookie("userId", uuid.v1(), 30);
        return true;
      }
    },
    async checkFiles() {
      var sizeFile = true
      if (this.files.length === 0) {
        sizeFile = true
      }else{
        sizeFile = false
      }
      await axios
        .post('http://localhost:8080/checkUpdate/'+this.getCookie("userId"), sizeFile)
        .then((response) => {
          console.log(response);
        })
        .catch((error) => {
          this.errorMessage = error.message;
          console.error("There was an error!", error);
        });
    },
  },
  cron: {
    time: 1000,
    method: "checkFiles",
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
.container {
  display: flex;
  flex-wrap: wrap;
}
.subContainer {
  flex: 1 1 200px;
  padding: 16px;
}
.cardFiles {
  height: 100%;
  background-color: ghostwhite;
  padding: 16px;
  border-radius: 15px;
  border: solid 1px;
  box-sizing: border-box;
}
.channels {
  display: flex;
  margin: 8px;
  justify-content: center;
  font-weight: bold;
}
.btnUpload {
  margin: auto;
  padding: 8px;
  width: 90%;
  background: #42b983;
  color: white;
  border-radius: 15px;
}
</style>
