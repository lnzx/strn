<template>
    <apexchart
        :options="options"
        :series="series"
        :height="height"
    ></apexchart>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
    data: Array
})

const data = ref(props.data)
const len = data.value.length
const height = len * 25

const options = {
    chart: {
        id: "isp",
        type: "bar",
    },
    plotOptions: {
        bar: {
            borderRadius: 4,
            horizontal: true,
        }
    },
    dataLabels: {
        enabled: true
    },
    xaxis: {
        categories: [],
    }
}

const series = [{
    name: 'nodes',
    data: []
}]

for (let i = 0; i < len; i++) {
    const e = data.value[i];
    options.xaxis.categories.push(e.name)
    series[0].data.push(e.count)
}

</script>