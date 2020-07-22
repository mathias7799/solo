<template>
  <div class="header-stats" id="main-header-stats">
    <div class="stat">
      <div class="sec">Workers Online/Offline</div>
      <div class="main high-letter-spacing" id="online-offline-workers">
        <mark>{{ activeWorkers }}</mark>
        <span>/{{ inactiveWorkers }}</span>
      </div>
    </div>
    <div class="stat">
      <div class="sec">Coinbase Balance</div>
      <div class="main">
        <div class="flexcol">
          <div>
            <span class="balance">{{ balance }}</span>
            <div class="lower">ETH</div>
          </div>
          <div class="gray subbalance">{{ balanceUSD }} USD</div>
        </div>
      </div>
    </div>
    <div class="stat">
      <div class="sec">Efficiency</div>
      <div class="main">
        <mark id="valid_shares_percentage_big">{{ validSharesPercentage }}</mark>%
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "StatsHeader",
  data() {
    return {
      activeWorkers: 5,
      inactiveWorkers: 1,
      balance: 1.23456,
      validSharesPercentage: 99.8,
    };
  },
};
</script>

<style lang="scss" scoped>
.header-stats {
  display: flex;
  justify-content: space-between;
}

.header-stats .stat {
  position: relative;
  display: flex;
  flex-direction: column;
  padding: 20px;
  background-color: white;
  width: 100%;
  border-radius: 10px;
  box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.1);
}

.header-stats .stat {
  position: relative;
  margin-left: 20px;
  display: flex;
  justify-content: center;
}

#payout-bar {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 6px;
  border-radius: 0px 0px 10px 10px;
  background-color: #eee;
}

#payout-bar .bar {
  background-color: #0069ff;
  border-radius: 0px 0px 10px 10px;
  z-index: 1;
}

#payout-bar.green .bar {
  background-color: #15cd72;
}

.header-stats .stat:first-child {
  margin-left: 0px;
}

.header-stats .stat > * {
  justify-content: flex-start;
}

.header-stats .stat .main {
  font-size: 40px;
  display: flex;
  font-weight: 600;
}

.header-stats .stat .main .lower {
  font-size: 15px;
  margin-left: 5px;
}

.header-stats .stat .sec {
  padding-bottom: 3px;
  font-size: 15px;
  font-weight: 600;
}

mark.red {
  color: #e74c3c;
}

.header-stats .stat .high-letter-spacing {
  letter-spacing: 3px;
}

.account {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
  background-color: white;
  border-radius: 10px;
  box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.account-stats {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.account .account-stats .wallet {
  display: flex;
  align-items: center;
  font-size: 30px;
  padding-top: 1px;
}

.account .account-stats .wallet .address {
  cursor: pointer;
  transition: 0.05s;
  font-weight: 500;
  width: 100%;
}

.account .account-stats .wallet .address {
  color: black;
}

.account .account-stats .wallet .address:hover {
  color: #0069ff;
  text-decoration: none;
}

.account .settings {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.main-stats-wrapper {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}

.main-stats {
  border-radius: 10px;
  background-color: white;
  box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.1);
  padding: 10px;
  width: 100%;
}

.main-stats {
  margin-left: 10px;
}

.main-stats:first-child {
  margin-left: 0px;
}

.main-stats .stats-title {
  padding: 5px 0px 0px 5px;
  font-size: 25px;
}

.main-stats .wrapper {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.main-stats .wrapper .stat {
  display: flex;
  flex-direction: column;
  padding: 20px;
}

.main-stats .wrapper .stat .title {
  margin-bottom: 5px;
  font-size: 15px;
  color: #222;
}

.main-stats .wrapper .stat .content {
  font-size: 12px;
  color: #777;
  font-weight: bold;
  white-space: nowrap;
}

.main-stats .wrapper .stat .content mark.big {
  font-size: 29px;
  color: black;
}

.chart {
  border-radius: 10px;
  background-color: white;
  box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.1);
  margin-bottom: 10px;
  width: 100%;
  height: 250px;
}

.select-dashboard-category ul {
  list-style-type: none;
  display: flex;
  justify-content: center;
}

.select-dashboard-category ul li {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  padding: 15px;
  font-size: 20px;
  border-bottom: 3px solid #00000000;
  transition: color 0.05s;
  cursor: pointer;
  font-weight: 700;
  text-align: center;
  width: 100px;
}

@media (max-width: 500px) {
  .select-dashboard-category ul li {
    flex-grow: 1;
  }

  .select-dashboard-category ul {
    justify-content: flex-start;
  }

  .account .account-stats .wallet .address {
    font-weight: 600;
  }
}

.select-dashboard-category ul li:hover {
  color: #0069ff;
  transition: color 0.05s;
}

.select-dashboard-category ul li.selector-selected {
  border-bottom: 3px solid #0069ff;
}

#payouts-chart {
  margin-top: 20px;
}

#payouts {
  display: none;
}

.main-stats {
  margin-bottom: 20px;
}

#blocks {
  display: none;
}

.worker-name.selected {
  color: #0069ff;
}

#worker-header {
  display: flex;
  visibility: hidden;
  height: 0px;
}

#worker-header .wrapper {
  display: flex;
  justify-content: center;
  flex-direction: column;
  height: 100px;
}

#worker-header .wrapper:first-child {
  margin-right: 30px;
}

#worker-header .title {
  font-size: 20px;
  font-weight: 600;
}

#worker-header .name {
  display: none;
  font-size: 35px;
  font-weight: bold;
  color: #0069ff;
}

@media (max-width: 980px) {
  .account {
    flex-direction: column;
    align-items: center;
  }

  .account .account-stats .wallet {
    font-size: 6vw;
    text-align: center;
  }

  .account .settings {
    margin-top: 20px;
    display: block;
    width: 100%;
    text-align: center;
  }

  .header-stats.sec-header-stats {
    overflow-x: auto;
    padding: 10px;
    white-space: nowrap;
  }
}

@media (max-width: 720px) {
  .account .account-stats .wallet {
    font-weight: 600;
    margin-bottom: 5px;
  }

  .header-stats .stat {
    margin-left: 5px;
  }

  .header-stats#main-header-stats .stat:last-child {
    display: none;
    justify-content: space-between;
  }
}

@media (max-width: 500px) {
  .main-stats .stats-title {
    font-size: 4vw;
  }

  .main-stats .wrapper .stat .content {
    font-size: 5vw;
  }

  .main-stats .wrapper .stat .content mark.big {
    font-size: 7vw;
  }

  .main-stats .wrapper .stat {
    padding: 6px;
  }

  .main-stats .wrapper .stat .title {
    font-size: 3vw;
  }

  .header-stats {
    justify-content: space-between;
  }

  .header-stats .stat {
    margin-top: 10px;
    width: 85%;
  }

  .header-stats .stat:first-child {
    margin-right: 10px;
  }

  .header-stats .stat .main {
    font-size: 10vw;
  }

  .header-stats .stat .main .lower {
    font-size: 3vw;
    margin-left: 2px;
  }

  .header-stats .stat .sec {
    padding-bottom: 1px;
    font-size: 3.5vw;
  }

  .select-dashboard-category {
    overflow-x: scroll;
  }
}

.nostats {
  width: 700px;
  margin: 0px auto;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  padding: 45px;
}

.nostats img {
  max-width: 550px;
  margin-bottom: 20px;
}

.nostats span {
  color: #333;
  text-align: center;
}

.nostats span.primary {
  margin-top: 20px;
  font-size: 20px;
  font-weight: bold;
}

.nostats span.secondary {
  font-weight: 600;
  margin-top: 10px;
  font-size: 15px;
  font-weight: 100;
  color: #555;
}

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

  .nostats {
    width: 90%;
    padding: 15px;
  }

  .nostats img {
    width: 80%;
  }
}

.settings-modal-bg {
  z-index: 100;
  visibility: hidden;
  width: 100%;
  height: 100vh;
  position: fixed;
  top: 0;
  left: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: rgba(0, 0, 0, 0.75);
}

.settings-modal {
  position: relative;
  width: 300px;
  display: flex;
  justify-content: center;
  flex-direction: column;
  background-color: white;
  padding: 20px;
  border-radius: 15px;
  box-shadow: 0px 0px 60px 0px rgba(0, 0, 0, 0.5);
  border: 3px solid #0069ff;
}

.modal-window > div {
  margin-top: 15px;
}

.modal-input-group input {
  margin-top: 5px;
}

.modal-text-lower {
  margin-top: 10px;
}

.modal-input-group:first-child {
  margin-top: 0px;
}

.settings-modal-close {
  position: absolute;
  top: 10px;
  right: 10px;
  font-size: 40px;
  color: #888;
  cursor: pointer;
  width: 40px;
  height: 40px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.settings-modal-close:hover {
  color: black;
}

.modal-window *:first-child {
  margin-top: 0px;
}

.modal-input-group {
  display: flex;
  flex-direction: column;
}

.modal-input-group input {
  border: 1px solid #dfdfdf;
  border-radius: 5px;
  padding: 10px;
  font-size: 15px;
}

.modal-label-default {
  color: #aaa;
}

.modal-image-wrapper {
  width: 100%;
  display: flex;
  justify-content: center;
  margin: 20px 0px 20px 0px;
}

.modal-image {
  height: 150px;
}

.modal-text-h1 {
  display: flex;
  width: 100%;
  justify-content: center;
  font-size: 18px;
  /*font-weight: 100;*/
}

.modal-text-lower {
  display: flex;
  width: 100%;
  justify-content: center;
  font-size: 14px;
  color: #aaa;
}

.modal-text-lower:hover {
  cursor: pointer;
  text-decoration: underline;
}

.lds-dual-ring {
  display: inline-block;
}
.lds-dual-ring:after {
  content: " ";
  display: block;
  width: 100px;
  height: 100px;
  border-radius: 50%;
  border: 12px solid #0069ff;
  border-color: #0069ff transparent;
  animation: lds-dual-ring 1.2s linear infinite;
}
@keyframes lds-dual-ring {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

#settings-loading-window {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
}

.modal-settings-status-comment {
  display: flex;
  justify-content: center;
  flex-direction: column;
  font-size: 15px;
  text-align: center;
}

.modal-settings-status-comment span {
  text-align: center;
}

.modal-settings-status-comment a {
  cursor: pointer;
}

.modal-settings-status-comment a:hover {
  text-decoration: underline;
}

.account-info {
  display: flex;
  margin-top: 5px;
}

.account-info .info-card {
  font-weight: 600;
  margin-right: 25px;
  color: #aaa;
}

.info-card > mark {
  color: black;
  font-weight: bold;
}

.account-info .info-card:last-child {
  margin-right: 0px;
}

.btn-disabled {
  background-color: #bbb;
  cursor: not-allowed;
}

.btn-disabled:hover {
  background-color: #ccc;
}

@media (max-width: 980px) {
  .account-info {
    flex-direction: row;
    justify-content: center;
    flex-wrap: wrap;
    border-radius: 5px;
    background-color: white;
    border: 0.5px solid #eaeaea;
  }

  .account-info .info-card {
    color: #111;
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 10px;
    margin-right: 0px;
  }

  .account-info .info-card > *:last-child {
    margin-left: 5px;
  }
}

@media (max-width: 500px) {
  .account-info {
    border: 0px solid;
    border-top: 0.5px solid #eaeaea;
    border-radius: 0px;
  }

  .account-info .info-card {
    font-size: 15px;
    padding: 10px;
  }

  .account .settings {
    margin-top: 15px;
  }
}

@media (max-width: 980px) and (min-width: 500px) {
  .account-info .info-card {
    flex: 1 0 42%;
  }
}

@media (max-width: 500px) {
  .account-info {
    flex-direction: column;
    flex: 0 1 auto;
  }

  .account-info .info-card {
    border-bottom: 0.8px solid #efefef;
  }

  .account-info .info-card:last-child {
    border: none;
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

#calculator {
  display: none;
}

.workerstitle-with-filter,
.workers-search {
  display: flex;
  align-items: center;
}

@media (max-width: 980px) {
  .workerstitle-with-filter {
    justify-content: space-between;
  }

  .main-stats .wrapper .stat .content {
    font-size: 15px;
  }

  .wallet {
    justify-content: center;
  }

  .account .account-stats .wallet .address {
    width: auto;
  }

  .main-stats .wrapper {
    flex-direction: column;
  }
}

.workerstitle-with-filter img {
  margin-left: 25px;
  margin-right: 10px;
  height: 20px;
  width: 20px;
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

.account-copy svg {
  width: 15px;
  height: 15px;
}

.account-copy:hover {
  background-color: #0069ff;
  border: 1px solid #fff;
  transition: 0.1s;
  fill: white;
}

.copy-address-wrapper {
  display: flex;
  align-items: center;
  position: relative;
}

.balance .blank-element {
  height: 30px;
  width: 150px;
  margin: 5px 0px;
}

@media (max-width: 500px) {
  .account-stats {
    width: 100%;
  }

  .balance .blank-element {
    width: 100px;
  }
}

#stats {
  visibility: hidden;
}

.invalid-field {
  border: 1px solid #ed4f32 !important;
  background-color: #ed4f3210 !important;
}
</style>