require("expose-loader?$!expose-loader?jQuery!jquery");
// require("bootstrap/dist/js/bootstrap.bundle.js");

$(function() {
  function init() {
    // 通过 data-xx 可以读取后台数据，注意是 单引号
    // data-server='<%= json(p) %>'
    const server = $("#idData");
    console.log(server.data("server"));
  }
  init();

  // 年份选择按钮
  $(".btn.year").on("click", function() {
    const allBtns = $(".btn.year");
    allBtns.addClass("btn-secondary");
    allBtns.removeClass("btn-primary");

    const self = $(this);
    self.addClass("btn-primary");
    self.removeClass("btn-secondary");

    const year = self.data("year");
    document.getElementById("selectedYear").innerText = year;
    document.getElementById("period").value = year;
  });

  // 任何一个值修改，都显示 已修改的标记
  $("input.indexValue").on("keydown", function() {
    $("#idChangedMark").css("display", "inline");
  });

  // 公司下拉框
  $("#companyOptions").on("change", function() {
    const company = this.value;
    document.getElementById("company").value = company;

    const period = document.getElementById("period").value;
    const matrixNum = document.getElementById("num").value;
    const target = `/matrix/openTypeA/?num=${matrixNum}&company=${company}&period=${period}`;
    // console.log(target);

    window.location = target;
  });
});
