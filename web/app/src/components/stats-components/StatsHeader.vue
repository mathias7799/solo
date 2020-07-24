<template>
  <div class="header-stats" id="main-header-stats">
    <div class="stat">
      <div class="sec">Workers Online/Offline</div>
      <div class="main high-letter-spacing" id="online-offline-workers">
        <mark>{{ workersOnline }}</mark>
        <span>/{{ workersOffline }}</span>
      </div>
    </div>
    <div class="stat">
      <div class="sec">Coinbase Balance</div>
      <div class="main">
        <div class="flexcol">
          <div>
            <span class="balance">{{ balance }}</span>
            <div class="lower">{{ ticker }}</div>
          </div>
          <div class="gray subbalance">{{ Math.round(price * balance * 100) / 100 }} USD</div>
        </div>
      </div>
    </div>
    <div class="stat">
      <div class="sec">Efficiency</div>
      <div class="main">
        <mark id="valid_shares_percentage_big">{{ efficiency }}</mark>%
      </div>
    </div>
  </div>
</template>

<script>
import $ from "jquery";
import { getCurrencyDetails, getCurrencyPrice } from "@/utils/currency.js";

export default {
  name: "StatsHeader",
  data() {
    return {
      workersOnline: 0,
      workersOffline: 1,
      balance: 0,
      efficiency: 0,
      ticker: "ETH",
      price: 0,
    };
  },
  created() {
    const updateData = (data) => {
      this.balance =
        Math.round(data.coinbaseBalance / Math.pow(10, 12)) / Math.pow(10, 6);
      if (window.innerWidth < 500) {
        this.balance =
          Math.round(this.balance * Math.pow(10, 3)) / Math.pow(10, 3);
      }

      this.onlineWorkers = data.onlineWorkers;
      this.offlineWorkers = data.workersOffline;
      this.efficiency = Math.round(data.efficiency * 10) / 10;

      var currencyDetails = getCurrencyDetails(data.chainId);
      this.ticker = currencyDetails.ticker;
      if (this.ticker == undefined) {
        this.ticker = "---";
      }

      getCurrencyPrice(
        currencyDetails.coingeckoId,
        (data) => {
          this.price = data[currencyDetails.coingeckoId]["usd"];
        },
        (data) => {
          if (currencyDetails.coingeckoId != undefined) {
            alert("Unable to get currency price: " + data.responseText);
          } else {
            this.price = NaN;
          }
        }
      );
    };

    $.get("http://localhost:8000/api/v1/headerStats", {}, function (data) {
      updateData(data.result);
    }).fail(function (data) {
      alert("Unable to fetch coinbase balance: " + data.responseJSON.error);
    });
  },
};
</script>

<style lang="scss" scoped>
@import "@/style/_colors.scss";
.header-stats {
  display: flex;
  justify-content: space-between;

  .stat {
    position: relative;
    display: flex;
    flex-direction: column;
    padding: 20px;
    background-color: white;
    width: 100%;
    border-radius: 10px;
    box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.1);
    position: relative;
    margin-left: 20px;
    display: flex;
    justify-content: center;

    .main {
      font-size: 40px;
      display: flex;
      font-weight: 600;
    }

    .main .lower {
      font-size: 15px;
      margin-left: 5px;
    }
    .stat .sec {
      padding-bottom: 3px;
      font-size: 15px;
      font-weight: 600;
    }
    .high-letter-spacing {
      letter-spacing: 3px;
    }
  }

  .stat:first-child {
    margin-left: 0px;
  }
  .stat > * {
    justify-content: flex-start;
  }
}

@media (max-width: 720px) {
  .header-stats .stat {
    margin-left: 5px;
  }

  .header-stats#main-header-stats .stat:last-child {
    display: none;
    justify-content: space-between;
  }
}

@media (max-width: 500px) {
  .stat {
    font-size: 7vw;
    padding: 6px;
  }

  .stat .title {
    font-size: 3vw;
  }

  .header-stats {
    justify-content: space-between;

    .stat {
      margin-top: 10px;
      width: 85%;

      .sec {
        padding-bottom: 1px;
        font-size: 3.5vw;
      }

      .main {
        font-size: 10vw;

        .lower {
          font-size: 3vw;
          margin-left: 2px;
        }
      }
    }
    .stat:first-child {
      margin-right: 10px;
    }
  }
}

.flexcol {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  width: auto;
}

.flexcol > div {
  display: flex;
  width: 100%;
}

.flexcol .subbalance {
  font-size: 15px;
}

@media (max-width: 980px) {
  .stat {
    font-size: 15px;
  }
}

#worker-search {
  height: 30px;
  border: none;
  border-bottom: 1px solid black;
  border-radius: 0px;
  font-size: 15px;
  width: 140px;
}

#worker-search:hover {
  border-bottom: 2px solid black;
}

#worker-search:focus {
  border-bottom: 2px solid #0069ff;
}

.account-help {
  height: 100%;
  cursor: pointer;
}

.account-help img {
  height: 20px;
  cursor: pointer;
}

.account-copy {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  width: 30px;
  height: 30px;
  margin-left: 10px;
  background-color: #efeff1;
  border-radius: 15px;
  border: 1px solid #00000000;
  transition: 0.1s;
  cursor: pointer;
}
</style>