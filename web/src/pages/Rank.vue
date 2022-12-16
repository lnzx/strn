<template>
  <div class="container">
    <div class="columns">
      <!-- region -->
      <div class="column">
        <span class="title is-6" style="margin-left: 50%;">Countries ({{ regions.length }})</span>
        <BarChart :data="regions" />
      </div>

      <!-- isp -->
      <div class="column">
        <span class="title is-6" style="margin-left: 55%;">ISPs ({{ isps.length }})</span>
        <BarChart :data="isps" />
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref,inject } from 'vue'
  import BarChart from '@comps/BarChart.vue'

  const api = inject('api')
  const isps = ref([])
  const regions = ref([])

  await api.get('/api').then(res => {
    const data = res.data
    isps.value = data.isps
    regions.value = data.regions
  })
</script>

