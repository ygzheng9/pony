require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");

import { el, mount, setChildren } from "redom";

$(function() {
  // 显示提示信息
  $('[data-toggle="popover"]').popover();

  function init() {
    // 当前年份
    const period = document.getElementById("period").value;
    const s = `.btn.year[data-year="${period}"]`;
    const e = $(s);
    e.addClass("btn-primary");
    e.removeClass("btn-secondary");

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

  // 动态创建 dom
  function createUI() {
    // 创建多个 dom，并且绑定事件
    var arr = ["first", "second", "third"];
    const b = arr.map(a => {
      const t = el("div.idDummy", a, { "data-tst": a });
      t.onclick = clickHandler(a);
      return t;
    });

    setChildren(document.getElementById("idTest"), b);
  }
  createUI();

  function clickHandler(name) {
    return function() {
      console.log("clicked", name);
    };
  }
});
