<template>
  <tr>
    <td :id="'worker-' + workerName">
      <!-- TODO: onclick="">-->
      <div class="space-between">
        <span class="worker-name black-underline">{{ workerName }}</span>
        <span class="rig-offline" v-if="lastSeen > 600">Offline</span>
      </div>
    </td>
    <td :class="{bluegray: lastSeen > 600}">
      {{ reportedHashrate }}
      <span class="bluegray">{{ reportedHashrateSIChar }}H/s</span>
    </td>
    <td :class="{bluegray: lastSeen > 600}">
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
import humanizeDuration from "humanize-duration";

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
    lastSeen: Number,
  },
  data() {
    return { lastSeenHuman: "now" };
  },
  mounted() {
    if (this.lastSeen < 1) {
      this.lastSeenHuman = "now";
    } else {
      this.lastSeenHuman =
        humanizeDuration(this.lastSeen * 1000, {
          round: true,
          largest: 1,
        }) + " ago";
    }
  },
};
</script>

<style lang="scss" scoped>
@import "@/style/_tables.scss";
.rig-offline {
  font-size: 10px;
  background-color: #ed4f32;
  color: white;
  margin-left: 10px;
  padding: 5px;
  border-radius: 5px;
}

.rig-offline:hover {
  text-decoration: none;
}

@media (max-width: 720px) {
  .rig-offline {
    font-size: 2vw;
  }
}
</style>