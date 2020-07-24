<template>
  <div class="table-wrapper mt20" id="workers-table">
    <table id="rigstats">
      <thead id="rigstats-thead">
        <tr>
          <th class="black-underline noselect" @click="sort('workerName')">
            Name
            <WorkerListSortIcon :sortValue="sortKeys.workerName" />
          </th>
          <th class="black-underline noselect" @click="sort('reportedHashrate')">
            Reported
            <WorkerListSortIcon :sortValue="sortKeys.reportedHashrate" />
          </th>
          <th class="black-underline noselect" @click="sort('effectiveHashrate')">
            Effective
            <WorkerListSortIcon :sortValue="sortKeys.effectiveHashrate" />
          </th>
          <th class="black-underline noselect" @click="sort('validShares')">
            Valid
            <WorkerListSortIcon :sortValue="sortKeys.validShares" />
          </th>
          <th class="black-underline noselect" @click="sort('staleShares')">
            Stale
            <WorkerListSortIcon :sortValue="sortKeys.staleShares" />
          </th>
          <th class="black-underline noselect" @click="sort('invalidShares')">
            Invalid
            <WorkerListSortIcon :sortValue="sortKeys.invalidShares" />
          </th>
          <th class="black-underline noselect" @click="sort('lastSeen')">
            Last Seen
            <WorkerListSortIcon :sortValue="sortKeys.lastSeen" />
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
import WorkerListSortIcon from "./WorkerListSortIcon.vue";

export default {
  name: "WorkerList",
  components: {
    WorkerListItem,
    WorkerListSortIcon,
  },
  data() {
    return {
      ascending: false,
      sortedKey: "",
      hasInitialized: false,
      sortKeys: {
        // Descending: -1, No sort: 0, Ascending: 1
        workerName: 0,
        reportedHashrate: 0,
        effectiveHashrate: 0,
        validShares: 0,
        staleShares: 0,
        invalidShares: 0,
        lastSeen: 0,
      },
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
  methods: {
    sort: function (key) {
      var aSort = -1;
      var bSort = 1;
      if (key != this.sortedKey) {
        this.sortKeys[this.sortedKey] = 0;
        this.sortedKey = key;
        if (this.hasInitialized) {
          this.ascending = true;
        } else {
          this.ascending = false;
          this.hasInitialized = true;
        }
      }
      if (this.ascending) {
        aSort = 1;
        bSort = -1;
      }
      this.ascending = !this.ascending;

      this.sortKeys[key] = bSort;
      this.workers.sort(function (a, b) {
        // Compare the 2 dates
        if (a[key] < b[key]) return aSort;
        if (a[key] > b[key]) return bSort;
        return 0;
      });
    },
  },
  beforeMount() {
    this.sort("workerName");
  },
};
</script>

<style lang="scss" scoped>
@import "../style/_utils.scss";
@import "../style/_tables.scss";
</style>