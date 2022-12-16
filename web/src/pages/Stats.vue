<template>
<main class="bg-color" id="main">
  <div class="table-container box px-0 py-2">
    <div class="px-1" style="padding-bottom: 0.5rem;">
      <input class="input is-small" type="text" v-model="ids" @keydown.enter="filter" placeholder="ID, ID, Press Enter to search" style="max-width: 97.5%">

      <div class="buttons are-small is-pulled-right">
        <button class="button is-white" title="refresh" @click="refresh">
            <span class="icon is-small i-color">
              <img src="https://img.icons8.com/material-outlined/24/null/update-left-rotation.png"/>
            </span>
        </button>
      </div>
    </div>
    <table class="table is-fullwidth is-hoverable is-striped">
        <thead style="background: #f8f8fb; font-size:14px;">
            <tr>
                <th>ID ({{ count }})</th>
                <th title="Title and version">Type</th>
                <th>IP & ISP</th>
                <th>Location</th>
                <th>Disk Usage</th>
                <th>Memory Usage</th>
                <th>CPU</th>
                <th title="Time to first byte stats">TTFB</th>
                <th>Cache rate</th>
                <th>Error rate</th>
                <th>Upload</th>
                <th>Download</th>
                <th>Last Registration</th>
                <th>Created</th>
                <th title="DNS weight from -1 to 100">Weight</th>
            </tr>
        </thead>
        <tbody style="color:#757981; font-size:14px;" v-html="trs"></tbody>
    </table>
  </div>

</main>  
</template>

<script>
import { ref, inject } from 'vue'

export default {
    setup() {
        const api = inject('api')
        const count = ref(0)
        const trs = ref('')
        const ids = ref('' || localStorage.getItem("node-ids"))
        var nodes = []
        var actives = 0

        const getData = () => {
          api.get('/api').then(res => {
            trs.value = ''
            const data = res.data
            nodes = data.nodes
            actives = data.count
            count.value = actives
            
            for (let i = 0,j=nodes.length; i < j; i++) {
              const node = nodes[i];
              trs.value += node.html;
            }
         })
        }

        const refresh = () => {
          getData()
        }

        const filter = () => {
          const val = ids.value
          if(!val){
            trs.value = ''
            count.value = actives

            for (let i = 0,j=nodes.length; i < j; i++) {
              const node = nodes[i];
              trs.value += node.html;
            }
            return
          }

          // other search
          if(val.indexOf('s:') == 0){
            let key = val.split(':')[1]
            
            trs.value = ''
            count.value = 0

            for (let i = 0,j=nodes.length; i < j; i++) {
              const node = nodes[i];
              if(node.html.indexOf(key) > -1){
                trs.value += node.html;
                count.value += 1
              }
            }

            return
          }

          // search ids
          if(val.trim().length > 7){
            trs.value = ''
            count.value = 0
            const ids = val.split(',')
            for (let i = 0,j=nodes.length; i < j; i++) {
              const node = nodes[i];
              if(ids.includes(node.id)){
                trs.value += node.html;
                count.value += 1
              }
            }
            localStorage.setItem("node-ids", val)
          }

        }

        getData()
        return {trs, count, refresh, ids, filter}
    },
}
</script>

<style scoped>

</style>