<template>
  <v-container class="fill-height">
    <v-responsive class="align-centerfill-height mx-auto">
      <div class="text-center">
        <v-btn @click="list">list</v-btn>
      </div>
      <v-row>
        <v-col>
          <v-form>
            <v-card>
              <v-row>
                <v-col no-gutters class="mr-0 pr-1 d-flex" density="compact" cols="10" sm="10" md="10">
                  <v-text-field v-model="refs.new_url" label="New URL" type="text"></v-text-field>
                </v-col>
                <v-col no-gutters cols="2" sm="2" md="2">
                  <v-btn @click="create">create</v-btn>
                </v-col>
              </v-row>
            </v-card>
          </v-form>
        </v-col>
      </v-row>
      <v-row v-for="(index) in refs.urls" :key="index">
        <v-col>
          <v-col><a :href="index.href">{{ index.url }}</a></v-col>
        </v-col>

      </v-row>

    </v-responsive>
  </v-container>
</template>

<script setup>
import { ref } from 'vue';

const refs = ref(
  {
    urls: [],
    new_url: "",
  })


function list() {
  console.log("list")
  const headers = {};
  headers["Authorization"] = "Bearer admin@secret";

  const options = {
    method: "GET",
    headers: headers,
    // credentials: "omit",
  };
  fetch("http://localhost:8080/admin/v1", options).
    then((response) => {
      return response.json()
        .then(data => {
          refs.value.urls = []
          data.forEach(element => {
            console.log("json", element)
            if (element.long_url != "") {
              refs.value.urls.push({ href: "http://localhost:8080/" + element.key, url: element.long_url })
            }
          });
        }).catch(err => {
          console.log("json err", err)
        })
    }).
    catch(err => {
      console.log(err)
    })
}

function create() {
  console.log("list")
  const headers = {};
  headers["Authorization"] = "Bearer admin@secret";

  const options = {
    method: "POST",
    headers: headers,
    body: JSON.stringify({ long_url: refs.value.new_url }),
  };
  fetch("http://localhost:8080/admin/v1", options).
    then((response) => {
      return response.json()
        .then(data => {
          console.log("json", data)
        }).catch(err => {
          console.log("json err", err)
        })
    }).
    catch(err => {
      console.log(err)
    })
}

</script>
