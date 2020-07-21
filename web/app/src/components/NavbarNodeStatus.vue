<template>
  <div class="node-status">
    <div
      class="status-indicator"
      :class="[isActive ? 'active' : 'inactive', isSyncing ? 'syncing' : '']"
    ></div>
    <span class="block-number">{{ blockNumber }}</span>
  </div>
</template>

<script>
export default {
  name: "NavbarNodeStatus",
  props: {
    isActive: Boolean,
    isSyncing: Boolean,
    blockNumber: {
      type: Number,
      default: 0
    }
  }
};
</script>

<style lang="scss" scoped>
@import "../style/_colors.scss";
@import "../style/_animations.scss";

.status-indicator {
  background-color: $gray;
  width: 7px;
  height: 7px;
  margin-right: 7px;
  border-radius: 7px;
}

.status-indicator.active {
  background-color: $green;
}

.status-indicator.inactive {
  background-color: $red;
}

.status-indicator.syncing {
  animation: blink 2s linear infinite;
}

.block-number {
  color: #444;
}
</style>