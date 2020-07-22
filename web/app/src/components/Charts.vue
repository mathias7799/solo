<template>
  <highcharts class="chart" :options="hashrateChartOptions"></highcharts>
</template>

<script>
export default {
  data() {
    var hashrateChartOptions = {
      colors: ["#0069ff", "#2c3e50", "#3498db"],
      chart: {
        type: "spline"
      },
      title: {
        text: "Hashrate"
      },
      xAxis: {
        type: "datetime",
        dateTimeLabelFormats: {
          month: "%e. %b",
          year: "%b"
        },
        title: {
          text: "Date"
        }
      },
      yAxis: {
        title: {
          text: "Hashrates"
        },
        min: 0
      },
      tooltip: {
        headerFormat: "<b>{series.name}</b><br>",
        pointFormat: `{point.x:%e. %b %H:%M}: {point.y:.2f} ${window.chart_hasrate_si}H/s`
      },
      time: {
        timezoneOffset: new Date().getTimezoneOffset()
      },

      series: [
        {
          name: "Effective Hashrate",
          id: "ehashrate",
          data: []
        },
        {
          name: "Reported Hashrate",
          data: []
        }
      ]
    };
    return { hashrateChartOptions };
  }
};
</script>