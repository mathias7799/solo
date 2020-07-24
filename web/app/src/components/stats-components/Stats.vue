<template>
  <div>
    <h1 class="text-big center m">Statistics</h1>
    <StatsHeader />
    <MainStats />
    <Charts />
    <WorkerList />
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
};
</script>