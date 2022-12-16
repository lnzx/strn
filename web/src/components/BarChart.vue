<template>
    <apexchart
        :options="options"
        :series="series"
        :height="height"
    ></apexchart>
</template>

<script setup>
const props = defineProps({
    data: Array
})

const data = props.data
const len = data.length
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
    const e = data[i];
    let name = e.name
    if(!e.name){
        name = '-'
    }
    options.xaxis.categories.push(name)
    series[0].data.push(e.count)
}

</script>