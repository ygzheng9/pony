import G2 from "@antv/g2";
import DataSet from "@antv/data-set";

$(() => {
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

  // word cloud
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

        forceFit: true,
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

  // frequence
  function renderFreq() {
    $.getJSON("/chart/get_wordfreq", function(data) {
      var chart = new G2.Chart({
        container: 'idWordFreq',
        forceFit: true,
        // width: 1000,
        height: 600,
        padding: [10, 10, 50, 124]
      });

      data = data.filter(a => a.level > 1);
      // console.log(data.length);

      chart.source(data, {
        level: {
          alias: '提及的人数'
        },
        count: {
          max: 1000,
          min: 0,
          nice: false,
          alias: '提及的次数'
        }
      });

      chart.axis('level', {
        label: {
          textStyle: {
            fill: '#8d8d8d',
            fontSize: 12
          },
        },
        tickLine: {
          alignWithLabel: false,
          length: 0
        },
        line: {
          lineWidth: 0
        },
        title: {
          position: 'center',
        }
      });
      chart.axis('count', {
        label: null,
        title: {
          offset: 30,
          textStyle: {
            fontSize: 12,
            fontWeight: 300
          },
          position: "center",
        }
      });
      chart.legend(false);
      chart.coord().transpose();
      chart.interval().position('level*count').size(26).opacity(1).label('count', {
        textStyle: {
          fill: '#8d8d8d'
        },
        offset: 10
      });
      chart.render();
    });
  }
  renderFreq();

  function renderDist() {
    $.getJSON('/chart/get_worddist', function(data) {
      data = data.filter(d => d.wc_count > 20);

      var ds = new DataSet();
      var dv = ds.createView('dist').source(data);
      dv.transform({
        type: 'bin.histogram',
        field: 'wc_count',
        binWidth: 20, // 在此修改矩形的宽度，代表真实数值的大小
        as: ['wc_count', 'count']
      });

      var chart = new G2.Chart({
        container: 'idWordDist',
        forceFit: true,
        height: window.innerHeight
      });

      chart.source(dv, {
        depth: {
          tickInterval: 4
        },
        count: {
          alias: '关键词数量'
        },
        wc_count: {
          alias: '重复次数'
        },
      });

      chart.axis('count', {
        title: {
          position: "center",
        }
      });
      chart.axis('wc_count', {
        title: {
          position: "center",
        }
      });

      chart.tooltip({
        crosshairs: false,
        position: 'top',
        inPlot: false
      });


      chart.interval().position('wc_count*count');
      chart.render();
    });
  }
  renderDist();
});
