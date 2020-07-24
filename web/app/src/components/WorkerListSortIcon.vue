<template>
  <img src="../assets/sort/sort-up.svg" alt="Sort Icon" v-if="sortValue == 1" />
  <img src="../assets/sort/sort-down.svg" alt="Sort Icon" v-else-if="sortValue == -1" />
  <img src="../assets/sort/sort.svg" alt="Sort Icon" v-else />
</template>

<script>
export default {
  name: "WorkerListSortIcon",
  props: {
    sortValue: Number,
  },
};
</script>