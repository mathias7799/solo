<template>
  <div class="worker-header">
    <div class="wrapper">
      <div class="title">WORKER</div>
      <div class="name">{{ workerName }}</div>
    </div>
    <div class="wrapper">
      <div class="btn" @click="reset()">Reset</div>
    </div>
  </div>
</template>

<script>
import $ from "jquery";

export default {
  name: "WorkerHeader",
  data() {
    return { workerName: "" };
  },
  mounted() {
    this.element = $(".worker-header");
  },
  methods: {
    open: function (workerName) {
      if (workerName == "") {
        return;
      }
      this.workerName = workerName;
      this.element.css("visibility", "visible");
      this.element.animate({ height: 100 }, 500);
      $("html, body").animate({ scrollTop: 0 }, 500);
    },
    reset: function () {
      this.element.animate({ height: 0 }, 500, () => {
        this.element.css("visibility", "hidden");
      });
      this.$emit("workerDeselected");
    },
  },
};
</script>

<style lang="scss" scoped>
@import "@/style/_buttons.scss";

.worker-header {
  display: flex;
  visibility: hidden;
  height: 0px;
}

.worker-header .wrapper {
  display: flex;
  justify-content: center;
  flex-direction: column;
  height: 100px;
}

.worker-header .wrapper:first-child {
  margin-right: 30px;
}

.worker-header .title {
  font-size: 20px;
  font-weight: 600;
}

.worker-header .name {
  font-size: 35px;
  font-weight: bold;
  color: #0069ff;
}
</style>