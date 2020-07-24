<template>
  <div>
    <div class="simple-wrapper">
      <highcharts class="chart" :options="hashrateChartOptions"></highcharts>
    </div>
    <div class="simple-wrapper mt20">
      <highcharts class="chart" :options="sharesChartOptions"></highcharts>
    </div>
  </div>
</template>

<script>
import $ from "jquery";
import { getSi } from "@/utils/format.js";

export default {
  props: {
    selectedWorker: String,
  },
  data() {
    var hashrateChartOptions = {
      colors: ["#0069ff", "#2c3e50"],
      chart: {
        height: 200,
        type: "spline",
      },
      title: {
        text: "Hashrate",
      },
      xAxis: {
        type: "datetime",
        dateTimeLabelFormats: {
          month: "%e. %b",
          year: "%b",
        },
        title: {
          text: "Date",
        },
      },
      yAxis: {
        title: {
          text: "Hashrates",
        },
        min: 0,
      },
      tooltip: {
        headerFormat: "<b>{series.name}</b><br>",
        pointFormat: ``,
      },
      time: {
        timezoneOffset: new Date().getTimezoneOffset(),
      },

      series: [
        {
          name: "Effective Hashrate",
          data: [],
        },
        {
          name: "Reported Hashrate",
          data: [],
        },
      ],
    };

    var sharesChartOptions = {
      colors: ["#1633ff", "#0069ff", "#031b4e"],
      chart: {
        type: "column",
        height: 200,
      },
      title: {
        text: "Shares",
      },
      yAxis: {
        title: {
          text: "Shares",
        },
        min: 0,
      },
      xAxis: {
        type: "datetime",
        dateTimeLabelFormats: {
          month: "%e. %b",
          year: "%b",
        },
        title: {
          text: "Date",
        },
      },
      tooltip: {
        headerFormat: "<b>{series.name}</b><br>",
        pointFormat: "{point.x:%e. %b %H:%M}: {point.y} Shares",
      },
      time: {
        timezoneOffset: new Date().getTimezoneOffset(),
      },

      series: [
        {
          name: "Valid shares",
          data: [],
        },
        {
          name: "Stale shares",
          data: [],
        },
        {
          name: "Invalid shares",
          data: [],
        },
      ],
    };
    return { hashrateChartOptions, sharesChartOptions };
  },
  watch: {
    // whenever question changes, this function will run
    selectedWorker: function (newWorkerName) {
      this.updateChart(newWorkerName);
    },
  },

  mounted() {
    this.updateChart("");
  },
  methods: {
    updateChart: function (workerName) {
      console.log('updating chart for "' + workerName + '"');
      const updateData = (
        effectiveHashrate,
        reportedHashrate,
        validShares,
        staleShares,
        invalidShares,
        siChar
      ) => {
        this.hashrateChartOptions.series[0].data = effectiveHashrate;
        this.hashrateChartOptions.series[1].data = reportedHashrate;
        this.sharesChartOptions.series[0].data = validShares;
        this.sharesChartOptions.series[1].data = staleShares;
        this.sharesChartOptions.series[2].data = invalidShares;
        this.hashrateChartOptions.tooltip.pointFormat =
          `{point.x:%e. %b %H:%M}: {point.y:.2f} ` + siChar + `H/s`;
      };
      console.log("params", { workerName });
      $.get("http://localhost:8000/api/v1/history", { workerName }, function (
        data
      ) {
        var avgEffectiveHashrate = [];

        data.result.forEach((item) => {
          avgEffectiveHashrate.push(item.effectiveHashrate);
        });

        // Calculating average
        avgEffectiveHashrate =
          avgEffectiveHashrate.reduce((a, b) => a + b, 0) -
          avgEffectiveHashrate.length;

        var si = getSi(avgEffectiveHashrate);

        var effectiveHashrateHistory = [];
        var reportedHashrateHistory = [];
        var validSharesHistory = [];
        var staleSharesHistory = [];
        var invalidSharesHistory = [];

        data.result.forEach((item) => {
          effectiveHashrateHistory.push([
            item.timestamp * 1000,
            item.effectiveHashrate / si[0],
          ]);
          reportedHashrateHistory.push([
            item.timestamp * 1000,
            item.reportedHashrate / si[0],
          ]);
          validSharesHistory.push([item.timestamp * 1000, item.validShares]);
          staleSharesHistory.push([item.timestamp * 1000, item.staleShares]);
          invalidSharesHistory.push([
            item.timestamp * 1000,
            item.invalidShares,
          ]);
        });

        updateData(
          effectiveHashrateHistory,
          reportedHashrateHistory,
          validSharesHistory,
          staleSharesHistory,
          invalidSharesHistory,
          si[1]
        );
      }).fail(function (data) {
        alert("Unable to fetch history: " + data.responseJSON.error);
      });
    },
  },
};
</script>

<style lang="scss" scoped>
@import "@/style/_wrappers.scss";
@import "@/style/_charts.scss";
</style>