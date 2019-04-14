import G2 from "@antv/g2";
import DataSet from "@antv/data-set";

$(() => {
  console.log("first chart");

  function renderFirst() {
    const data = [
      { genre: "Sports", sold: 275 },
      { genre: "Strategy", sold: 115 },
      { genre: "Action", sold: 120 },
      { genre: "Shooter", sold: 350 },
      { genre: "Other", sold: 150 }
    ]; // G2 对数据源格式的要求，仅仅是 JSON 数组，数组的每个元素是一个标准 JSON 对象。
    // Step 1: 创建 Chart 对象
    const chart = new G2.Chart({
      container: "c1", // 指定图表容器 ID
      width: 600, // 指定图表宽度
      height: 300 // 指定图表高度
    });
    // Step 2: 载入数据源
    chart.source(data);
    // Step 3：创建图形语法，绘制柱状图，由 genre 和 sold 两个属性决定图形位置，genre 映射至 x 轴，sold 映射至 y 轴
    chart
      .interval()
      .position("genre*sold")
      .color("genre");
    // Step 4: 渲染图表
    chart.render();
  }
  // renderFirst();

  function registerCloud() {
    function getTextAttrs(cfg) {
      return _.assign({}, cfg.style, {
        fillOpacity: cfg.opacity,
        fontSize: cfg.origin._origin.size,
        rotate: cfg.origin._origin.rotate,
        text: cfg.origin._origin.text,
        textAlign: "center",
        fontFamily: cfg.origin._origin.font,
        fill: cfg.color,
        textBaseline: "Alphabetic"
      });
    }

    // 给point注册一个词云的shape
    G2.Shape.registerShape("point", "cloud", {
      drawShape: function drawShape(cfg, container) {
        var attrs = getTextAttrs(cfg);
        return container.addShape("text", {
          attrs: _.assign(attrs, {
            x: cfg.x,
            y: cfg.y
          })
        });
      }
    });
  }
  registerCloud();

  function renderCloud() {
    $.getJSON("/assets/data/world-population.json", function(data) {
      var dv = new DataSet.View().source(data);
      var range = dv.range("value");
      var min = range[0];
      var max = range[1];
      dv.transform({
        type: "tag-cloud",
        fields: ["x", "value"],
        // size: [window.innerWidth, window.innerHeight],
        size: [600, 300],

        font: "Verdana",
        padding: 0,
        timeInterval: 5000, // max execute time
        rotate: function rotate() {
          var random = ~~(Math.random() * 4) % 4;
          if (random == 2) {
            random = 0;
          }
          return random * 90; // 0, 90, 270
        },
        fontSize: function fontSize(d) {
          if (d.value) {
            return ((d.value - min) / (max - min)) * (80 - 24) + 24;
          }
          return 0;
        }
      });
      var chart = new G2.Chart({
        container: "idCloud",
        // width: window.innerWidth,
        // height: window.innerHeight,

        width: 600,
        height: 300,

        padding: 0
      });
      chart.source(dv, {
        x: {
          nice: false
        },
        y: {
          nice: false
        }
      });
      chart.legend(false);
      chart.axis(false);
      chart.tooltip({
        showTitle: false
      });
      chart.coord().reflect();
      chart
        .point()
        .position("x*y")
        .color("category")
        .shape("cloud")
        .tooltip("value*category");
      chart.render();
    });
  }
  // renderCloud();

  function renderWordCloud() {
    $.getJSON("/chart/get_wordcloud", function(data) {
      var dv = new DataSet.View().source(data);
      var range = dv.range("count");
      var min = range[0];
      var max = range[1];
      dv.transform({
        type: "tag-cloud",
        fields: ["word", "count"],

        size: [1000, 600],

        font: "Verdana",
        padding: 0,
        timeInterval: 5000, // max execute time
        rotate: function rotate() {
          var random = ~~(Math.random() * 4) % 4;
          if (random == 2) {
            random = 0;
          }
          return random * 90; // 0, 90, 270
        },
        fontSize: function fontSize(d) {
          if (d.value) {
            return ((d.value - min) / (max - min)) * (80 - 24) + 24;
          }
          return 0;
        }
      });
      var chart = new G2.Chart({
        container: "idWordCloud",

        width: 1000,
        height: 600,

        padding: 0
      });
      chart.source(dv, {
        x: {
          nice: false
        },
        y: {
          nice: false
        }
      });
      chart.legend(false);
      chart.axis(false);
      chart.tooltip({
        showTitle: false
      });
      chart.coord().reflect();
      chart
        .point()
        .position("x*y")
        .color("word")
        .shape("cloud")
        .tooltip("count");
      chart.render();
    });
  }
  renderWordCloud();
});
