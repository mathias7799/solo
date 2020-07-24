<template>
  <div class="main-stats-wrapper">
    <div class="main-stats">
      <div class="stats-title">Hashrate</div>
      <div class="wrapper">
        <div class="stat">
          <div class="content" id="effective_hashrate">
            <mark class="big">{{ effective }}</mark>
            {{ siChar }}H/s
          </div>
          <div class="title">Effective</div>
        </div>
        <div class="stat">
          <div class="content" id="average_hashrate">
            <mark class="big">{{ average }}</mark>
            {{ siChar }}H/s
          </div>
          <div class="title">Average</div>
        </div>
        <div class="stat">
          <div class="content" id="reported_hashrate">
            <mark class="big">{{ reported }}</mark>
            {{ siChar }}H/s
          </div>
          <div class="title">Reported</div>
        </div>
      </div>
    </div>
    <div class="main-stats">
      <div class="stats-title">Shares</div>
      <div class="wrapper">
        <div class="stat">
          <div class="content" id="valid_shares">
            <mark class="big">{{ validShares }}</mark>
          </div>
          <div class="title">
            <span>Valid (</span>
            <span
              id="valid_shares_percentage"
            >{{ Math.round(validShares / totalShares * 100) / 100 }}</span>%)
          </div>
        </div>
        <div class="stat">
          <div class="content" id="stale_shares">
            <mark class="big">{{ staleShares }}</mark>
          </div>
          <div class="title">
            <span>Stale (</span>
            <span
              id="stale_shares_percentage"
            >{{ Math.round(staleShares / totalShares * 100) / 100 }}</span>%)
          </div>
        </div>
        <div class="stat">
          <div class="content" id="invalid_shares">
            <mark class="big">{{ invalidShares }}</mark>
          </div>
          <div class="title">
            <span>Invalid (</span>
            <span
              id="invalid_shares_percentage"
            >{{ Math.round(invalidShares / totalShares * 100) / 100 }}</span>%)
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import $ from "jquery";

export default {
  name: "MainStats",
  data() {
    return {
      // Hashrate
      siChar: "",
      effective: 0,
      reported: 0,
      average: 0,
      // Shares
      validShares: 0,
      staleShares: 0,
      invalidShares: 0,
      totalShares: 0,
    };
  },
  created() {
    const updateData = (data) => {
      this.siChar = data.si.char;
      var siDiv = data.si.div;
      this.effective = Math.round((data.hashrate.effective / siDiv) * 10) / 10;
      this.reported = Math.round((data.hashrate.reported / siDiv) * 10) / 10;
      this.average = Math.round((data.hashrate.average / siDiv) * 10) / 10;

      this.validShares = data.shares.valid;
      this.staleShares = data.shares.stale;
      this.invalidShares = data.shares.invalid;
      this.totalShares =
        this.validShares + this.staleShares + this.invalidShares;
    };
    $.get("http://localhost:8000/api/v1/stats", {}, function (data) {
      updateData(data.result);
    }).fail(function (data) {
      alert("Unable to fetch stats: " + data.responseJSON.error);
    });
  },
};
</script>

<style lang="scss" scoped>
.main-stats-wrapper {
  display: flex;
  justify-content: space-between;
  margin: 20px 0px;
}

.main-stats {
  border-radius: 10px;
  background-color: white;
  box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.1);
  padding: 10px;
  width: 100%;
  margin-left: 10px;
  .stats-title {
    padding: 5px 0px 0px 5px;
    font-size: 25px;
  }
  .wrapper {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
  }
  .stat {
    display: flex;
    flex-direction: column;
    padding: 20px;
    .title {
      margin-bottom: 5px;
      font-size: 15px;
      color: #222;
    }

    .content {
      font-size: 12px;
      color: #777;
      font-weight: bold;
      white-space: nowrap;

      mark.big {
        font-size: 29px;
        color: black;
      }
    }
  }
}

.main-stats:first-child {
  margin-left: 0px;
}

.select-dashboard-category ul li.selector-selected {
  border-bottom: 3px solid #0069ff;
}

@media (max-width: 500px) {
  .main-stats .stats-title {
    font-size: 4vw;
  }

  .main-stats .wrapper .stat {
    padding: 6px;

    .title {
      font-size: 3vw;
    }

    .content {
      font-size: 5vw;

      mark.big {
        font-size: 7vw;
      }
    }
  }
}

@media (max-width: 980px) {
  .main-stats .wrapper {
    flex-direction: column;
    .stat .content {
      font-size: 15px;
    }
  }
}
</style>