<template>
  <div>
    <h1 class="text-big center m">Stats</h1>
    <StatsHeader />
    <WorkerHeader ref="workerHeader" v-on:workerDeselected="updateWorker('')" />
    <MainStats :selectedWorker="selectedWorker" />
    <Charts :selectedWorker="selectedWorker" />
    <WorkerList v-on:workerSelected="updateWorker($event)" />
  </div>
</template>

<script>
import Charts from "./Charts.vue";
import StatsHeader from "./StatsHeader.vue";
import MainStats from "./MainStats.vue";
import WorkerList from "./WorkerList.vue";
import WorkerHeader from "./WorkerHeader.vue";

export default {
  Name: "Stats",
  components: {
    Charts,
    StatsHeader,
    MainStats,
    WorkerList,
    WorkerHeader,
  },
  data() {
    return { selectedWorker: "" };
  },
  methods: {
    updateWorker(workerName) {
      this.selectedWorker = workerName;
      this.$refs.workerHeader.open(workerName);
    },
  },
};
</script>