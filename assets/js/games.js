// require("expose-loader?$!expose-loader?jQuery!jquery");

import axios from "axios";
import * as R from "ramda";

$(() => {
    function initScreen() {
        $("#criterionRaw").hide();
        $("#criterionDisplay").show();

        $("#optionsRaw").hide();
        $("#optionsDisplay").show();
    }
    initScreen();

    function noInstance() {
        const game_id = document.getElementById("game_id").innerText.trim();
        // console.log("game", game_id);
        return (game_id.length === 0);
    }

    // create or find project by name
    $("#btnCreate").on("click", createHandle);
    function createHandle() {
        const name = document.getElementById("game_name").value;
        if (name.length === 0) {
            layer.alert("请输入待决策问题");
            return;
        }

        const param = {
            name: name
        };
        // console.log(param);
        axios.post("/games/create", param).then(function(res) {
            // console.log(res);
            // 2019-04-19T09:16:28.44001Z  -> 2019-04-19 09:16:28
            const { data } = res;
            const raw = data.created_at;
            const dt = raw.substring(0, 19).replace("T", " ");
            document.getElementById("created_at").innerText = dt;
            document.getElementById("game_id").innerText = data.id;
            document.getElementById("game_inst").innerText = JSON.stringify(data);

            const critList = data.criterion.split("\n").filter(i => i.length > 0);
            renderCriterion({ raw: data.criterion, items: critList });

            const optionList = data.options.split("\n").filter(i => i.length > 0);
            renderOptions({ raw: data.options, items: optionList });

            $("#btnCriterionReload").trigger("click");
            $("#btnOptionsFire").trigger("click");
            $("#navItem-0").trigger("click");
        });
    }
    
    // criterion: save & display
    function renderCriterion({ raw, items }) {
        document.getElementById("criterion_text").value = raw;

        const tags = items.map((i, idx) => {
            return `<tr>
                    <td>${idx + 1}</td>
                    <td>${i}</td>
                </tr>`;
        });
        $("#criterionTable").html(tags);
    }

    $("#btnCriterion").on("click", criterionHandle);
    function criterionHandle() {
        const criterion = document.getElementById("criterion_text").value.trim();

        if (criterion.length === 0) {
            layer.alert("请输入选择标准，每行一个");
            return;
        }

        const param = {
            id: document.getElementById("game_id").innerText,
            criterion
        };
        // console.log(param);
        axios.post("/games/saveCriterion", param).then(function(res) {
            const { data } = res;

            renderCriterion(data);

            $("#criterionRaw").hide();
            $("#criterionDisplay").show();
        });
    }
    $("#btnCriterionModify").on("click", criterionModHandle);
    function criterionModHandle() {
        if (noInstance()) {
            layer.alert("请先输入问题");
            return;
        }
        $("#criterionRaw").show();
        $("#criterionDisplay").hide();
    }

    //  options: save & display
    function renderOptions({ raw, items }) {
        document.getElementById("options_text").value = raw;

        const tags = items.map((i, idx) => {
            return `<tr>
                    <td>${idx + 1}</td>
                    <td>${i}</td>
                </tr>`;
        });
        $("#optionsTable").html(tags);
    }

    $("#btnOptions").on("click", optionsHandle);
    function optionsHandle() {
        const options = document.getElementById("options_text").value.trim();

        if (options.length === 0) {
            layer.alert("请输入可选项，每行一个");
            return;
        }

        const param = {
            id: document.getElementById("game_id").innerText,
            options: options
        };
        // console.log(param);
        axios.post("/games/saveOptions", param).then(function(res) {
            // console.log(res);
            const { data } = res;
            renderOptions(data);

            $("#optionsRaw").hide();
            $("#optionsDisplay").show();
        });
    }
    $("#btnOptionsModify").on("click", optionsModHandle);
    function optionsModHandle() {
        if (noInstance()) {
            layer.alert("请先输入问题");
            return;
        }

        $("#optionsRaw").show();
        $("#optionsDisplay").hide();
    }

    // genEmptyMatrix generate matrix element with empty string
    function genEmptyMatrix(count) {
        const result = [];
        for (let i = 1; i <= count; i++) {
            const row = [];
            for (let j = 1; j <= count; j++) {
                if (i === j) {
                    row.push("1");
                } else {
                    row.push("");
                }
            }
            result.push(row);
        }
        return result;
    }

    // allowed number in compare matrix
    const permitted = ["1", "2", "3","4", "5", "6", "7","8", "9"];
    // -2.2342354234 -> "2.2342"
    function fmtString(s, p) {
        const a = s > 0 ? s : -s;
        return a.toFixed(p).toString();
    }

    // criterion priority
    function criterionHot() {
        let hot;

        function createPairs(reset) {
            const count = $("#criterionTable tr").length;
            if (count === 0) {
                console.log("no criterion.");
                return;
            }

            const labels = [];
            for (let i = 1; i <= count; i++) {
                labels.push(`${i}`);
            }
            // console.log(labels);

            var container = document.getElementById("tblCriterion");

            if (hot) {
                hot.destroy();
            }

            if (reset) {
                // 重新来过
                const pairs = genEmptyMatrix(count);

                hot = new Handsontable(container, {
                    licenseKey: "non-commercial-and-evaluation",
                    data: pairs,
                    rowHeaders: labels,
                    colHeaders: labels,
                    filters: false,
                    dropdownMenu: false,
                    beforeChange,
                    beforeKeyDown
                });
                renderResult("");
                return;
            }

            // 从数据库中加载
            const param = {
                id: document.getElementById("game_id").innerText
            };
            // axios.get(`/games/byid?id=${param.id}`)
            axios.post("/games/byid", param).then(function(res) {
                const { data } = res;
                const { entry, result } = data;

                let pairs;
                if (entry.pairs === "") {
                    pairs = genEmptyMatrix(count);
                } else {
                    pairs = JSON.parse(entry.pairs);
                    if (pairs.length === 0) {
                        pairs = genEmptyMatrix(count);
                    }
                }

                hot = new Handsontable(container, {
                    licenseKey: "non-commercial-and-evaluation",
                    data: pairs,
                    rowHeaders: labels,
                    colHeaders: labels,
                    filters: false,
                    dropdownMenu: false,
                    beforeChange,
                    beforeKeyDown
                });

                renderResult(result);
            });
        }

        function beforeChange(changes, source) {
            // [[row, prop, oldVal, newVal], ...]
            // console.log(changes[0]);
            const [r, c, d, v] = changes[0];
            const a = permitted.findIndex(i => i === v);
            if (a < 0) {
                return;
            }
            hot.setDataAtRowProp(c, r, `1/${v}`);
        }

        function beforeKeyDown(evt) {
            const a = permitted.findIndex(i => i === evt.key);
            // console.log(evt.key, a);
            if (a < 0) {
                evt.stopImmediatePropagation();
                return false;
            }
        }

        $("#btnCriterionReset").on("click", resetCriterionHandle);
        function resetCriterionHandle() {
            createPairs(true);
        }

        $("#btnCriterionReload").on("click", reloadCriterionHandle);
        function reloadCriterionHandle() {
            createPairs(false);
        }

        $("#btnCriterionSave").on("click", criterionSave);
        function criterionSave() {
            const d = hot.getData();
            const p = JSON.stringify(d);
            // console.log(p);

            const param = {
                id: document.getElementById("game_id").innerText,
                pairs: p
            };
            // console.log(param);
            axios.post("/games/saveCriterionPairs", param).then(function(res) {
                // console.log(res);
                const { data } = res;
                renderResult(data.result);
            });
        }

        function renderResult(result) {
            // console.log(result);
            const target = $("#criterionResult");
            if (result === "") {
                target.html("");
                return;
            }

            if (result.lambdaMax !== 0) {
                const a = `<p>CR = ${fmtString(result.cr, 4)}</p>`;
                const w = result.weights.map(w => fmtString(w, 4)).join(" ");
                const msg = `${a}<p>Weight: [${w}]</p>`;
                target.html(msg);
                return;
            }
            target.html("");
        }
    }
    criterionHot();


    // 在某个标准下，选项间的重要度
    function optionsCtrl() {
        function renderNav() {
            const l =  $("#criterionTable tr");
            const msgs = [];
            l.each((idx, item) => {
                const a = $(item);
                const b = a.children();
                const c = b[1].innerText;

                const s = `<li class="nav-item"><a class="nav-link navCritItem" href="#" id="navItem-${idx}" data-seq=${idx}>${c}</a></li>`
                msgs.push(s);
            });
            const tags = msgs.join(" ");
            $("#navCriterion").html(tags);
            $(".navCritItem").on("click", clickNavHandle);
        }
        $("#btnOptionsFire").on("click", renderNav);

        function clickNavHandle() {
            const self = $(this);
            const currSeq = self.data("seq");
            // console.log('clicked: ', self.data("seq"));
            document.getElementById("currCriteria").innerText = currSeq;

            if (self.hasClass('active')) {
                console.log('already current');
                return false;
            }

            $(".navCritItem").removeClass("active");
            self.addClass("active");

            // 切换标签后，模拟点击 “加载”
            reloadOptionsHandle();

            return false;
        }

        let hot;

        function createPairs(reset) {
            const count = $("#optionsTable tr").length;
            if (count === 0) {
                console.log("no criterion.");
                return;
            }

            const labels = [];
            for (let i = 1; i <= count; i++) {
                labels.push(`${i}`);
            }
            // console.log(labels);

            const container = document.getElementById("tblOptionCompare");

            if (hot) {
                hot.destroy();
            }

            if (reset) {
                // 重新来过
                const pairs = genEmptyMatrix(count);

                hot = new Handsontable(container, {
                    licenseKey: "non-commercial-and-evaluation",
                    data: pairs,
                    rowHeaders: labels,
                    colHeaders: labels,
                    filters: false,
                    dropdownMenu: false,
                    beforeChange,
                    beforeKeyDown
                });

                renderResult("");
                return;
            }

            // 从数据库中加载
            const param = {
                id: document.getElementById("game_id").innerText,
                seq: document.getElementById("currCriteria").innerText,
            };

            axios.post("/games/loadOptionPairs", param).then(function(res) {
                const { data } = res;
                const { item, result }  = data;
                let pairs;
                if (item.pairs === "") {
                    pairs = genEmptyMatrix(count);
                } else {
                    pairs = JSON.parse(item.pairs);
                    if (pairs.length === 0) {
                        pairs = genEmptyMatrix(count);
                    }
                }

                hot = new Handsontable(container, {
                    licenseKey: "non-commercial-and-evaluation",
                    data: pairs,
                    rowHeaders: labels,
                    colHeaders: labels,
                    filters: false,
                    dropdownMenu: false,
                    beforeChange,
                    beforeKeyDown
                });

                renderResult(result);
            });
        }

        function beforeChange(changes, source) {
            // [[row, prop, oldVal, newVal], ...]
            // console.log(changes[0]);
            const [r, c, d, v] = changes[0];
            const a = permitted.findIndex(i => i === v);
            if (a < 0) {
                return;
            }
            hot.setDataAtRowProp(c, r, `1/${v}`);
        }

        function beforeKeyDown(evt) {
            const a = permitted.findIndex(i => i === evt.key);
            // console.log(evt.key, a);
            if (a < 0) {
                evt.stopImmediatePropagation();
                return false;
            }
        }

        $("#btnOptionsReset").on("click", resetOptionsHandle);
        function resetOptionsHandle() {
            createPairs(true);
        }

        $("#btnOptionsReload").on("click", reloadOptionsHandle);
        function reloadOptionsHandle() {
            createPairs(false);
        }

        $("#btnOptionsSave").on("click", optionsSave);
        function optionsSave() {
            const d = hot.getData();
            const p = JSON.stringify(d);
            // console.log(p);

            const param = {
                id: document.getElementById("game_id").innerText,
                seq: document.getElementById("currCriteria").innerText,
                pairs: p
            };
            // console.log(param);
            axios.post("/games/saveOptionPairs", param).then(function(res) {
                // console.log(res);
                const { data } = res;

                renderResult(data.result);
            });
        }

        function renderResult(result) {
            // console.log(result);
            const target = $("#optionsResult");
            if (result === "") {
                target.html("");
                return;
            }

            if (result.lambdaMax !== 0) {
                const a = `<p>CR = ${fmtString(result.cr, 4)}</p>`;
                const w = result.weights.map(w => fmtString(w, 4)).join(" ");
                const msg = `${a}<p>Weight: [${w}]</p>`;
                target.html(msg);
                return;
            }
            target.html("");
        }
    }
    optionsCtrl();

    // calc final result
    function finalComponent() {
        $("#btnFinal").on("click", finalHandle);
        function finalHandle() {
            // console.log("final clicked");
            const param = {
                id: document.getElementById("game_id").innerText,
            };
            axios.post("/games/calcFinal", param).then(function(res) {
                const {data} = res;
                renderResult(data);
            });
        }

        function renderResult(result) {
            // console.log(result);
            const {score, options, status} = result;

            const errMsg = [
                {code: 5, message: "请先输入问题"},
                {code: 6, message: "请先输入问题"},
                {code: 7, message: "请先填写完成 选择标准的重要度"},
                {code: 10, message: "请先填写完成 选择标准的重要度"},
                {code: 20, message: "可选项在每个标准下都需要有重要度"},
                {code: 30, message: "请重新检查选项重要性，确保每个标准下均已完成"},
            ];

            const err = errMsg.find(e => e.code ===status );
            if (err !== undefined) {
                // console.log(err.message);
                layer.alert(err.message);
                return;
            }

            const list = score.map((s, i)=> ({s, o: options[i]}));
            const final = R.sort((a, b) => b.s - a.s, list);

            const s = final.map((i, idx) => {
                return `<tr>
                    <td>${idx + 1}</td>
                    <td>${i.o}</td>
                    <td>${fmtString(i.s, 4)}</td>
                </tr>`
            });
            const tags = s.join(" ");
            $("#finalResult").html(tags);
        }
    }
    finalComponent();
});
