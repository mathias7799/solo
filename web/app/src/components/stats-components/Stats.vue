<template>
  <div>
    <h1 class="text-big center m">Stats</h1>
    <StatsHeader />
    <MainStats />
    <Charts :selectedWorker="selectedWorker" />
    <WorkerList v-on:workerSelected="updateWorker($event)" />
  </div>
</template>

<script>
import Charts from "./Charts.vue";
import StatsHeader from "./StatsHeader.vue";
import MainStats from "./MainStats.vue";
import WorkerList from "./WorkerList.vue";

export default {
  Name: "Stats",
  components: {
    Charts,
    StatsHeader,
    MainStats,
    WorkerList,
  },
  data() {
    return { selectedWorker: "" };
  },
  methods: {
    updateWorker(workerName) {
      this.selectedWorker = workerName;
    },
  },
};
</script>