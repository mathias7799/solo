<template>
  <div class="table-wrapper mt20" id="workers-table">
    <table id="rigstats">
      <thead id="rigstats-thead">
        <tr>
          <th class="black-underline noselect">
            Name
            <img src="../assets/sort/sort-up.svg" alt="Sort Icon" />
          </th>
          <th class="black-underline noselect">
            Reported
            <img src="../assets/sort/sort.svg" alt="Sort Icon" />
          </th>
          <th class="black-underline noselect">
            Effective
            <img src="../assets/sort/sort.svg" alt="Sort Icon" />
          </th>
          <th class="black-underline noselect">
            Valid
            <img src="../assets/sort/sort.svg" alt="Sort Icon" />
          </th>
          <th class="black-underline noselect">
            Stale
            <img src="../assets/sort/sort.svg" alt="Sort Icon" />
          </th>
          <th class="black-underline noselect">
            Invalid
            <img src="../assets/sort/sort.svg" alt="Sort Icon" />
          </th>
          <th class="black-underline noselect">
            Last Seen
            <img src="../assets/sort/sort.svg" alt="Sort Icon" />
          </th>
        </tr>
      </thead>
      <tbody id="rigstats-tbody">
        <template v-for="worker in workers">
          <WorkerListItem
            :key="worker.workerName"
            :workerName="worker.workerName"
            :reportedHashrate="worker.reportedHashrate"
            :reportedHashrateSIChar="worker.reportedHashrateSIChar"
            :effectiveHashrate="worker.effectiveHashrate"
            :effectiveHashrateSIChar="worker.effectiveHashrateSIChar"
            :validShares="worker.validShares"
            :staleShares="worker.staleShares"
            :invalidShares="worker.invalidShares"
            :lastSeen="worker.lastSeen"
            :lastSeenHuman="worker.lastSeenHuman"
          />
        </template>
      </tbody>
    </table>
  </div>
</template>

<script>
import WorkerListItem from "./WorkerListItem.vue";

export default {
  name: "WorkerList",
  components: {
    WorkerListItem,
  },
  data() {
    return {
      workers: [
        {
          workerName: "rig1",
          reportedHashrate: 13.37,
          reportedHashrateSIChar: "M",
          effectiveHashrate: 23.37,
          effectiveSIChar: "M",
          validShares: 1,
          staleShares: 2,
          invalidShares: 3,
          lastSeen: 32,
          lastSeenHuman: "32 seconds ago",
        },
        {
          workerName: "rig2",
          reportedHashrate: 14.47,
          reportedHashrateSIChar: "M",
          effectiveHashrate: 21.17,
          effectiveSIChar: "M",
          validShares: 3,
          staleShares: 2,
          invalidShares: 1,
          lastSeen: 44,
          lastSeenHuman: "44 seconds ago",
        },
      ],
    };
  },
};
</script>

<style lang="scss" scoped>
@import "../style/_utils.scss";
@import "../style/_tables.scss";
</style>