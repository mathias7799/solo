<template>
  <tr>
    <td :id="'worker-' + workerName">
      <!-- TODO: onclick="">-->
      <div class="space-between">
        <span class="worker-name black-underline">{{ workerName }}</span>
        <span class="rig-offline" v-if="lastSeenSecs > 600">Offline</span>
      </div>
    </td>
    <td class="bluegray">
      {{ reportedHashrate }}
      <span class="bluegray">{{ reportedHashrateSIChar }}H/s</span>
    </td>
    <td class="bluegray">
      {{ effectiveHashrate }}
      <span class="bluegray">{{ effectiveHashrateSIChar }}H/s</span>
    </td>
    <td>{{ validShares }}</td>
    <td>{{ staleShares }}</td>
    <td>{{ invalidShares }}</td>
    <td>{{ lastSeenHuman }}</td>
  </tr>
</template>

<script>
export default {
  props: {
    workerName: String,
    reportedHashrate: Number,
    reportedHashrateSIChar: String,
    effectiveHashrate: Number,
    effectiveHashrateSIChar: String,
    validShares: Number,
    staleShares: Number,
    invalidShares: Number,
    lastSeenSecs: Number,
    lastSeenHuman: String,
  },
};
</script>

<style lang="scss" scoped>
@import "../style/_tables.scss";
</style>