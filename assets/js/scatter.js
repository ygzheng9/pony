import G2 from "@antv/g2";

import * as R from "ramda";

$(() => {
    function renderScore() {
        $.getJSON('/scatter/data', function(data) {
            renderData(data);
        });

        function renderData(data) {
            // console.log(data);
            const tx = transform(data.items);
            const transData = tx.data;

            const chart = new G2.Chart({
                container: 'scatterData',
                forceFit: true,
                height: window.innerHeight
            });
            chart.source(transData);

            const itemCnt = transData.length;
            // 为各个字段设置别名
            chart.scale({
                x: {
                    type: "linear",
                    alias: '重要性',
                    // min: 1, max: itemCnt,
                },
                y: {
                    type: "linear",
                    alias: '紧迫度',
                    // min: 1, max: itemCnt,
                },
                sd_score: {
                    type: "linear",
                    alias: '集中度'
                },
                item: {
                    alias: '项目'
                },
            });

            chart.tooltip({
                showTitle: false // 不显示默认标题
            });

            chart.axis('x', {
                title: {
                    position: "center",
                    textStyle: {
                        fontSize: 12,
                        fill: '#404040', // 文本的颜色
                    },
                },
                label: {
                    formatter: function formatter(value) {
                        return (value * 1.0).toFixed(2);
                    }
                }
            });

            chart.axis('y', {
                title: {
                    position: "center",
                    textStyle: {
                        fontSize: 12,
                        fill: '#404040', // 文本的颜色
                    },
                },
                label: {
                    formatter: function formatter(value) {
                        return (value * 1.0).toFixed(2);
                    }
                }
            });

            const c = chart.point().position('x*y')
                .shape('circle').tooltip('item*x*y*sd_score')
                .label('item', {
                    type: "scatter",
                    // 设置文本的显示样式，还可以是个回调函数，回调函数的参数为该坐标轴对应字段的数值
                    textStyle: {
                        // textAlign: 'center', // 文本对齐方向，可取值为： start middle end
                        // textBaseline: 'bottom', // 文本基准线，可取 top middle bottom，默认为middle
                        // fill: '#404040', // 文本的颜色
                        fill: 'red', // 文本的颜色
                        fontSize: '16', // 文本大小
                        fontWeight: 'bold', // 文本粗细
                        // rotate: 30,
                    }
                });

            // 差别是是否有 size
            if (tx.type === "SD") {
                c.size('sd_score', [10, 50]);
            }

            chart.render();
        }

        function transform(items) {
            // 形成所有 item 的 list
            const itemList = items.reduce((acc, curr) => {
                const exist = acc.findIndex(i => i.item === curr.item);
                if (exist < 0) {
                    acc.push({item: curr.item, al: [], bl: []});
                }
                return acc;
            }, []);
            // console.log("itemList init: ", itemList);

            // 分组，每组内，对每个标准，从小到大排序，取序号+1
            // 小代表靠前，代表重要，在图形上需要显示到 右上角，所以，后面还有反转
            const byGroup = R.groupBy(i => i.group, items);
            const entries = Object.entries(byGroup);
            const seqList = entries.map(([k, v]) => {
                const a = R.sort((a1, a2) => a1.ind_a - a2.ind_a, v);
                const al = a.map((i, idx) => ({...i, seq_a: idx+1}));

                const b = R.sort((b1, b2) => b1.ind_b - b2.ind_b, v);
                const bl = b.map((i, idx) => ({...i, seq_b: idx+1}));

                return {
                    k, al, bl
                }
            });

            // 对每个 item，取出所有 组 的排序；
            seqList.forEach( i => {
                i.al.forEach(a => {
                    const idx = itemList.findIndex(i => i.item === a.item);
                    itemList[idx].al.push(a.seq_a);
                });

                i.bl.forEach(b => {
                    const idx = itemList.findIndex(i => i.item === b.item);
                    itemList[idx].bl.push(b.seq_b);
                });
            });
            // console.log("itemList after: ", itemList);

            // 计算均值，方差
            const groupCnt = entries.length;
            const itemCnt = itemList.length;
            const avgList = itemList.map(i => {
                const avg_a = R.sum(i.al) / groupCnt;
                const avg_b = R.sum(i.bl) / groupCnt;

                const sd_a = R.sum(i.al.map(a => Math.pow(Math.abs(a - avg_a), 2)));
                const sd_b = R.sum(i.bl.map(b => Math.pow(Math.abs(b - avg_b), 2)));

                // 反转坐标轴，排序小的，转换到右上角
                return {...i,
                    x: itemCnt - avg_a,
                    y: itemCnt - avg_b,
                    sd_x: sd_a,
                    sd_y: sd_b,
                    sd: sd_a + sd_b};
            });

            // 计算总方差
            const sd_total = avgList.reduce((acc, curr) => {
                return acc + curr.sd
            }, 0);
            // console.log("total sd: ", sd_total);

            // const r = avgList.map(i => ({...i, sd_score: sd_total / i.sd}));
            const r = avgList.map(i => {
                return {
                    item: i.item,
                    x: fmtNumber(i.x, 4),
                    y: fmtNumber(i.y, 4),
                    sd_score: fmtNumber(sd_total / i.sd, 4),
                };
            });

            let type = "SD";
            // 只有一个数据，没有方差，显示时，都是点，没有大小 size
            if (sd_total === 0) {
                type = "NONE";
            }

            return {
                data: r,
                type,
            };
        }
    }
    renderScore();

    // 格式化数字
    function fmtNumber(n, p) {
        return n.toFixed(p);
    }
});