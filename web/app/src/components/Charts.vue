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
export default {
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
        pointFormat: `{point.x:%e. %b %H:%M}: {point.y:.2f} ${window.chart_hasrate_si}H/s`,
      },
      time: {
        timezoneOffset: new Date().getTimezoneOffset(),
      },

      series: [
        {
          name: "Effective Hashrate",
          id: "ehashrate",
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
};
</script>

<style lang="scss" scoped>
@import "../style/_wrappers.scss";
@import "../style/_charts.scss";
</style>